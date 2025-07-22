#!/bin/bash

set -e

CWD=`pwd`

echo "Downloading CORE.json"
curl https://raw.githubusercontent.com/nsip/registration-data-set/master/core.json > app/napval/schemas/core.json
curl https://raw.githubusercontent.com/nsip/registration-data-set/master/core_parent2.json > app/napval/schemas/core_parent2.json
echo "Downloading nats-streaming-server"
# go get github.com/nats-io/nats-streaming-server
# go get github.com/nats-io/stan.go
rm -rf nats-streaming-server
git clone https://github.com/nats-io/nats-streaming-server
cd nats-streaming-server
cd ..

do_build() {
	echo "Building NAPVAL..."
	mkdir -p $OUTPUT
        rm -rf $OUTPUT/*.csv
	cd nats-streaming-server
	# cd ../../nats-io/stan.go
	GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUTPUT/$GNATS -mod=mod
	cd $CWD
	cd ./app/napval
	go get
	GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUTPUT/$HARNESS -mod=mod
	cd ..
	rsync -a napval/napval.toml napval/napval_nss.cfg napval/public napval/schemas napval/schoolslist napval/templates $OUTPUT/
	rsync -a napval/students.csv $OUTPUT/
}

do_shells() {
	cd $CWD
	cp bin/gonapval.sh $OUTPUT/
	cp bin/stopnapval.sh $OUTPUT/
}

do_bats() {
	cd $CWD
	cp bin/gonapval.bat $OUTPUT/
	cp bin/stopnapval.bat $OUTPUT/
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
	#zip -qr ../$ZIP napval
	cd $CWD
}

build_mac64() {
	# MAC OS X (64 only)
	echo "Building Mac binaries..."
	GOOS=darwin
	GOARCH=amd64
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Mac/napval
	GNATS=nats-streaming-server
	HARNESS=napval
	ZIP=nias-napval-Mac.zip
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
	OUTPUT=$CWD/build/Win64/napval
	GNATS=nats-streaming-server.exe
	HARNESS=napval.exe
	ZIP=nias-napval-Win64.zip
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
	OUTPUT=$CWD/build/Linux64/napval
	GNATS=nats-streaming-server
	HARNESS=napval
	ZIP=nias-napval-Linux64.zip
	do_build
	#do_goupx
	do_shells
	# do_zip
	echo "...all Linux64 binaries built..."
}

# TODO ARM
# GOOS=linux GOARCH=arm GOARM=7 go build -o $CWD/build/LinuxArm7/go-nias/aggregator

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
    build_mac64
    build_windows64
    build_linux64
fi
