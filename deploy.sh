    echo "First arg: $1"
    TAG_NAME=$1
     #export namefx=( $(grep -Eo '[[:digit:]]+|[^[:digit:]]+' <<<'$TAG_NAME') )
    #Print current tag
    echo "current tag: ${TAG_NAME}"

    # === EXTRACT FUNCTION NAME AND VERSION FROM TAG ===========================
    # Read TAG_NAME splitting by the IFS separator into the needed variables
    # ex.: broker/1.0.0.dev => [fx_name=broker, tag=1.0.0.dev]
    IFS='/'
    read -a strarr <<< $TAG_NAME
    fx_name=$(echo ${strarr[0]})
    tag=$(echo ${strarr[1]})
    echo "fx_name: ${fx_name}"
    echo "tag: ${tag}"

    # === FX ENTRYPOINT AND TAG VERSION ========================================
    # Camel case function entrypoint. Ex.: broker -> Broker 
    fx_entrypoint=$(sed -r 's/(^|-)(\w)/\U\2/g' <<<"${fx_name}")
    # Replace '.' with '_' for tag version. Ex.: 1.2.3.dev -> 1_2_3_dev 
    tag_version=$(echo "${tag}" | sed -r 's/\./_/g')
    echo "fx_entrypoint: ${fx_entrypoint}"
    echo "tag_version: ${tag_version}"

    # === SET DEV ENV VARS =====================================================
    if [[ "${tag}" == *"dev"* ]]; then
      echo "setting DEV environment variables..."
      bucket=function-data
      region=europe-west1
      project=positive-apex-350507
      env=dev
      genFx=--gen2
      sa=wopta-dev-cloudbuild-sa@positive-apex-350507.iam.gserviceaccount.com
      timeout=60
      vpc=functions-connector
    fi
    # === SET DEV UAT VARS =====================================================
    if [[ "${tag}" == *"uat"* ]]; then
      echo "setting UAT environment variables..."
      bucket=function-data
      region=europe-west1
      mem=256Mb
      project=core-452909
      env=uat
      genFx=--gen2
      sa=wopta-dev-cloudbuild-sa@positive-apex-350507.iam.gserviceaccount.com
      timeout=60
      vpc=functions-connector
    fi
    # === SET PROD ENV VARS ====================================================
    if [[ "${tag}" == *"prod"* ]]; then
      echo "setting PROD environment variables..."
      bucket=core-350507-function-data
      region=europe-west1
      project=core-350507
      env=prod
      genFx=""
      sa=wopta-prod-cloud-function@core-350507.iam.gserviceaccount.com
      timeout=520
      vpc=functions-connector
    fi

    # === SET INGRESS SETTINGS BY FUNCTION =====================================
    ingress=internal-and-gclb
    if [[ "${fx_name}" == "callback" ]]; then
      ingress=all
    fi

    # === SET MEMORY BY FUNCTION ===============================================
    mem=256Mb
    if [[ "${fx_name}" == "broker" || "${fx_name}" == "quote" || "${fx_name}" == "mga" || "${fx_name}" == "payment" || "${fx_name}" == "callback" || "${fx_name}" == "inclusive" ]]; then
      mem=1Gb
    elif [[ "${fx_name}" == "companydata" ]]; then
      mem=2Gb
    fi

    # === COPY ASSETS FROM BUCKET ==============================================
    echo "copying assets..."
    mkdir -p /workspace/${fx_name}/tmp/assets
    gsutil -m cp -r gs://${bucket}/assets/documents/** /workspace/${fx_name}/tmp/assets
    cp /workspace/.gcloudignore /workspace/${fx_name}/.gcloudignore 

    # === DEPLOY ===============================================================
    echo "deploying function ${fx_name}..."
    gcloud functions deploy ${fx_name} \
      --project=${project} \
      --region=${region} \
      --source=/workspace/${fx_name} \
      --entry-point=${fx_entrypoint} \
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
   
    # === APPLY TAG VERSION LABEL ==============================================
    if [[ "${env}" == "dev" ]]; then
      echo "DEV run services update..."
      gcloud run services update ${fx_name} --update-labels tagversion=${tag_version} --service-account=${sa} --region=${region}
    fi