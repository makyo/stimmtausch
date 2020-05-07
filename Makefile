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
	echo "---\nlayout: default\ntitle: \"Command: stimmtausch\"\n---\n\n" > docs/cmd/stimmtausch.md.bak
	cat docs/cmd/stimmtausch.md | sed -e 's/.md)/)/g' | sed -e 's/](st/](\/cmd\/st/g' >> docs/cmd/stimmtausch.md.bak
	mv docs/cmd/stimmtausch.md.bak docs/cmd/stimmtausch.md
	echo "---\nlayout: default\ntitle: \"Command: stimmtausch headless\"\n---\n\n" > docs/cmd/stimmtausch_headless.md.bak
	cat docs/cmd/stimmtausch_headless.md | sed -e 's/.md)/)/g' | sed -e 's/](st/](\/cmd\/st/g' >> docs/cmd/stimmtausch_headless.md.bak
	mv docs/cmd/stimmtausch_headless.md.bak docs/cmd/stimmtausch_headless.md
	echo "---\nlayout: default\ntitle: \"Command: stimmtausch strip-ansi\"\n---\n\n" > docs/cmd/stimmtausch_strip-ansi.md.bak
	cat docs/cmd/stimmtausch_strip-ansi.md | sed -e 's/.md)/)/g' | sed -e 's/](st/](\/cmd\/st/g' >> docs/cmd/stimmtausch_strip-ansi.md.bak
	mv docs/cmd/stimmtausch_strip-ansi.md.bak docs/cmd/stimmtausch_strip-ansi.md
