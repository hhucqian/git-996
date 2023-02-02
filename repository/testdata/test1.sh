cd $1

git init
git config user.name test
git config user.email test@test.com

echo "123" > test.txt
git add test.txt
git commit -m "commit 1"

echo "234" > test.txt
git add test.txt
git commit -m "commit 2"

echo "456" > test.txt
echo "567" >> test.txt
git add test.txt
git commit -m "commit 3"

rm test.txt
git add test.txt
git commit -m "commit 4"