.PHONY: clean
clean:
	rm -rf .bundle .sass-cache _site vendor

.PHONY: deps
deps:
	go build ./...

.PHONY: package
package: clean docs
	# Build binary package for GitHub.
	gbp buildpackage --git-ignore-branch  # allow building from any branch
	# Build source package for PPA
	debuild -S -sa
	# Build snap
	snapcraft

.PHONY: docs
docs: deps
	go run _util/generate-stimmtausch-cmd-docs/main.go
	echo "---\nlayout: default\ntitle: \"Command: stimmtausch\"\n---\n\n" > docs/cmd/index.md
	cat docs/cmd/stimmtausch.md | sed -e 's/.md)/)/g' | sed -e 's/](st/](\/cmd\/st/g' >> docs/cmd/index.md
	mv docs/cmd/index.md docs/cmd/stimmtausch.md
	echo "---\nlayout: default\ntitle: \"Command: stimmtausch headless\"\n---\n\n" > docs/cmd/headless.md
	cat docs/cmd/stimmtausch_headless.md | sed -e 's/.md)/)/g' | sed -e 's/](st/](\/cmd\/st/g' >> docs/cmd/headless.md
	mv docs/cmd/headless.md docs/cmd/stimmtausch_headless.md
	echo "---\nlayout: default\ntitle: \"Command: stimmtausch strip-ansi\"\n---\n\n" > docs/cmd/strip-ansi.md
	cat docs/cmd/stimmtausch_strip-ansi.md | sed -e 's/.md)/)/g' | sed -e 's/](st/](\/cmd\/st/g' >> docs/cmd/strip-ansi.md
	mv docs/cmd/strip-ansi.md docs/cmd/stimmtausch_strip-ansi.md
