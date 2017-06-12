include $(GOPATH)/src/github.com/digineo/goldflags/goldflags.mk

RELEASE_PATTERN   = bin/ubnt-tools_$(COMMIT_ID)_%.tar.bz2
RELEASE_DIRS     = $(shell find bin/ -mindepth 1 -maxdepth 1 -type d)
RELEASE_ARCHIVES = $(patsubst bin/%,$(RELEASE_PATTERN),$(RELEASE_DIRS))


.PHONY: all
all: discovery provisioner

.PHONY: discovery
discovery:
	cd cmd/ubnt-discovery && $(MAKE)

.PHONY: provisioner
provisioner:
	cd cmd/ubnt-provisioner && $(MAKE)

.PHONY: clean
clean:
	rm -f bin/*.tar.bz2

.PHONY: clobber
clobber: clean
	rm -rf $(RELEASE_DIRS)

.PHONY: release
release: $(RELEASE_ARCHIVES)
	sha256sum $(RELEASE_ARCHIVES)

$(RELEASE_PATTERN): bin/%
	cp resources/config.yml $</provisioner.config.yml
	tar cvjf $@ --directory=$< $(notdir $(wildcard $</*))
