    echo "First arg: $1"
    TAG_NAME=$1
     #export namefx=( $(grep -Eo '[[:digit:]]+|[^[:digit:]]+' <<<'$TAG_NAME') )
    #Print current tag
    echo "current tag: ${TAG_NAME}"
   #-------------SPLIT Tag-------------------------------------------------------
    #Read the string value
    echo $namefx
    read text
    # Set comma as delimiter
    IFS='/'
    #Read the split words into an array based on comma delimiter
    read -a strarr <<< $TAG_NAME
    #Print the splitted words
    echo "namefx : ${strarr[0]}"
    echo "tag : ${strarr[1]}"
   #-------------------------NAME FX Cammel case ad tag ver-------------------------------------------
    namefx=$strarr[0]
    name_camel=$(sed -r 's/(^|-)(\w)/\U\2/g' <<<"${strarr[0]}")
    tagVersion=$(echo "${strarr[1]}" | sed -r 's/\./_/g')
    t=$(tr -s . _ <<< "${strarr[1]}")
    echo "namefx camel : ${name_camel}"
    echo "tagVersion camel : ${tagVersion}"
    echo "tagVersion camel : ${t}"
   #-----------------SET VAR DEV---------------------------------------------------
    if [[ "${strarr[1]}" == *"dev"* ]]; then
      echo "dev enviroment"
      bucket=function-data
      project=positive-apex-350507
      env=dev
      genFx=--gen2
      sa=wopta-dev-cloudbuild-sa@positive-apex-350507.iam.gserviceaccount.com
      timeout=60
      vpc=wopta-dev-custom-vpc
    fi
   #---------------------SET VAR PROD-----------------------------------------------
    if [[ "${strarr[1]}" == *"prod"* ]]; then
    echo "prod enviroment"
      bucket=core-350507-function-data
      project=core-350507
      env=prod
      genFx=""
      sa=wopta-prod-cloud-function@core-350507.iam.gserviceaccount.com
      timeout=520
      vpc=prod-custom-vpc
    fi

    #--------- Copy assets folder from Google Bucket to directory----------------------------------
    mkdir -p /workspace/${strarr[0]}/tmp/assets
    gsutil -m cp -r gs://${bucket}/assets/documents/** /workspace/${strarr[0]}/tmp/assets
    cp /workspace/.gcloudignore /workspace/${strarr[0]}/.gcloudignore 
    #----------------------------------
    #----------DEPLOY FX------------------------
    gcloud functions deploy ${strarr[0]} \
    --project=${project} \
    --region=europe-west1 \
    --source=/workspace/${strarr[0]} \
    --entry-point=${name_camel} \
    --trigger-http \
    --allow-unauthenticated \
    --run-service-account=${sa} \
    --runtime=go121 \
    --env-vars-file ${env}.yaml \
    --timeout=${timeout} \
    --vpc-connector=${vpc} \
    --egress-settings=all \
    ${genFx} 
   
      #----------------------------------
      if [[ "${strarr[1]}" == *"dev"* ]]; then
    echo "dev run services update"
    gcloud run services update ${strarr[0]}  --update-labels tagversion=${tagVersion} --region=europe-west1 --service-account=${sa}
    fi
   
    
     #gcloud functions add-iam-policy-binding ${strarr[0]} \
    # --region='europe-west1' \
    # --member='serviceAccount:wopta-dev-cloudbuild-sa@positive-apex-350507.iam.gserviceaccount.com' \
    # --role='roles/cloudfunctions.invoker'