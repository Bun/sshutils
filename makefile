.PHONY: sshutils
sshutils:
	go install -v ./cmd/pssh ./cmd/prun

test:
	go test ./...

.PHONY: deb
deb:
	rm -rf build/
	mkdir -p build/DEBIAN
	mkdir -p build/usr/bin
	cp debian/control build/DEBIAN/control
	go build -v -o build/usr/bin/pssh ./cmd/pssh
	go build -v -o build/usr/bin/prun ./cmd/prun
	fakeroot dpkg-deb -z2 --build build/ sshutils-0.1.deb
