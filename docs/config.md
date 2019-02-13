---
layout: default
title: Configuration
---

## Configuration

Stimmtausch is designed to be configurable from the ground up, as much as possible. It relies on configuration files in one of three formats: YAML, TOML, or JSON. These files are read in order from:

* `/etc/stimmtausch/st.yaml`
* `/etc/stimmtausch/conf.d/*`
* `$HOME/.config/stimmtausch/*.st.*`
* `$HOME/.config/stimmtausch/*/*.st.*`
* `$HOME/.strc`

Any values showing up in later files override values in earlier files. See [snuffler](https://github.com/makyo/snuffler) for more on that.

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

**Example**

```yaml
stimmtausch:
    triggers:
        - name: "Hilite all my usernames"
          type: hilite
          matches: ["[Ff]oxface", "[Rr]udderbutt"]
          attributes: "bold+green"
        - name: "I hate this guy..."
          type: gag
          match: "(?i)bad-wolf"
        # More triggers...
```

### Client

*Documentation on this section will be coming soon!*

#### Syslog

#### Logging

#### UI
