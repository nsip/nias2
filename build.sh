#!/bin/bash

set -e

CWD=`pwd`

echo "Downloading CORE.json"
curl https://raw.githubusercontent.com/nsip/registration-data-set/master/core.json > harness/schemas/core.json
curl https://raw.githubusercontent.com/nsip/registration-data-set/master/core_parent2.json > harness/schemas/core_parent2.json
echo "Downloading gnatsd"
go get github.com/nats-io/gnatsd

do_build() {
	mkdir -p $OUTPUT
	cd ../../nats-io/gnatsd
	GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUTPUT/$GNATS
	cd $CWD
	cd ./harness
	go get
	GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUTPUT/$HARNESS
	cd ..
	rsync -a harness/nias.toml harness/public harness/schemas harness/schoolslist harness/templates harness/privacyfilters $OUTPUT/
}

do_shells() {
	cp bin/gonias.sh $OUTPUT/
	cp bin/stopnias.sh $OUTPUT/
}

do_bats() {
	cp bin/gonias.bat $OUTPUT/
	cp bin/stopnias.bat $OUTPUT/
}

do_upx() {
	upx $OUTPUT/$GNATS
	upx $OUTPUT/$HARNESS
}

do_goupx() {
	goupx $OUTPUT/$GNATS
	goupx $OUTPUT/$HARNESS
}

do_zip() {
	cd $OUTPUT
	cd ..
	zip -qr ../$ZIP go-nias
	cd $CWD
}

build_mac64() {
	# MAC OS X (64 only)
	echo "Building Mac binaries..."
	GOOS=darwin
	GOARCH=amd64
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Mac/go-nias
	GNATS=gnatsd
	HARNESS=harness
	ZIP=go-nias-Mac.zip
	do_build
	#do_upx
	do_shells
	do_zip
	echo "...all Mac binaries built..."
}


build_windows64() {
	# WINDOWS 64
	echo "Building Windows64 binaries..."
	GOOS=windows
	GOARCH=amd64
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Win64/go-nias
	GNATS=gnatsd.exe
	HARNESS=harness.exe
	ZIP=go-nias-Win64.zip
	do_build
	do_upx
	do_bats
	do_zip
	echo "...all Windows64 binaries built..."
}

build_windows32() {
	# WINDOWS 32
	echo "Building Windows32 binaries..."
	GOOS=windows
	GOARCH=386
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Win32/go-nias
	GNATS=gnatsd.exe
	HARNESS=harness.exe
	ZIP=go-nias-Win32.zip
	do_build
	do_upx
	do_bats
	do_zip
	echo "...all Windows32 binaries built..."
}

build_linux64() {
	# LINUX 64
	echo "Building Linux64 binaries..."
	GOOS=linux
	GOARCH=amd64
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Linux64/go-nias
	GNATS=gnatsd
	HARNESS=harness
	ZIP=go-nias-Linux64.zip
	do_build
	do_goupx
	do_shells
	do_zip
	echo "...all Linux64 binaries built..."
}

build_linux32() {
	# LINUX 32
	echo "Building Linux32 binaries..."
	GOOS=linux
	GOARCH=386
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Linux32/go-nias
	GNATS=gnatsd
	HARNESS=harness
	ZIP=go-nias-Linux32.zip
	do_build
	do_goupx
	do_shells
	do_zip
	echo "...all Linux32 binaries built..."
}

# TODO ARM
# GOOS=linux GOARCH=arm GOARM=7 go build -o $CWD/build/LinuxArm7/go-nias/aggregator

build_mac64
build_windows64
build_windows32
build_linux64
build_linux32

