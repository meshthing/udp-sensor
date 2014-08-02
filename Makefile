build:
	scripts/build.sh

clean:
	rm -f bin/udp-sensor || true
	rm -rf .gopath || true

.PHONY: build clean