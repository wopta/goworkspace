
cd ..
echo $1

git add .
git commit -m "$1"
git push origin master 
git push google master  
 