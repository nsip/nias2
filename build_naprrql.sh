#!/bin/bash

set -e

CWD=`pwd`


    do_build() {
        # comment out line below to exclude naprrqlhp from build
        include_hp
        echo "Building NAPRRQL..."
        mkdir -p $OUTPUT
    	cd $CWD
    	cd ./app/naprrql
    	go get
    	GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUTPUT/$HARNESS
        cd $CWD
        cd ./app
        mkdir -p naprrql/in
        rm -rf $OUTPUT/in/*.xml $OUTPUT/in/*.zip $OUTPUT/kvs
        rsync -a naprrql/naprrql.toml naprrql/gql_schemas naprrql/in naprrql/public naprrql/reporting_templates naprrql/*.pdf $OUTPUT/
    }


    include_hp() {
        echo "Including NAPRRQLHP..."
        cd $CWD
        cd ./app/naprrqlhp
        go get
        GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUTPUT/$HPHARNESS
    }



      do_zip() {
        cd $OUTPUT
        cd ..
        zip -qr ../$ZIP naprrql
        cd $CWD
      }

      build_mac64() {
        # MAC OS X (64 only)
        echo "Building Mac binaries..."
        GOOS=darwin
        GOARCH=amd64
        LDFLAGS="-s -w"
        OUTPUT=$CWD/build/Mac/naprrql
        # GNATS=nats-streaming-server
        HARNESS=naprrql
        AUDITDIFFHARNESS=napcomp
        HPHARNESS=naprrqlhp
        ZIP=naprrql-Mac.zip
        do_build
        #do_upx
        # do_shells
        # do_zip
        echo "...all Mac binaries built..."
      }


      build_windows64() {
        # WINDOWS 64
        echo "Building Windows64 binaries..."
        GOOS=windows
        GOARCH=amd64
        LDFLAGS="-s -w"
        OUTPUT=$CWD/build/Win64/naprrql
        # GNATS=nats-streaming-server.exe
        HARNESS=naprrql.exe
        HPHARNESS=naprrqlhp.exe
        AUDITDIFFHARNESS=napcomp.exe
        ZIP=naprrql-Win64.zip
        do_build
        #do_upx
        # do_bats
        # do_zip
        echo "...all Windows64 binaries built..."
      }

      build_linux64() {
        # LINUX 64
        echo "Building Linux64 binaries..."
        GOOS=linux
        GOARCH=amd64
        LDFLAGS="-s -w"
        OUTPUT=$CWD/build/Linux64/naprrql
        # GNATS=nats-streaming-server
        HARNESS=naprrql
        HPHARNESS=naprrqlhp
        AUDITDIFFHARNESS=napcomp
        ZIP=naprrql-Linux64.zip
        do_build
        #do_goupx
        # do_shells
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
