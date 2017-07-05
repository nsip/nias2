#!/bin/bash

set -e

CWD=`pwd`

echo "Downloading CORE.json"
curl https://raw.githubusercontent.com/nsip/registration-data-set/master/core.json > app/napval/schemas/core.json
curl https://raw.githubusercontent.com/nsip/registration-data-set/master/core_parent2.json > app/napval/schemas/core_parent2.json
curl http://specification.sifassociation.org/Implementation/AU/3.4/XSD/SIF_Message/SIF_Message.xsd > app/sms/SIF_Message.xsd
#echo "Downloading nats-streaming-server"
#go get github.com/nats-io/nats-streaming-server


do_build() {
	cd ../../nats-io/nats-streaming-server
	GOOS="$GOOS" GOARCH="$GOARCH" go build -i -ldflags="$LDFLAGS" -o $OUTPUT/$GNATS
	cd $CWD
}

do_zip() {
        cd $OUTPUT
        cd ..
	rm -f ../$ZIP
        zip -qr ../$ZIP nias
        zip -qr ../$ZIP napval
        zip -qr ../$ZIP naprr
	cd $CWD
}

do_zip_qlserver() {
        cd $OUTPUT
        cd ..
        zip -qr ../$ZIP qlserver
	cd $CWD
}

build_mac64() {
	# MAC OS X (64 only)
	echo "Building Mac binaries..."
	GOOS=darwin
	GOARCH=amd64
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Mac/nias
	GNATS=nats-streaming-server
	ZIP=go-nias-Mac.zip
	do_build
	do_zip
	do_zip_qlserver
	echo "...all Mac binaries built..."
}


build_windows64() {
	# WINDOWS 64
	echo "Building Windows64 binaries..."
	GOOS=windows
	GOARCH=amd64
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Win64/nias
	GNATS=nats-streaming-server.exe
	ZIP=go-nias-Win64.zip
	do_build
	do_zip
	do_zip_qlserver
	echo "...all Windows64 binaries built..."
}

build_windows32() {
	# WINDOWS 32
	echo "Building Windows32 binaries..."
	GOOS=windows
	GOARCH=386
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Win32/nias
	GNATS=nats-streaming-server.exe
	ZIP=go-nias-Win32.zip
	do_build
	do_zip
	echo "...all Windows32 binaries built..."
}

build_linux64() {
	# LINUX 64
	echo "Building Linux64 binaries..."
	GOOS=linux
	GOARCH=amd64
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Linux64/nias
	GNATS=nats-streaming-server
	ZIP=go-nias-Linux64.zip
	do_build
	do_zip
	echo "...all Linux64 binaries built..."
}

build_linux32() {
	# LINUX 32
	echo "Building Linux32 binaries..."
	GOOS=linux
	GOARCH=386
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Linux32/nias
	OUTPUTNAPRR=$CWD/build/Linux32/naprr
	OUTPUTNAPVAL=$CWD/build/Linux32/napval
	GNATS=nats-streaming-server
	NAPVALHARNESS=napval
	SMSHARNESS=sms
	NAPRRHARNESS=naprr
	NAPYR3WHARNESS=napyr3w
	ZIP=go-nias-Linux32.zip
	do_build
	do_zip
	echo "...all Linux32 binaries built..."
}

# TODO ARM
# GOOS=linux GOARCH=arm GOARM=7 go build -o $CWD/build/LinuxArm7/go-nias/aggregator

sh build_sms.sh
sh build_napval.sh
sh build_naprr.sh
sh build_sifql.sh

build_mac64
build_windows64
build_windows32
build_linux64
build_linux32

