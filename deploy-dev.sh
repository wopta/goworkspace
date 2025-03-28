echo "TAG_NAME: $1"
echo "FX_NAME: $2"
echo "FX_ENTRYPOINT: $3"
echo "TAG_VERSION: $4"
TAG_NAME=$1
FX_NAME=$2
FX_ENTRYPOINT=$3
TAG_VERSION=$4

echo "current tag: ${TAG_NAME}"
echo "setting DEV environment variables..."
bucket=function-data
region=europe-west1
project=positive-apex-350507
env=dev
genFx=--gen2
sa=wopta-dev-cloudbuild-sa@positive-apex-350507.iam.gserviceaccount.com
timeout=60
vpc=functions-connector

# === COPY ASSETS FROM BUCKET ==================================================
echo "copying assets..."
mkdir -p /workspace/${FX_NAME}/tmp/assets
gsutil -m cp -r gs://${bucket}/assets/documents/** /workspace/${FX_NAME}/tmp/assets
cp /workspace/.gcloudignore /workspace/${FX_NAME}/.gcloudignore

# === SET INGRESS SETTINGS BY FUNCTION =========================================
ingress=internal-and-gclb
if [[ "${FX_NAME}" == "callback" ]]; then
  ingress=all
fi

# === SET MEMORY BY FUNCTION ===================================================
mem=256Mb
if [[ "${FX_NAME}" == "broker" || "${FX_NAME}" == "quote" ]]; then
  mem=1Gb
fi

# === DEPLOY ===================================================================
echo "deploying function ${FX_NAME}..."
gcloud functions deploy ${FX_NAME} \
  --project=${project} \
  --region=${region} \
  --source=/workspace/${FX_NAME} \
  --entry-point=${FX_ENTRYPOINT} \
  --trigger-http \
  --allow-unauthenticated \
  --run-service-account=${sa} \
  --runtime=go121 \
  --env-vars-file ${env}.yaml \
  --timeout=${timeout} \
  --ingress-settings=${ingress} \
  --egress-settings=all \
  --memory=${mem} \
  --vpc-connector=${vpc} \
  ${genFx} 

# === APPLY TAG VERSION LABEL ==================================================
echo "DEV run services update..."
gcloud run services update ${FX_NAME} --update-labels tagversion=${TAG_VERSION} --service-account=${sa} --region=${region}
