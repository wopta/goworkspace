
cd ..
echo $1
echo $2

git add .
git commit -m "$2"
git push work master 
git push google master  
git push wopta master  
git tag $1 -m "$2"
git push work $1
git push google $1
git push wopta $1

 
