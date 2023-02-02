cd $1

git init
git config user.name test
git config user.email test@test.com

echo "123" > test.txt
git add test.txt
git commit -m "commit 1"

cp /bin/sh ./test.bin
git add test.bin
git commit -m "commit 2"