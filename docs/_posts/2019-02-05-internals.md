---
layout: post
title: Some internals
author: Maddy
---

[![Stimmtausch, system logs, and htop](/assets/2019-02-05.png)](/assets/2019-02-05.png)

Lots going on, but not all of it visible! I'll boil it down to a few points, though!

The [terminal UI library](https://github.com/jroimartin/gocui) I was using seems to have stalled and I ran into a few bugs with it. Unwilling to be enough of a pest to force the author to come back just to land a PR, I did the FOSS thing and forked it into [GoTUI](https://github.com/makyo/gotui), then immediately made a whole bunch of other changes I was worried I was going to have to do in my code, anyway, so hey, that's a win! One of those changes is wrapping on word (rather than on character) in a view, which you can kiiinda see in that little screenshot. Also: nice memory and CPU usage :D

The other big thing is that the [config management library](https://github.com/spf13/viper) I was using wasn't cutting it. It only supported one config file, and I wanted many. After all, maybe multiple users of a computer want to run Stimmtausch, but don't want to share a config (since that'd have their passwords and all)? And certainly they don't want to have to keep every configurable option, server, world, trigger, and so on in one enormous config file! Many things will append/include config files - both nginx and Apache do this - and some things will clobber all changes - see: shells - but I wanted something in between: I wanted global config that could be overridden by local config, n layers deep. I'm positive there's some libraries out there for doing this, but it was a fun problem, so I figured I'd tackle it myself.

Thus [Snuffler](https://github.com/makyo/snuffler), which snuffles about in the locations you give it for config files, then loads them all in order from global to local. (You can also "Snorfle" into a separate object. And all the testing uses "snoot" as a variable. I shouldn't be allowed near a computer.)

Anyway, now that both of these are fully plugged in, I'll get to work on implementing triggers, highlights, and gags!
