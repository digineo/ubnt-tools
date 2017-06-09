include $(GOPATH)/src/github.com/digineo/goldflags/goldflags.mk

NAME      = ubnt-dicovery
TARGET    = ../../bin/$(NAME)
DEPS      = $(shell find ../../discovery -type f -name '*.go')
DEPS     += $(shell find . -type f -name '*.go')
LDFLAGS	  = -s -w $(GOLDFLAGS)

$(NAME): $(DEPS)
	mkdir -p $(dir $TARGET)
	go build -ldflags "$(LDFLAGS)" -o "$(TARGET)"
