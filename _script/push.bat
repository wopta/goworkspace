
Rem .\push.bat "init_rules"
cd ..
echo %1

git add .
git commit -m %1
git push work master 
git push google master  
 
