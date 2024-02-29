
cd ..
echo $1
echo $2
echo $3
git add .

cd $1

cd ..
git commit -m "$3"
git push work master 
git push google master 
git tag -a $1/$2 -m "$3"
git push google $1/$2