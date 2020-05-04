.PHONY: clean
clean:
	rm -rf .bundle .sass-cache _site vendor

.PHONY: deps
deps:
	go build ./...

.PHONY: deb
deb: clean deps
	# Build binary package for GitHub.
	gbp buildpackage
	# Build source package for PPA
	debuild -S -sa
