// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package help

var HelpMessages = map[string]Help{
	"help": Help{
		Name:      "/help",
		ShortDesc: "command help",
		Synopsis: map[string]string{
			"":          "show this help",
			"<command>": "show help for <command> (without the `/`)",
		},
		Overview:    "Give help on an internal command",
		Description: "Find out more information on topics in Stimmtausch by using the /help command. For instance, to learn about the `/log` command, type `/help log`.\n\n`Available topics:`\n\n{HELPTOPICS}",
	},

	"log": Help{
		Name:      "/log",
		ShortDesc: "connection logging",
		Synopsis: map[string]string{
			"":             "show this help",
			"--help":       "show this help",
			"--list":       "list open log files",
			"<file>":       "start logging the current world to the specified file",
			"--off <file>": "stop logging to the specified file",
		},
		Overview:    "Command to control logging output from worlds.",
		Description: "Logging in Stimmtausch is controlled through the /log command. Invoked with a file name, it starts logging the current world's output to the specified file (absolute, or relative to the directory in which Stimmtausch was started). You can turn logging off at any time by calling `/log --off <file>`. To list what logs are open, you can call `/log --list`",
	},

	"fg": Help{
		Name:      "/fg",
		ShortDesc: "bring world to the foreground",
		Synopsis: map[string]string{
			"":        "rotate to the next active world to the right (same as `/]`)",
			">":       "rotate to the next world to the right",
			"<":       "rotate to the next world to the left",
			"]":       "rotate to the next active world to the right",
			"[":       "rotate to the next active world to the left",
			"<world>": "switch to the named world",
		},
		Overview:    "Command to control moving between worlds.",
		Description: "Moving between worlds in Stimmtausch is accomplished with the /fg command. You can rotate between worlds by using the special world names > and <, otherwise you can specifi which world you would like to bring to the foreground.",
		SeeAlso:     "`/>` (same as `/fg >`), `/<` (same as `/fg <`), `/]` (same as `/fg ]`), `/[` (same as `/fg [`)",
	},

	"connect": Help{
		Name:      "/connect",
		ShortDesc: "connect to worlds",
		Synopsis: map[string]string{
			"<named world>": "connect to the named world",
			//"<named server>": "connect to the named server without a username",
			//"<address>:<port>": "connect to the server address and port specified",
		},
		Overview:    "Command to connect to worlds.",
		Description: "Connecting to worlds in Stimmtausch is accomplished with the /connect command. You can connect to worlds named in your configuration files.", // Additionally, you can connect to servers named in your configuration without user information, or to a specified address and port. In each of the latter two cases, you will be given a temporary world name which you can use with other commands.
		SeeAlso:     "`/disconnect`, `/fg`",
	},

	"disconnect": Help{
		Name:      "/disconnect",
		ShortDesc: "disconnect from worlds",
		Synopsis: map[string]string{
			"<world>": "disconnect from the world specified",
		},
		Overview:    "Command to disconnect from worlds.",
		Description: "Disconnecting from worlds in Stimmtausch is accomplished with the /disconnect command. It accepts a world name.",
		SeeAlso:     "`/connect`",
	},

	"quit": Help{
		Name:      "/quit",
		ShortDesc: "quit Stimmtausch",
		Synopsis: map[string]string{
			"": "disconnect from all worlds and quit Stimmtausch",
		},
		Overview:    "Command to quit Stimmtausch.",
		Description: "Quitting Stimmtausch is accomplished to the /quit command.", // Note that if you send `/quit` from _any_ client attached to Stimmtausch (e.g: if you're using Stimmtausch in headless mode or as a server), it will quit, detaching every connected client.",
	},

	"syslog": Help{
		Name:      "/syslog",
		ShortDesc: "log to the system log",
		Synopsis: map[string]string{
			"<level> <message>": "You may log arbitrary information to the system log via the /syslog command. Why? We're sure you have your reasons! The log level is the first argument, and may be one of 'TRACE', 'DEBUG', 'INFO', 'WARNING', 'ERROR', or 'CRITICAL'.",
		},
	},
}
