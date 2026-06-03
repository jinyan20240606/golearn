PKGS := github.com/joyZF/errors
SRCDIRS := $(shell go list -f '{{.Dir}}' $(PKGS))
GO := go

check: test vet gofmt misspell unconvert staticcheck ineffassign unparam

test:
	$(GO) test $(PKGS)

vet: | test
	$(GO) vet $(PKGS)

staticcheck:
	$(GO) get honnef.co/go/tools/cmd/staticcheck
	staticcheck -checks all $(PKGS)

misspell:
	$(GO) get github.com/client9/misspell/cmd/misspell
	misspell \
		-locale GB \
		-error \
		*.md *.go

unconvert:
	$(GO) get github.com/mdempsky/unconvert
	unconvert -v $(PKGS)

ineffassign:
	$(GO) get github.com/gordonklaus/ineffassign
	find $(SRCDIRS) -name '*.go' | xargs ineffassign

pedantic: check errcheck

unparam:
	$(GO) get mvdan.cc/unparam
	unparam ./...

errcheck:
	$(GO) get github.com/kisielk/errcheck
	errcheck $(PKGS)

gofmt:
	@echo Checking code is gofmted
	@test -z "$(shell gofmt -s -l -d -e $(SRCDIRS) | tee /dev/stderr)"


list-tags:
	@echo "Listing Git project history tags:"
	@git tag -l

push-tag:
	@if [ -z "$(TAG)" ]; then \
		echo "Error: Please provide a TAG parameter. Example: make push-tag TAG=v1.0.0"; \
	else \
		echo "Pushing new tag $(TAG) with message: $(MESSAGE)"; \
		git tag -a $(TAG) -m "$(MESSAGE)"; \
		git push origin $(TAG); \
	fi

release: list-tags
	@read -p "Enter the new tag: " TAG; \
	read -p "Enter the tag message: " MESSAGE; \
	make push-tag TAG=$$TAG MESSAGE="$$MESSAGE"