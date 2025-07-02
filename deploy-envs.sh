TAG_NAME=$1

# === EXTRACT FUNCTION NAME AND VERSION FROM TAG ===============================
# Read TAG_NAME splitting by the IFS separator into the needed variables
IFS='/'
read -a strarr <<< $TAG_NAME
FX_NAME=$(echo ${strarr[0]})
TAG=$(echo ${strarr[1]})
echo "FX_NAME: ${FX_NAME}"
echo "TAG: ${TAG}"

# === FX ENTRYPOINT ============================================================
# Camel case function entrypoint. Ex.: broker -> Broker 
FX_ENTRYPOINT=$(sed -r 's/(^|-)(\w)/\U\2/g' <<<"${FX_NAME}")
echo "FX_ENTRYPOINT: ${FX_ENTRYPOINT}"

# === TAG VERSION ==============================================================
# Replace '.' with '_' for tag version. Ex.: 1.2.3.dev -> 1_2_3_dev 
TAG_VERSION=$(echo "${TAG}" | sed -r 's/\./_/g')
echo "TAG_VERSION: ${TAG_VERSION}"

# === DEV ENV VARS =============================================================
if [[ "${TAG}" == *"dev"* ]]; then
    echo "Setting DEV environment variables..."
    BUCKET=function-data
    REGION=europe-west1
    PROJECT=positive-apex-350507
    ENV=dev
    GEN_FX=--gen2
    SERVICE_ACCOUNT=wopta-dev-cloudbuild-sa@positive-apex-350507.iam.gserviceaccount.com
    TIMEOUT=60
    VPC=functions-connector
fi

# === UAT ENV VARS =============================================================
if [[ "${TAG}" == *"uat"* ]]; then
    echo "Setting UAT environment variables..."
    BUCKET=core-452909-function-data
    REGION=europe-west1
    PROJECT=core-452909
    ENV=uat
    GEN_FX=--gen2
    SERVICE_ACCOUNT=wopta-uat-cloudbuild-sa@core-452909.iam.gserviceaccount.com
    TIMEOUT=60
    VPC=functions-connector
fi

# === PROD ENV VARS ============================================================
if [[ "${TAG}" == *"prod"* ]]; then
    echo "Setting PROD environment variables..."
    BUCKET=core-350507-function-data
    REGION=europe-west1
    PROJECT=core-350507
    ENV=prod
    GEN_FX=""
    SERVICE_ACCOUNT=wopta-prod-cloudbuild-sa@core-350507.iam.gserviceaccount.com
    TIMEOUT=520
    VPC=functions-connector
fi

GCPTOKEN=$(gcloud auth print-access-token)

VGOPROXY="https://oauth2accesstoken:${GCPTOKEN}@europe-west1-go.pkg.dev/positive-apex-350507/goworkspace"
VGOPROXY+=,https://proxy.golang.org,direct
VGONOSUMDB="gitlab.dev.wopta.it/goworkspace/*"

BUILD_ENV_VARS="^#^GOPROXY=${VGOPROXY}"
BUILD_ENV_VARS+="#GONOSUMDB=${VGONOSUMDB}"

# === COPY ASSETS FROM BUCKET ==================================================
echo "Copying assets..."
mkdir -p /workspace/${FX_NAME}/tmp/assets
gsutil -m cp -r gs://${BUCKET}/assets/documents/** /workspace/${FX_NAME}/tmp/assets
cp /workspace/.gcloudignore /workspace/${FX_NAME}/.gcloudignore

# === SET INGRESS SETTINGS BY FUNCTION =========================================
INGRESS=internal-and-gclb
if [[ "${FX_NAME}" == "callback" ]]; then
    INGRESS=all
fi

# === SET MEMORY BY FUNCTION ===================================================
MEMORY=256Mb
if [[ "${FX_NAME}" == "broker" || "${FX_NAME}" == "quote" || "${FX_NAME}" == "mga" || "${FX_NAME}" == "payment" || "${FX_NAME}" == "callback" || "${FX_NAME}" == "inclusive" ]]; then
    MEMORY=1Gb
elif [[ "${FX_NAME}" == "companydata" ]]; then
    MEMORY=2Gb
fi

# === DEPLOY ===================================================================
echo "Deploying function ${FX_NAME}..."
gcloud functions deploy ${FX_NAME} \
    --project=${PROJECT} \
    --region=${REGION} \
    --source=/workspace/${FX_NAME} \
    --entry-point=${FX_ENTRYPOINT} \
    --trigger-http \
    --allow-unauthenticated \
    --run-service-account=${SERVICE_ACCOUNT} \
    --build-service-account=projects/${PROJECT}/serviceAccounts/${SERVICE_ACCOUNT} \
    --runtime=go121 \
    --env-vars-file ${ENV}.yaml \
    --timeout=${TIMEOUT} \
    --ingress-settings=${INGRESS} \
    --egress-settings=all \
    --memory=${MEMORY} \
    --vpc-connector=${VPC} \
    --verbosity=debug \
    --update-build-env-vars "${BUILD_ENV_VARS}" \
    ${GEN_FX} 

# === APPLY TAG VERSION LABEL ==================================================
if [[ "${GEN_FX}" == "--gen2" ]]; then
    echo "Run services update..."
    gcloud run services update ${FX_NAME} \
        --update-labels tagversion=${TAG_VERSION} \
        --service-account=${SERVICE_ACCOUNT} \
        --region=${REGION}
fi
