.PHONY: clean
clean:
	rm -rf .bundle .sass-cache _site vendor

.PHONY: deps
deps:
	dep ensure -v
