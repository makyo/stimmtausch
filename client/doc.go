// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

// Package client contains everything required to connect to a server.
//
// Stimmtausch uses a hierarchy of concepts within a client:
//
//     client
//     ├── server type
//     │   ├── server
//     │   │   ├── world
//     │   │   │   └── conection
//     │   │   ├── world
//     │   │   ...
//     │   ├── server
//     │   ...
//     ├── server type
//    ...
//
// A client can contain any number of server types. These are the various kinds
// of games out there, such as MUCKs, MUDs, MOOs, etc. Each server type can, in
// turn, contain any number of servers. Each server can then contain any number
// of worlds, which are the union of a character and a server. Finally, each
// world can contain any number of connections (though, in practice, usually
// only one).
//
// This maps directly to a good chunk of the configuration file, handily.
//
// Responsibility is passed throughout this chain. For instance, you call
// connect on the client, which decides how to do so by looking at the worlds
// and servers it knows about. Ditto closing connections. Similarly, errors
// percolate through the tree.
//
// The client is designed to be agnostic as far as interacting with it goes.
// Stimmtausch has a UI, but that is not required; you can run the client in
// headless mode and it will maintain the connection for you, and you can
// interact with it without using the termbox UI.
//
// This works on the fact that Stimmtausch's client manages the connection and
// all of its interaction through two files: an output file which streams
// data received from the connection, and an input file which is actually a
// FIFO. When you write to the input file, that data gets sent to the server
// while any responses are written to the output file.
//
// This mirrors the IRC client [`ii`](https://tools.suckless.org/ii/), and,
// in turn, the MUCK client [`mm`](https://github.com/onlyhavecans/mm). The
// latter was a heavy inspiration for Stimmtausch.
//
// The upside to this is that the combination of everything being an instance
// of io.Writer - the output file, the FIFO, the history buffers, the termbox
// views, everything - is that the whole thing is just a well-organized set of
// fmt.Fprint*() statements once it's all set up.
package client
