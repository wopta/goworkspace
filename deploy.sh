export namefx=( $(grep -Eo '[[:digit:]]+|[^[:digit:]]+' <<<'$TAG_NAME') )
    #Print current tag
    echo "current tag: ${TAG_NAME}"
  
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
    namefx=$strarr[0]

    name_camel=$(sed -r 's/(^|-)(\w)/\U\2/g' <<<"${strarr[0]}")
    tagVersion=$(echo "${strarr[1]}" | sed -r 's/\./_/g')
    t=tr -s . _ <<< "${strarr[1]}"
    echo "namefx camel : ${name_camel}"
    echo "tagVersion camel : ${tagVersion}"
    echo "tagVersion camel : ${t}"

    if [[ "${strarr[1]}" == *"dev"* ]]; then
    echo "dev enviroment"
      bucket=gs://function-data
      project=positive-apex-350507
      env=dev
      genFx=--gen2
      sa=wopta-dev-cloudbuild-sa@positive-apex-350507.iam.gserviceaccount.com
      timeout=60
    fi

    if [[ "${strarr[1]}" == *"prod"* ]]; then
    echo "prod enviroment"
      bucket=gs://core-350507-function-data
      project=core-350507
      env=prod
      genFx=""
      sa=wopta-prod-cloud-function@core-350507.iam.gserviceaccount.com
      timeout=520
    fi
        #--ingress-settings internal-only \
    #--timeout=540 \
    # Save a value to persistent volume mount: "/workspace"
    echo $namefx > /workspace/namefx.txt &&
    # Save another
    echo $name_camel > /workspace/entrypoint.txt
    #ls -la
    
    # Copy assets folder from Google Bucket to directory
    mkdir -p /workspace/${strarr[0]}/tmp/assets
    if [[ "${strarr[1]}" == *"dev"* ]]; then

    echo "dev enviroment"
      gsutil -m cp -r gs://function-data/assets/documents/** /workspace/${strarr[0]}/tmp/assets
    fi
     if [[ "${strarr[1]}" == *"prod"* ]]; then
    echo "prod enviroment"
      gsutil -m cp -r gs://core-350507-function-data/assets/documents/** /workspace/${strarr[0]}/tmp/assets
    fi
    
    cp /workspace/.gcloudignore /workspace/${strarr[0]}/.gcloudignore 
    
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
    --timeout=${timeout
    ${genFx} 
   
   
      if [[ "${strarr[1]}" == *"dev"* ]]; then
    echo "dev run services update"
    gcloud run services update ${strarr[0]}  --update-labels tagversion=${tagVersion} --region=europe-west1 --service-account=${sa}
    fi
   
    
     #gcloud functions add-iam-policy-binding ${strarr[0]} \
    # --region='europe-west1' \
    # --member='serviceAccount:wopta-dev-cloudbuild-sa@positive-apex-350507.iam.gserviceaccount.com' \
    # --role='roles/cloudfunctions.invoker'