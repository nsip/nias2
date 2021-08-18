#!/bin/bash

set -e

CWD=`pwd`

#rm -rf dev-nrt
#git clone https://github.com/nsip/dev-nrt

# Linux64
build_linux64() {
    
  echo "Building Linux64 binaries..."
  cd ./dev-nrt
  rm -rf ./build && ./build.sh Linux64

  OUTPUT="$CWD"/build/Linux64/naprrql
  rm -rf "$OUTPUT"
  mkdir -p "$OUTPUT"
  cp -r ./build/* "$OUTPUT"
  
  ./clean.sh

  cd "$OUTPUT"
  mv testdata in
  HARNESS=naprrql
  mv nrt $HARNESS  

  # zip
  ZIP=naprrql-Linux64.zip
  cd ..
  #zip -qr ../$ZIP naprrql   

  cd "$CWD"
  echo "...all Linux64 binaries built..." 

}

# Mac
build_mac64() {
    
  echo "Building Mac binaries..."
  cd ./dev-nrt
  rm -rf ./build && ./build.sh Mac

  OUTPUT="$CWD"/build/Mac/naprrql
  rm -rf "$OUTPUT"
  mkdir -p "$OUTPUT"
  cp -r ./build/* "$OUTPUT"

  ./clean.sh
  
  cd "$OUTPUT"
  mv testdata in
  HARNESS=naprrql
  mv nrt $HARNESS

  # zip
  ZIP=naprrql-Mac.zip
  cd ..
  #zip -qr ../$ZIP naprrql   
  
  cd "$CWD"
  echo "...all Mac binaries built..." 
  
}

# Windows64
build_windows64() {
    
  echo "Building Windows64 binaries..."
  cd ./dev-nrt
  rm -rf ./build && ./build.sh Windows64

  OUTPUT="$CWD"/build/Win64/naprrql
  rm -rf "$OUTPUT"
  mkdir -p "$OUTPUT"
  cp -r ./build/* "$OUTPUT"

  ./clean.sh
  
  cd "$OUTPUT"
  mv testdata in
  HARNESS=naprrql.exe
  mv nrt.exe $HARNESS

  # zip
  ZIP=naprrql-Windows64.zip
  cd ..
  #zip -qr ../$ZIP naprrql
  
  cd "$CWD"
  echo "...all Windows64 binaries built..." 
  
}

if [ "$1" = "L64"  ]
then
    build_linux64
elif [ "$1" = "W64"  ]
then
    build_windows64
elif [ "$1" = "M64"  ]
then
    build_mac64
else
    build_linux64
    build_mac64
    build_windows64    
fi

cd "$CWD"

# rm -rf ./dev-nrt
