#!/bin/bash

set -e

CWD=`pwd`

do_build() {
	echo "Building Documentation..."
	mkdir -p $OUTPUT
	rsync -a app/documentation/*.pdf $OUTPUT/
}

build_mac64() {
	# MAC OS X (64 only)
	OUTPUT=$CWD/build/Mac/documentation
	do_build
}


build_windows64() {
	OUTPUT=$CWD/build/Win64/documentation
	do_build
}

build_linux64() {
	OUTPUT=$CWD/build/Linux64/documentation
	do_build
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
    build_mac64
    build_windows64
    build_linux64
fi

