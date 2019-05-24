#!/bin/bash

set -e

CWD=`pwd`

curl http://specification.sifassociation.org/Implementation/AU/3.4/XSD/SIF_Message/SIF_Message.xsd > app/sms/SIF_Message.xsd

do_clear() {
	echo "clear"
	rm -rf $OUTPUT/
}

do_build() {
	echo "Building SMS..."
	mkdir -p $OUTPUT
	cd $CWD
	cd ./app/sms
	go get
	GOOS="$GOOS" GOARCH="$GOARCH" go build -i -ldflags="$LDFLAGS" -o $OUTPUT/$SMSHARNESS
	cd $CWD
	cd ./app
	rsync -a ../test_data sms/nias_nss.cfg sms/privacyfilters sms/SIF_Message.xsd $OUTPUT/
}

do_shells() {
	cd $CWD
	cp bin/gonias.sh $OUTPUT/
	cp bin/stopnias.sh $OUTPUT/
}

do_bats() {
	cd $CWD
	cp bin/gonias.bat $OUTPUT/
	cp bin/stopnias.bat $OUTPUT/
}

do_upx() {
	upx $OUTPUT/$GNATS
	upx $OUTPUT/$SMSHARNESS
}

do_goupx() {
	goupx $OUTPUT/$GNATS
	goupx $OUTPUT/$SMSHARNESS
}

do_zip() {
	cd $OUTPUT
	cd ..
	#zip -qr ../$ZIP go-nias8
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
	SMSHARNESS=nias
	ZIP=go-nias-Mac.zip
	do_clear
	do_build
	#do_upx
	do_shells
	# do_zip
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
	SMSHARNESS=nias.exe
	ZIP=go-nias-Win64.zip
	do_clear
	do_build
	#do_upx
	do_bats
	# do_zip
	echo "...all Windows64 binaries built..."
}

build_linux64() {
	# LINUX 64
	echo "Building Linux64 binaries..."
	GOOS=linux
	GOARCH=amd64
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Linux64/nias
	GNATS=nats-streaming-server
	SMSHARNESS=nias
	ZIP=go-nias-Linux64.zip
	do_clear
	do_build
	#do_goupx
	do_shells
	# do_zip
	echo "...all Linux64 binaries built..."
}

# TODO ARM
# GOOS=linux GOARCH=arm GOARM=7 go build -o $CWD/build/LinuxArm7/go-nias/aggregator

build_mac64
build_windows64
build_linux64
