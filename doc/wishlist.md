---
layout: default
title: Wishlist
---

## Wishlist

* Global config with common settings/servers (viper can handle it)
* Config editor in the UI
    * Create server/serverType/world entries in views, fiddle with client settings
    * Write to config file
    * Reload config
* Menu bar with items for each command (or at least most of the common ones). Should be disable-able, or at least hideable
* Something for max post size
    * At least a post size character count
    * Maybe highlight characters over limit
    * Break it up at post limit, with command to run before each successive post?
        * Break at one word before limit
        * Use command prefix to send next chunk
        * e.g: `post_limit_chunk_prefix: "spoof ...%s"` or something?
* On that note, means of staging posts (ctrl+enter?)
* List of open logfiles per world
* Timestamps (both ui and logging)
* Connection tiling
* Using ui only mode as a front end to headless mode
    * Involves looser coupling between UI and `client`. Maybe a socket?
* MCP support
