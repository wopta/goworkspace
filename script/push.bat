
Rem .\tag_update.bat rules 0.0.1.dev init_rules
cd ..
echo %1
echo %2
echo %3
git add .
git commit -m %1
git push origin master 
git push google master  
 
