DEPS      = $(shell find src -type f)
DEPS     += index.html
GOBINDATA = $(GOPATH)/bin/go-bindata

ASSETS    = \
	dist/views/index.html \
	dist/assets/application.js

bindata.go: $(GOBINDATA) $(ASSETS)
	$(GOBINDATA) -pkg ui -prefix dist -o $@ dist/...

$(GOBINDATA):
	go get -u github.com/jteeuwen/go-bindata/...

dist/views/index.html: index.html
	mkdir -p ../../provisioner/ui/dist/views
	cp index.html dist/views/

dist/assets/application.js: $(DEPS)
	yarn install
	npm run build

.PHONY: clean
clean:
	rm -rf dist
	rm bindata.go

.PHONY: clobber
clobber: clean
	rm -rf node_modules
