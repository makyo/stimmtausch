# Stimmtausch

Stimmtausch (`st`) is a curses-based MUCK client written in Go.

## Reasoning

TinyFugue (`tf`) has been the gold-standard terminal MU\* client for decades. However, it has largely been abandoned (thankfully in a very stable state!). Rather than attempt to take over the project and restart it, it was decided that a new client would be the best path forward: a modern client using modern tooling and design.

## The name

*Stimmtausch* is a compositional technique that shows up in fugue. Often called "voice exchange", it refers to the changing of motifs between the voices in a fugue which often occurs before the recapitulation. The goal was to have a name that expressed something detailed and well-thought-out. "SmallSonata" was also considered, but discarded as too twee :)

## Inspiration

TinyFugue, of course, was a big inspiration for this project, but [`mm`](https://github.com/onlyhavecans/mm) played a role as well, as a MUCK client written in Go using some interesting ideas. Many portions of `mm` are found in Stimmtausch, though, rather than use it as a library, it was modified to fit the design goals of the project.

## Libraries

* gocui
* cobra
* viper
* goconvey

## Contributors

* Madison Scott-Clary (@makyo; Maddy on FurryMUCK and Tapestries).
