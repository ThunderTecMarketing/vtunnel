#!/usr/bin/env bash

TARGET="vtunnel_darwin"

cd ./main/
gox -osarch darwin/amd64 -output $TARGET
cd ..
mv ./main/$TARGET .

