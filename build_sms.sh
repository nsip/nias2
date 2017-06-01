#!/bin/bash

set -e

CWD=`pwd`

curl http://specification.sifassociation.org/Implementation/AU/3.4/XSD/SIF_Message/SIF_Message.xsd > app/sms/SIF_Message.xsd

do_clear() {
	echo "clear"
	rm -r $OUTPUT/
}

do_build() {
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
	zip -qr ../$ZIP go-nias8
}

build_mac64() {
	# MAC OS X (64 only)
	echo "Building Mac binaries..."
	GOOS=darwin
	GOARCH=amd64
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Mac/go-nias8
	GNATS=nats-streaming-server
	SMSHARNESS=sms
	ZIP=go-nias-Mac.zip
	do_clear
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
	OUTPUT=$CWD/build/Win64/go-nias8
	GNATS=nats-streaming-server.exe
	SMSHARNESS=sms.exe
	ZIP=go-nias-Win64.zip
	do_clear
	do_build
	#do_upx
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
	OUTPUT=$CWD/build/Win32/go-nias8
	GNATS=nats-streaming-server.exe
	SMSHARNESS=sms.exe
	ZIP=go-nias-Win32.zip
	do_clear
	do_build
	#do_upx
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
	OUTPUT=$CWD/build/Linux64/go-nias8
	GNATS=nats-streaming-server
	SMSHARNESS=sms
	ZIP=go-nias-Linux64.zip
	do_clear
	do_build
	#do_goupx
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
	OUTPUT=$CWD/build/Linux32/go-nias8
	GNATS=nats-streaming-server
	SMSHARNESS=sms
	ZIP=go-nias-Linux32.zip
	do_clear
	do_build
	#do_goupx
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

