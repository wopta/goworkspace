
Rem .\tag.bat  v0.0.2 init_rules
cd ..
echo %1
echo %2
git add .
git commit -m %2
git push work master 
git push google master  
 git tag %1 -m %2
 git push work %1
 git push google %1

 
