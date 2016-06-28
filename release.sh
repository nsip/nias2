
cd tools; go build release.go; cd ..
./tools/release nias2 go-nias-Mac.zip build/go-nias-Mac.zip
./tools/release nias2 go-nias-Win64.zip build/go-nias-Win64.zip
./tools/release nias2 go-nias-Win32.zip build/go-nias-Win32.zip
./tools/release nias2 go-nias-Linux64.zip build/go-nias-Linux64.zip
./tools/release nias2 go-nias-Linux32.zip build/go-nias-Linux32.zip
./tools/release nias2 go-nias-LinuxArm7.zip build/go-nias-LinuxArm7.zip
