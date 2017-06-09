include $(GOPATH)/src/github.com/digineo/goldflags/goldflags.mk

NAME        = ubnt-provisioner
TARGET      = ../../bin/$(NAME)
DEPS        = $(shell find ../../discovery -type f -name '*.go')
DEPS       += $(shell find ../../provisioner -type f -name '*.go')
DEPS       += $(shell find . -type f -name '*.go')

ASSETS      = ../../provisioner/ui/bindata.go
INDEX       = ../../provisioner/ui/dist/views/index.html

ASSET_DEPS  = $(shell find ../../provisioner/ui/src -type f)
ASSET_DEPS += $(INDEX)

.PHONY: $(NAME)
$(NAME): $(DEPS) $(ASSETS)
	$(call goldbuild,$(TARGET),linux,amd64)
	$(call goldbuild,$(TARGET),darwin,amd64)
	$(call goldbuild,$(TARGET),windows,amd64)

$(ASSETS): $(ASSET_DEPS)
	cd ../../provisioner/ui && $(MAKE)
