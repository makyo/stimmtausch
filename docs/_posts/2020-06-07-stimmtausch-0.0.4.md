---
layout: post
title: Stimmtausch 0.0.4
excerpt: In which Maddy gets writers' block and works on software instead.
---

Woo wow hi yay!

Was running into some writers block (though patrons will get some poems in a bit!), so I've been doing some work on Stimmtausch. Have accomplished the following:

- Removed the reliance on `/etc` for global log files, meaning that a binary of the app can just be downloaded and run lickity split.
- Triggers such as hilites and gags can now be world-specific
- Help is now shown in a modal overlay, which is nice.
- Keybindings for worldswitching: `Esc+<arrows>` rotates (equiv to `/>` and `/<`), and `Ctrl+<square brackets>` rotates to the next world with new activity (equiv to the new `/]` and `/[`)
- A few bugfixes: (`/connect` without a world no longer crashes, `/help` while not connected no longer crashes, connecting to a disconnected world while not active does not cause world switching problems)

These are all released on the snap store in the `edge` and `beta` channels as 0.0.4rev1. Pushing is way easier now that there's an automatic override in place. There was a bit of confusion as the wrong branch was being built in the snap, but hey, that's fixed!

[![help!](/assets/2020-06-07.png)](/assets/2020-06-07.png)
