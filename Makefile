.PHONY: discovery
discovery:
	cd cmd/ubnt-discovery && $(MAKE)

.PHONY: clean
	rm -rf bin/
