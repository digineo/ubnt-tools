include $(GOPATH)/src/github.com/digineo/goldflags/goldflags.mk

NAME      = ubnt-dicovery
TARGET    = ../../bin/$(NAME)
DEPS      = $(shell find ../../discovery -type f -name '*.go')
DEPS     += $(shell find . -type f -name '*.go')

$(NAME): $(DEPS)
	$(call goldbuild,$(TARGET),linux,amd64)
	$(call goldbuild,$(TARGET),darwin,amd64)
	$(call goldbuild,$(TARGET),windows,amd64)
