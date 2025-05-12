CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
LINTVER=v1.62.2
LINTBIN=golangci-lint


bindir:
	mkdir -p ${BINDIR}


install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

define lint
	@if [ -f "$(1)/go.mod" ]; then \
		output=$$(${LINTBIN} --config=.golangci.yaml run $(1) 2>&1); \
		exit_code=$$?; \
		echo "$$output"; \
		if [ $$exit_code -ne 0 ]; then \
			if echo "$$output" | grep -q "no go files to analyze"; then \
				exit 0; \
			else \
				exit $$exit_code; \
			fi \
		fi \
	fi
endef


cart-lint:
	$(call lint,cart)

loms-lint:
	$(call lint,loms)

notifier-lint:
	$(call lint,notifier)

comments-lint:
	$(call lint,comments)
