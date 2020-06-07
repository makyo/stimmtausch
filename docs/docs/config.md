---
layout: default
title: Configuration
---

## Configuration

*TL;DR:* There are lots of configuration options, and most of them are set for you except for [worlds](#worlds) and [triggers](#triggers), which you'll probably want to put in your local configuration.

Stimmtausch is designed to be configurable from the ground up, as much as possible. It relies on configuration sources in one of three formats: YAML, TOML, or JSON. The default configuration sets up the very basics, and from there, it starts loading every file in `$HOME/.config/stimmtausch` whose name includes the string `.st.` - `worlds.st.yaml`, `triggers.st.yaml`, and so on.

Any values showing up in later files *override* values in earlier files. See [snuffler](https://github.com/makyo/snuffler) for more on that. You can see the effective configuration by running `stimmtausch config`, which will build the internal configuration that the client uses and then dump it back out to the screen. To see all of the builtin defaults by themselves, you can run `stimmtausch config --default`.

Configuration settings are broken down into a few categories:

* [Server types](#server-types), which basically boil down to the strings used to connect/disconnect from a server.
* [Servers](#servers), which are the addresses, ports and other such information for MU\*s.
* [Worlds](#worlds), which are how you log in - they associate usernames and passwords with servers.
* [Triggers](#triggers), which cover things that Stimmtausch should automatically do when something happens in the world, such as highlight a word or run a script.
* [Client](#client), which holds information about the Stimmtausch client itself. This is further broken down into a few categories:
    * [Syslog](#syslog), which details what to do with the logs that Stimmtausch itself generates.
    * [Logging](#logging), which holds information about logging from the connections.
    * [UI](#ui), which describes various bits of the user interface

All configuration files must have a top-level `stimmtausch` key. This helps the configuration software know that it's actually loading stuff for Stimmtausch and not some random program.

### Server types

About
:   `server_types` hold information about how to automatically connect to or disconnect from a world.

    Expects an object mapping server type names to values.

Values
:  
    * `connect_string` (*string*) - the string to send to the server when used to connect. The substrings `$username` and `$password` will be replaced with those values.

      Example: `connect_string: "connect $username $password"`

    * `disconnect_string` (*string*) - the string to send to the server when disconnecting.

      Example: `disconnect_string: QUIT`

**Default**


```yaml
stimmtausch:
    server_types:
        muck:
            name: "TinyMUCK, FuzzballMUCK, etc."
            connect_string: "connect $username $password"
            disconnect_string: "QUIT"
```

### Servers

About
:   `servers` holds connection details (minus username/password) for MU\*s.

    Expects an object mapping server names to values.

Values
:  
    * `host` (*string* required) - the domain name or IP for the MU\*".

      Example: `host: furrymuck.com`

    * `port` (*number* required) - the port for the MU\*.

      Example: `port: 8899`

    * `ssl` (*boolean*) - whether or not to use SSL when connecting.

      Example: `ssl: true`

    * `insecure` (*boolean*) - whether or not self-signed certs should be trusted.

      Example: `insecure: true`

    * `type` (*string* required) - the server type of the MU\*; must match the key of one of the listed server types.

      Example: `type: muck`

**Default**

```yaml
stimmtausch:
    servers:
        spr:
            host: muck.sprmuck.org
            port: 23551
            ssl: true
            insecure: true
            type: muck
        furrymuck:
            host: furrymuck.com
            port: 8899
            ssl: true
            insecure: true
            type: muck
        tapestries:
            host: tapestries.fur.com
            port: 6699
            ssl: true
            insecure: true
            type: muck
        spindizzy:
            host: muck.spindizzy.org
            port: 7073
            ssl: true
            type: muck
```

### Worlds

About
:   `worlds` holds user/character information for a server.

    Expects an object mapping world names to values.

Values
:  
    * `display_name` (*string*) - a free-form display name for the world to be used in the UI above the send buffer.

      Example: `displayname: "FurryMUCK: Foxface"`
    
    * `server` (*string* required) - the name of the server to connect to; must match the key of one of the listed servers.

      Example: `server: furrymuck`

    * `username` (*string* required) - the username to connect as.

      Example: `username: Foxface`

    * `password` (*string* required) - the password to use.

      Example: `password: ILoveSwishyTails`

    * `log` (*boolean*) - whether or not to keep the global logs after disconnecting.

      Example: `log: true`

**Example**

```yaml
stimmtausch:
    worlds:
        fm_foxface:
            display_name: "FurryMUCK: Foxface"
            server: furrymuck
            username: Foxface
            password: ILoveSwishyTails
            log: true
        # More worlds...
```

### Triggers

About
:   `triggers` holds information about automatic behaviors for the client to take.

    Expects a list of triggers. 

Values
:  
    * `name` (*string*) - the name of the trigger.

      Example: `name: "Hilite all my usernames"`

    * `type` (*string* required; one of `hilite`, `gag`, `script`, or `macro`) - what to do when the trigger matches: change the color/attributes of the text, don't show the line at all, run a script (*not implemented*), or run a macro (*not implemented*).

      Example: `type: hilite`

    * `world` (*string* optional; the name (not display name) of a world) - the world to which this trigger should apply. If none is specified, it will apply to every world.

      Example: `world: fm_foxface`

    * `match` (*string* either `match` or `matches` or both are required) - the [regular expression](https://golang.org/pkg/regexp/) to match in the line.

      Example: `match: "[Ff]oxface"`

    * `matches` (*list of strings* either `match` or `matches` or both are required) - a list of matches as above.

      Example: `matches: ["[Ff]oxface", "[Rr]udderbutt"]`

    * `attributes` (*string* required for hilites) - one or more attributes or colors, separated by `+`, which map to [an attribute/color string](https://ansigo.projects.makyo.io).

      Example: `attributes: "bold+bg:grey10+green"`

    * `log_anyway` (*boolean* only used for gags) - when a gag is triggered, it won't be displayed to the screen. It also won't be printed in any open log files, unless this is set to `true`

      Example: `log_anyway: false`

    * `script` (*string* required for scripts) - the path of a script/executable to run (*not implemented*)

    * `macro` (*string* required for macros) - the name of a macro to run (*not implemented*)

Notes
:   As a note, for the time being due to a [bug](https://github.com/makyo/stimmtausch/issues/62), list more specific (e.g: names, single words) hilite triggers after more general triggers that hilite the whole line, or hiliting may stop partway through.

**Example**

```yaml
stimmtausch:
    triggers:
        - name: "Hilite all my usernames"
          type: hilite
          matches: ["[Ff]oxface", "[Rr]udderbutt"]
          attributes: "bold+green"
        - name: "I hate this guy, but he's only on FM..."
          type: gag
          world: furrymuck
          match: "(?i)bad-wolf"
        # More triggers...
```

### Client

*Documentation on this section will be coming soon!*

#### Syslog

`show_syslog`
:   Whether or not to show the syslog in a pane in the UI *(not implemented)*

`log_level`
:   Minimum level of log to output out of TRACE, DEBUG, INFO, WARNING, ERROR, and CRITICAL --- *Default: INFO*

#### Profile

Only one may be set to true at once.

`mem`
:   Whether or not to profile memory --- *Default: false*

`cpu`
:   Whether or not to profile CPU usage --- *Default: false*

#### Logging

`time_string`
:   Date/time format to use in an example string. The numbers matter, because Golang. The date must follow the reference date/time of 3:04:05PM on January 2nd, 2006, Mountain Standard Time (-0700). 1-2 3:4:5 6 7. It's silly, but I don't make [the rules](https://golang.org/pkg/time/#Time.Format). --- *Default: 2006-01-02T150405*

`log_timestamps`
:   Whether or not to include the timestamp in the log files *(not implemented)*

`log_world`
:   Whether or not to keep the log for the connection to the world after disconnecting. --- *Default: true*

#### UI

`scrollback`
:   How many lines received from the connection to keep in memory. --- *Default: 5000*

`history`
:   How many lines sent to the connection to keep in memory. --- *Default: 500*

`unified_history_buffer`
:   Whether or not to keep a separate history buffer for each connection or to have one for all connections --- *(not implemented, effectively true)*

`vim_keybindings`
:   sigh... *(not implemented)*

`indent_first`
:   Number of spaces for indenting the first line of a wrapped line. --- *Default: 0*

`indent_subsequent`
:   Number of spaces for indenting the subsequent lines of a wrapped line --- *Default: 4*

`mouse`
:   Whether or not to capture mouse events (note that, if true, many terminals will not let you click URLs or select text) --- *Default: false*

`colors`
:   [Colors and attributes](https://ansigo.projects.makyo.io/) used in the client UI itself. *Default:*

    ```yaml
    send_title:
      # Focused world
      active: "bold+white"
      # Focused world, scrolled up with new activity
      active_more: "bold+underline+white"
      # Non-focused world
      inactive: "steelblue"
      # Non-focused world with new activity
      inactive_more: "steelblue+underline"
      # Disconnected world (non-focused)
      disconnected: "mediumvioletred"
      # Disconnected world with new activity
      disconnected_more: "mediumvioletred+underline"
      # Disconnected world (non-focused) with unread lines
      disconnected_active: "deeppink3"
    ```
