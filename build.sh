#!/bin/bash

set -e

CWD=`pwd`

echo "Building all components..."
# sh build_sms.sh
sh build_napval.sh
sh build_naprrql.sh
sh build_napcomp.sh

echo "Creating zip archives..."
cd $CWD/build
cd Mac
zip -qr ../nias-Mac.zip .
cd ..

cd Win32
zip -qr ../nias-Win32.zip .
cd ..

cd Win64
zip -qr ../nias-Win64.zip .
cd ..

cd Linux32
zip -qr ../nias-Linux32.zip .
cd ..

cd Linux64
zip -qr ../nias-Linux64.zip .
cd ..
echo "Zip archives created"

echo "Removing temporary build files"
rm -r Mac Win32 Win64 Linux32 Linux64

echo "Build Completed."

# build_mac64
# build_windows64
# build_windows32
# build_linux64
# build_linux32

