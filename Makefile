SUBDIRS := $(wildcard plugins/*/.)

all:
	for dir in $(SUBDIRS); do \
		cd $$dir && go build && cd ../../; \
	done

install:
	for dir in $(SUBDIRS); do \
		cd $$dir && go install && cd ../../; \
	done

clean:
	for dir in $(SUBDIRS); do \
		cd $$dir && go clean && cd ../../; \
	done

.PHONY: all install clean $(SUBDIRS)
