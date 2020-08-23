---
layout: post
title: Character counts ahoy!
author: Maddy
excerpt: In which Maddy helps our buffer-busting friends.
---

<script id="asciicast-etrki2gwV9hTvm8HnT3UOz9vv" src="https://asciinema.org/a/etrki2gwV9hTvm8HnT3UOz9vv.js?cols=112" async></script>

Lookie there! Now you can see how much you've typed and whether or not you're going over your buffer! Not only that, but you can now enter more than one line at a time!

There are a few additional fixes to go along with this, such as some hardened work around stopping errors and an improved error display mechanism (anything from `WARNING` and above gets shown to the user in a modal). There's still a bit of work to go before this all gets promoted to Beta, but you can preview these changes in the `edge` channel on the snap with:

    snap install stimmtausch --channel=edge

Or, if you already have it installed:

    snap refresh stimmtausch --channel=edge
