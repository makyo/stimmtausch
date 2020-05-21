---
layout: post
author: Maddy
title: Stimmtausch 0.0.3 and stimmtausch.vim
excerpt: In which Maddy wises from her gwave to provide s'more Stimmtauschness
---

[![Yay!](/assets/2020-05-08.png)](/assets/2020-05-08.png)

Well well well, if it isn't *Maddy*. Back at it again with the Stimmtauschery!

I know it's been a while. Things have been really dumb! 2020 has been the longest fucking decade! Anyway, Since last I wrote, I've switched jobs, lost my job, published four books, and have two more on the way, so I'm just all over the place!

Unfortunately, I also kind of lost the will to work on software for a bit, and Stimmtausch kind of languished. Just so much burnout going around that it was hard to get excited about software for a bit there. I still used the client daily, and I still got super frustrated at all the bugs I was lacking. Eventually, though, [@cyveris](https://github.com/cyveris) got me up and moving again (thanks!) and I got back to work.

So! Two big updates! The first, of course, is that the [milestone for 0.0.3](https://github.com/makyo/stimmtausch/projects/4) has been completed. There is one big feature and several bugfixes:

* **Major change:** the command is now `stimmtausch`, not `st`. This is due to a packaging clash with [simple terminal](https://st.suckless.org/). If you are not using that st, you can create an alias for yourself in your shell rc file, or via `update-alternatives --install`.
* Log output from current world to a file with `/log filename`. You can turn it off with `/log --off`, list current open logs with `/log --list`, and get help with `/log --help`, which uses the new-ish modal system.
* Home and end keys are now mapped.
* Bug fix: hilites that match more than once in a line no longer clash.
* Bug fix: screen is now cleared on quit, where views used to linger.
* Bug fix: don't barf if ~/.local/log/stimmtausch doesn't exist on startup.
* Bug fix: don't barf when switching away from a disconnected world and then trying to switch back.
* Bug fix: on that note, make sure that, when an inactive world disconnects, the title bar state changes.

@cyveris will be helping me out a lot more in terms of project organization, so hopefully there will be more work coming up, too.

Now, a bit of bad news on this release: something about the combination of launchpad and my Debian packaging setup stopped playing well together, so there is no release in the PPA at this time. ~~If you want to install from a .deb, you will have to do so by [downloading it](https://github.com/makyo/stimmtausch/releases/tag/0.0.3) and installing via `dpkg -i`.~~ I gave up on this! Loser!

Instead, I'm moving to snaps! The upside to this is that development can move a lot faster because releases don't hurt quite so much. It will soon be available on the snap store. It requires one manual override in order to use `~/.config/stimmtausch` for your configuration files, but they're working on making that automatic. See the latest post for more information.

Anyway!

There's another new addition to Stimmtausch, which is [stimmtausch.vim](https://github.com/makyo/stimmtausch.vim). This lets you use vim as the input field for Stimmtausch if you're running it in headless mode. You can find out much more on how to use it on the project page.

Once there's an update on either the snap or deb front, I'll make another post!
