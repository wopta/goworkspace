
Rem .\push.bat  v0.0.2 init_rules
set location=callback
set v=1.0.0.prod
set m="major allineament"


git add .
git commit -m %m%
git tag -a document/%v% -m %m%
git push google document/%v%


git tag -a rules/%v% -m %m%
git push google rules/%v%


git tag -a broker/%v% -m %m%
git push google broker/%v%


git tag -a claim/%v% -m %m%
git push google claim/%v%

git tag -a enrich/%v% -m %m%
git push google enrich/%v%

git tag -a payment/%v% -m %m%
git push google payment/%v%


git tag -a product/%v% -m %m%
git push google product/%v%


git tag -a %location%/%v% -m %m%
git push google %location%/%v%


git tag -a quote/%v% -m %m%
git push google quote/%v%

git tag -a user/%v% -m %m%
git push google user/%v%

git tag -a mail/%v% -m %m%
git push google mail/%v%