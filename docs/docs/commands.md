---
layout: default
title: Commands
---

## Commands

As with TinyFugue, one communicates with Stimmtausch via commands starting with `/`. When you send a line that starts with `/`, the client first searches for a builtin by that name, then searches for a macro by that name.

### Builtins

`/connect [connectStr]`, `/c [connectStr]`
:   Connect to the specified world/server/connection string. The argument is required

`/disconnect [-r] [connectionName]`, `/dc [-r] [connectionName]`
:   Disconnect from the specified connection. If no world is provided, it disconnects from the current world. If `-r` is provided, it also removes the world from the UI.

`/remove [world]`, `/r [world]`
:   Remove the current world from the UI. **Warning:** this does not disconnect the world, and there is no way to re-attach a world to the UI yet, so use with care. Provided mostly for removing stale, disconnected worlds from the UI.

`/fg [directionOrConnectionName]`
:   Switch to the given world. If called as `/fg <` or `/fg >`, it switches one world in the given direction. As a shortcut for those, you can also use `/<` or `/>`. `/fg` without a direction or connection is equivalent to `/]` below.

`/]` and `/[`
:   Rotate to the next active world in that direction. For example, `/]` keeps calling `/>` until it hits a world with more lines (stopping at the current world if it doesn't find it).

`/quit`
:   Disconnects from all worlds and quits the program.

`/syslog [level] [message]`
:   Simple test command that logs a message to the system log at the given level (which can be `trace`, `debug`, `info`, `warning`, `error`, `critical`).`

## Key bindings

`Ctrl+C`
:   Quit Stimmtausch

`Ctrl+L`
:   Redraw the screen

`Esc+→` and `Esc+←`
:   Rotate one world to the right or left (equivalent to `/>` and `/<`)

`Ctrl+]` and `Ctrl+[`
:   Rotate to the next active world to the right or left (equivalent to `/]` and `/[`)
