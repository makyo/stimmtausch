---
layout: post
title: Sending works!
author: Maddy
excerpt: In which Maddy gets really excited about being able to log into a MUCK.
---

[![Screenshot of Stimmtausch connected and talking to a MUCK](/assets/2019-02-01.png)](/assets/2019-02-01.png)

Sending now works! We can talk back and forth with the server.

There's a few immediate problems that jump out to me from today's work. First, visible here, is that there's no wrapping smarts for the text - it just wraps at whatever character. What's also not shown is that if you resize the window, the app resizes fine, but the text does not reflow. That will take a bit of work.

More immediately, though, and hopefully easier to fix is that if you disconnect in the MUCK, when you quit the app, it tries to disconnect again. Since you're not connected, though, it kinda hangs. I know what's going on there, at least!

Here's to more work!
