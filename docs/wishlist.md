---
layout: default
title: Wishlist
---

## Wishlist

* ~~Global config with common settings/servers (viper can handle it)~~ Viper couldn't handle it, but snuffler can.
* Config editor in the UI
    * Create server/serverType/world entries in views, fiddle with client settings
    * Write to config file
    * Reload config
* Menu bar with items for each command (or at least most of the common ones). Should be disable-able, or at least hideable³
* Something for max post size
    * At least a post size character count
    * Maybe highlight characters over limit
    * Break it up at post limit, with command to run before each successive post?
        * Break at one word before limit
        * Use command prefix to send next chun
        * e.g: `post_limit_chunk_prefix: "spoof ...%s"` or something?
* On that note, means of staging posts (ctrl+enter?)
* Similarly, ability for input buffer to optionally expand with long posts.
* List of open logfiles per world
* Timestamps (both ui and logging)
* Connection tiling
* Using ui only mode as a front end to headless mode
    * Involves looser coupling between UI and `client`. Maybe a socket?
* MCP support
* Triggers¹
    * ~~Hilite~~
    * ~~Gag~~
    * Script
    * Macro
* ~~Allow multiple matches per trigger~~ (plus case-sensitivity?)⁵
* Macros
    * [Zygo?](https://github.com/glycerine/zygomys)² - either way, use a predefined embedded language for better docs early on
    * Use `/<macro>` to call that macro
    * Have a standard library in `/etc/stimmtausch/macro`
    * Predefined macros either via stdlib or in go land:
        * `new-trigger` creates a new trigger (e.g: `wf-partial` for creating new temporary triggers
        * `call` calls an existing macro
        * `connections` gets a list of current connections (struct of name, host, etc?)
        * `connect` adds connection
        * `go-to-connection`/`fg` go to world (optionally `fg "->"`/`fg "<-"` á là TF)


## Notes

### ^1

```yaml
triggers:
    - type: hilite
      match: "[Mm]addy"
      attributes: [green, bold]
    - type: gag
      match: "^## Saving changed objects, just a moment... ##$"
    - type: script
      # Matches are passed in as arguments, e.g:
      #   Maddy mew, "!MaddyLog: asdfasdf!"
      # calls script with arguments:
      #   !MaddyLog: asdfasdf!
      #   asdfasdf
      match: "!MaddyLog: (.+)!"
      # Script is any executable
      script: ~/.config/stimmtausch/scripts/foo.py
    - type: macro
      # As above with arguments
      match: "!MaddyMacro: (.+)!"
      macro: do-a-thing
```

### 2

Looks like Zygo only returns the top item of the stack, which means we can't just add arbitrary things to the stack to have them show up. Solution is to maybe have a `register-macro` function:

```zygo
(defn maddy-log [match0, match1]
  (...))

(register-macro maddy-log)
```

Still, maybe Zygo's not quite what we want. More research required.

### 3

Use fn keys. (F1=help, F2=file, etc...)

* if vim mode⁴:
    * Show only in normal mode
* else:
    * Show always

### 4

* Configurable
* Separately: Alt mode
    * Alt + → : Next world
    * Alt + ← : Prev world

### 5

```
triggers:
    # My characters
    - type: hilite
      case-sensitive: false
      matches:
        - Maddy
        - Makyo
        - Younes
        - Happenstance
