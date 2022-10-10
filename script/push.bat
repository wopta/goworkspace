
Rem .\push.bat  v0.0.1 init_rules
cd ..
echo %1
echo %2
echo %3
git add .
git commit -m %1
git tag -a %1/%2 -m %3
git push origin %1/%2
git push origin master 
git push google master  
 
