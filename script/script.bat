cd ..
echo $1
echo $2
git add .
git commit -m $2
git tag -a enrich-vat/$1 -m $2
git push google enrich-vat/$1