
Rem .\tag_update.bat rules 0.0.1.dev init_rules  .\tag_update.bat quoteAllrisk 0.0.1.dev init_quote_munich  .\tag_update.bat enrich-vat 0.0.20.dev fix _cors  
cd ..
echo %1
echo %2
echo %3

cd %1

cd ..
git add .
git commit -m %3
git push origin master 
git push google master  
git tag -a %1/%2 -m %3
git push google %1/%2