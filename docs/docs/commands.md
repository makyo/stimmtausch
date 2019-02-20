---
layout: default
title: Commands
---

## Commands

As with TinyFugue, one communicates with Stimmtausch via commands starting with `/`. When you send a line that starts with `/`, the client first searches for a builtin by that name, then searches for a macro by that name.

### Builtins

`/connect [connectStr]`, `/c [connectStr]`
:   Connect to the specified world/server/connection string. The argument is required

`/disconnect [connectionName]`, `/dc [connectionName]`
:   Disconnect from the specified connection. If no argument is provided, it disconnects from the current world.

`/fg [directionOrConnectionName]`
:   Switch to the given world. If called as `/fg <` or `/fg >`, it switches one world in the given direction. As a shortcut for those, you can also use `/<` or `/>`

`/quit`
:   Disconnects from all worlds and quits the program.

`/syslog [level] [message]`
:   Simple test command that logs a message to the system log at the given level (which can be `trace`, `debug`, `info`, `warning`, `error`, `critical`).`
