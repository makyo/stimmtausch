---
layout: default
title: Configuration
---

## Configuration

Files are read in order from:

* `/etc/stimmtausch/st.yaml`
* `/etc/stimmtausch/conf.d/*`
* `$HOME/.config/stimmtausch/*.st.*`
* `$HOME/.config/stimmtausch/*/*.st.*`
* `$HOME/.strc`

Any values showing up in later files override values in earlier files. See [snuffler](https://github.com/makyo/snuffler) for more on that.

Files can be YAML, TOML, or JSON.

Actual documentation forthcoming, but for now, you can check out [the godocs](https://godoc.org/github.com/makyo/st/config).
