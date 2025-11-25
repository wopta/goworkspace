
cd ..
echo $1
echo $2
echo $3
git add .

cd $1

cd ..
git add .
git commit -m "$3"
git push work master 
git push google master 
git push wopta master 
git tag -a $1/$2 -m "$3"
git push wopta $1/$2
git push work $1/$2