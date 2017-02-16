
cd tools; go build release.go; cd ..
./tools/release nias2 nias-naprr-Mac.zip build/nias-naprr-Mac.zip
./tools/release nias2 nias-naprr-Win64.zip build/nias-naprr-Win64.zip
./tools/release nias2 nias-naprr-Win32.zip build/nias-naprr-Win32.zip
./tools/release nias2 nias-naprr-Linux64.zip build/nias-naprr-Linux64.zip
./tools/release nias2 nias-naprr-Linux32.zip build/nias-naprr-Linux32.zip
