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
	rm -rf bin/
