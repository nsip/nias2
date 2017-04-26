
cd tools; go build release.go; cd ..
./tools/release nias2 nias-napval-Mac.zip build/nias-napval-Mac.zip
./tools/release nias2 nias-napval-Win64.zip build/nias-napval-Win64.zip
./tools/release nias2 nias-napval-Win32.zip build/nias-napval-Win32.zip
./tools/release nias2 nias-napval-Linux64.zip build/nias-napval-Linux64.zip
./tools/release nias2 nias-napval-Linux32.zip build/nias-napval-Linux32.zip
