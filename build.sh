#!/bin/bash

set -e

CWD=`pwd`

echo "Downloading CORE.json"
curl https://raw.githubusercontent.com/nsip/registration-data-set/master/core.json > app/napval/schemas/core.json
curl https://raw.githubusercontent.com/nsip/registration-data-set/master/core_parent2.json > app/napval/schemas/core_parent2.json
curl http://specification.sifassociation.org/Implementation/AU/3.4/XSD/SIF_Message/SIF_Message.xsd > app/sms/SIF_Message.xsd
echo "Downloading nats-streaming-server"
go get github.com/nats-io/nats-streaming-server

do_clear() {
	echo "clear"
	rm -r $OUTPUT/
	rm -r $OUTPUTNAPRR/
}

do_build() {
	mkdir -p $OUTPUT
	mkdir -p $OUTPUTNAPRR
	cd ../../nats-io/nats-streaming-server
	GOOS="$GOOS" GOARCH="$GOARCH" go build -i -ldflags="$LDFLAGS" -o $OUTPUT/$GNATS
	cd $CWD
	cd ./app/naprr
	go get
	GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUTPUTNAPRR/$NAPRRHARNESS
	cd $CWD
	cd ./app/napyr3w
	go get
	GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUTPUTNAPRR/$NAPYR3WHARNESS
	cd $CWD
	cd ./app/napval
	go get 
	GOOS="$GOOS" GOARCH="$GOARCH" go build -i -ldflags="$LDFLAGS" -o $OUTPUT/$NAPVALHARNESS
	cd $CWD
	cd ./app/sms
	go get
	GOOS="$GOOS" GOARCH="$GOARCH" go build -i -ldflags="$LDFLAGS" -o $OUTPUT/$SMSHARNESS
	cd $CWD
	cd ./app
	rsync -a ../test_data napval/nias8help.pdf napval/napval.toml sms/nias.toml napval/napval_nss.cfg sms/nias_nss.cfg napval/public napval/schemas napval/schoolslist napval/templates sms/privacyfilters sms/SIF_Message.xsd $OUTPUT/
	rsync -a ../test_data naprr/templates naprr/public naprr/naprr.toml $OUTPUTNAPRR/
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
	upx $OUTPUT/$NAPVALHARNESS
	upx $OUTPUTNAPRR/$NAPRRHARNESS
	upx $OUTPUTNAPRR/$NAPYR3WHARNESS
	upx $OUTPUT/$SMSHARNESS
}

do_goupx() {
	goupx $OUTPUT/$GNATS
	goupx $OUTPUT/$SMSHARNESS
	goupx $OUTPUT/$NAPVALHARNESS
	goupx $OUTPUTNAPRR/$NAPRRHARNESS
	goupx $OUTPUTNAPRR/$NAPYR3WHARNESS
}

do_zip() {
	cd $OUTPUT
	cd ..
	zip -qr ../$ZIP go-nias8
	cd $CWD
	cd $OUTPUTNAPRR
	cd ..
	zip -qr ../$ZIP naprr
	cd $CWD
}

build_mac64() {
	# MAC OS X (64 only)
	echo "Building Mac binaries..."
	GOOS=darwin
	GOARCH=amd64
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Mac/go-nias8
	OUTPUTNAPRR=$CWD/build/Mac/naprr
	GNATS=nats-streaming-server
	NAPVALHARNESS=napval
	SMSHARNESS=sms
	NAPRRHARNESS=naprr
	NAPYR3WHARNESS=napyr3w
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
	OUTPUTNAPRR=$CWD/build/Win64/naprr
	GNATS=nats-streaming-server.exe
	NAPVALHARNESS=napval.exe
	SMSHARNESS=sms.exe
	NAPRRHARNESS=naprr.exe
	NAPYR3WHARNESS=napyr3w
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
	OUTPUTNAPRR=$CWD/build/Win32/naprr
	GNATS=nats-streaming-server.exe
	NAPVALHARNESS=napval.exe
	SMSHARNESS=sms.exe
	NAPRRHARNESS=naprr.exe
	NAPYR3WHARNESS=napyr3w
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
	OUTPUTNAPRR=$CWD/build/Linux64/naprr
	GNATS=nats-streaming-server
	NAPVALHARNESS=napval
	SMSHARNESS=sms
	NAPRRHARNESS=naprr
	NAPYR3WHARNESS=napyr3w
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
	OUTPUTNAPRR=$CWD/build/Linux32/naprr
	GNATS=nats-streaming-server
	NAPVALHARNESS=napval
	SMSHARNESS=sms
	NAPRRHARNESS=naprr
	NAPYR3WHARNESS=napyr3w
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

