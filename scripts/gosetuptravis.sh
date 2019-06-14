#!/bin/bash

# Check gofmt
echo "==> Change checkout of code to work with go modules..."

git clone https://github.com/$TRAVIS_REPO_SLUG.git $TRAVIS_REPO_SLUG
cd $TRAVIS_REPO_SLUG
git checkout -qf $TRAVIS_COMMIT