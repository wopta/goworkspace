
Rem .\push.bat  v0.0.2 init_rules
cd ..
echo %1
echo %2
echo %3
git add .
git commit -m %2
Rem git tag -a %1 -m %2
Rem git push origin %1
git push origin master 
git push google master  
 
