---
layout: post
title: Stimmtausch 0.0.1
author: Maddy
excerpt: In which the very first release of Stimmtausch is announced.
---

[![Yay!](/assets/2019-02-13.png)](/assets/2019-02-13.png)

Yay! I did the thing! I finished all of the tasks on the [0.0.1 milestone](https://github.com/makyo/stimmtausch/milestone/1), beat my head against `debuild`, and managed to pull together the first release of Stimmtausch! You can now install it by doing the following:

```bash
sudo add-apt-repository -u ppa:makyo/st
sudo apt install stimmtausch
```

Or, go [here](https://github.com/makyo/stimmtausch/releases/tag/0.0.1) and follow the instructions for downloading a .deb and installing that.

This will get you up and running with the client, but you'll have to edit your `~/.strc` with some of [configuration information](/config). A good start for a `~/.strc` would look something like this:

```yaml
stimmtausch:
    worlds:
        fm:
            display_name: "FurryMUCK: Foxface"
            server: furrymuck
            username: Foxface
            password: Swishytail1
        taps:
            display_name: "Tapestries: Foxface"
            server: tapestries
            username: Foxface
            password: Swishytail1
    triggers:
        - type: hilite
          match: "[Ff]oxface"
          attributes: "bold+green"
```
