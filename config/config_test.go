package config_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/makyo/stimmtausch/config"
)

func stubConfig() *config.Config {
	return &config.Config{
		Version: 1,
		ServerTypes: map[string]config.ServerType{
			"stubtype": config.ServerType{
				ConnectString:    "connect $username $password",
				DisconnectString: "QUIT",
			},
		},
		Servers: map[string]config.Server{
			"stubserver": config.Server{
				Host:       "example.com",
				Port:       12345,
				SSL:        false,
				Insecure:   false,
				ServerType: "stubtype",
			},
		},
		Worlds: map[string]config.World{
			"stubworld": config.World{
				DisplayName: "TARDIS",
				Server:      "stubserver",
				Username:    "user",
				Password:    "pass",
				Log:         false,
			},
		},
		Triggers: []config.Trigger{
			config.Trigger{
				Type:       "hilite",
				Matches:    []string{"(?i)(the )?doctor", "(?i)rose( tyler)?"},
				Attributes: "bold",
			},
			config.Trigger{
				Type:  "gag",
				Match: "bad-wolf",
			},
			config.Trigger{
				Type:  "macro",
				Match: "Mickey Smith",
			},
			config.Trigger{
				Type:  "script",
				Match: "Donna Noble",
			},
		},
		Client: config.Client{
			Syslog: config.Syslog{
				ShowSyslog: false,
				LogLevel:   "INFO",
			},
			Logging: config.Logging{
				TimeString:    "",
				LogTimestamps: false,
				LogWorld:      false,
			},
			UI: config.UI{
				Scrollback:           100,
				History:              100,
				UnifiedHistoryBuffer: false,
				VimKeybindings:       false,
				IndentFirst:          0,
				IndentSubsequent:     4,
				Mouse:                false,
			},
		},
		HomeDir:    "/home/rose",
		ConfigDir:  "/home/rose/.config/stimmtausch",
		WorkingDir: "/home/rose/.local/share/stimmtausch",
		LogDir:     "/home/rose/.local/log/stimtausch",
	}
}

func TestConfig(t *testing.T) {
	Convey("When creating config", t, func() {

		Convey("It can be validated and finalized", func() {

			Convey("If valid, it sets names on worlds and servers and compiles triggers", func() {
				c := stubConfig()
				errs := c.FinalizeAndValidate()
				So(len(errs), ShouldEqual, 0)
				So(c.Worlds["stubworld"].Name, ShouldEqual, "stubworld")
				So(c.Servers["stubserver"].Name, ShouldEqual, "stubserver")
				So(len(c.CompiledTriggers), ShouldEqual, 4)
			})

			Convey("It requires a version greater than 0", func() {
				c := stubConfig()
				c.Version = 0
				errs := c.FinalizeAndValidate()
				So(len(errs), ShouldEqual, 1)
				So(errs[0].Error(), ShouldStartWith, "version key wasn't set")
			})

			Convey("Worlds must refer to existing servers", func() {
				c := stubConfig()
				w := c.Worlds["stubworld"]
				w.Server = "bad-wolf"
				c.Worlds["stubworld"] = w
				errs := c.FinalizeAndValidate()
				So(len(errs), ShouldEqual, 1)
				So(errs[0].Error(), ShouldEqual, "world stubworld refers to unknown server bad-wolf")
			})

			Convey("Servers must refer to existing server types", func() {
				c := stubConfig()
				s := c.Servers["stubserver"]
				s.ServerType = "bad-wolf"
				c.Servers["stubserver"] = s
				errs := c.FinalizeAndValidate()
				So(len(errs), ShouldEqual, 1)
				So(errs[0].Error(), ShouldEqual, "server stubserver refers to unknown server type bad-wolf")
			})
		})
	})
}
