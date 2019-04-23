
cd tools; go build release.go; cd ..
./tools/release nias2 > version/version.go
sh build.sh
./tools/release nias2 nias-Mac.zip build/nias-Mac.zip
./tools/release nias2 nias-Win64.zip build/nias-Win64.zip
./tools/release nias2 nias-Linux64.zip build/nias-Linux64.zip
