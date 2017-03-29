#!/bin/bash

rm -rf dist
if [ ! -d "curation" ]; then
  git clone https://github.com/turbotardigrade/agora-curation.git curation
else
  cd curation
  git pull
  cd ..
fi

cd curation
./install.sh
cp -R dist ../