.PHONY: default configure
default: configure
	pushd coin; go test -v; popd
configure:
	dep ensure
