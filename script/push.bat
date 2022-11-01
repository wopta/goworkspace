
Rem .\push.bat  v0.0.2 init_rules
cd ..
echo %1

git add .
git commit -m %1
git push origin master 
git push google master  
 
