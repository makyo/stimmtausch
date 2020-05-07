.PHONY: clean
clean:
	rm -rf .bundle .sass-cache _site vendor

.PHONY: deps
deps:
	go build ./...

.PHONY: deb
deb: clean docs
	# Build binary package for GitHub.
	gbp buildpackage --git-ignore-branch  # allow building from any branch
	# Build source package for PPA
	debuild -S -sa

.PHONY: docs
docs: deps
	go run _util/generate-stimmtausch-cmd-docs/main.go
	echo -e "---\nlayout: default\ntitle: \"Command: stimmtausch\"\n---\n\n" > docs/cmd/index.md
	cat docs/cmd/stimmtausch.md >> docs/cmd/index.md
	rm docs/cmd/stimmtausch.md
	echo -e "---\nlayout: default\ntitle: \"Command: stimmtausch headless\"\n---\n\n" > docs/cmd/headless.md
	cat docs/cmd/stimmtausch_headless.md >> docs/cmd/headless.md
	rm docs/cmd/stimmtausch_headless.md
	echo -e "---\nlayout: default\ntitle: \"Command: stimmtausch strip-ansi\"\n---\n\n" > docs/cmd/strip-ansi.md
	cat docs/cmd/stimmtausch_strip-ansi.md >> docs/cmd/strip-ansi.md
	rm docs/cmd/stimmtausch_strip-ansi.md
