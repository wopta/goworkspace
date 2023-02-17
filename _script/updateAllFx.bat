
Rem .\push.bat  v0.0.2 init_rules
set location=callback
set v="0.5.4.dev"
set m="major allineament"


git add .
git commit -m %m%
git tag -a document/%v% -m %m%
git push google document/%v%


git add .
git commit -m %m%
git tag -a rules/%v% -m %m%
git push google rules/%v%


git add .
git commit -m %m%
git tag -a broker/%v% -m %m%
git push google broker/%v%


git add .
git commit -m %m%
git tag -a claim/%v% -m %m%
git push google claim/%v%

git add .
git commit -m %m%
git tag -a enrich/%v% -m %m%
git push google enrich/%v%

git add .
git commit -m %m%
git tag -a payment/%v% -m %m%
git push google payment/%v%


git add .
git commit -m %m%
git tag -a product/%v% -m %m%
git push google product/%v%


 git add .
 git commit -m %m%
git tag -a %location%/%v% -m %m%
git push google %location%/%v%

git add .
git commit -m %m%
git tag -a quote/%v% -m %m%
git push google quote/%v%