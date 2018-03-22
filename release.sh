
cd tools; go build release.go; cd ..
./tools/release > version/version.go
sh build.sh
./tools/release nias2 nias-Mac.zip build/nias-Mac.zip
./tools/release nias2 nias-Win64.zip build/nias-Win64.zip
./tools/release nias2 nias-Win32.zip build/nias-Win32.zip
./tools/release nias2 nias-Linux64.zip build/nias-Linux64.zip
./tools/release nias2 nias-Linux32.zip build/nias-Linux32.zip
