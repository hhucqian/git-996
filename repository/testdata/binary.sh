#! /bin/sh

rm -rf repo
mkdir -p repo
cd repo

git init
git config user.name test
git config user.email test@test.com

cp /bin/sh ./test.bin
git add test.bin
git commit -m "commit 1"

rm ./test.bin
cp /bin/bash ./test.bin
git add test.bin
git commit -m "commit 2"

rm ./test.bin
git add test.bin
git commit -m "commit 3"

