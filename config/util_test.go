package config_test

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	homedir "github.com/mitchellh/go-homedir"

	"github.com/makyo/stimmtausch/config"
)

func TestUtils(t *testing.T) {
	Convey("When setting directories", t, func() {
		homeDir, err := homedir.Dir()
		So(err, ShouldBeNil)

		Convey("It respects snap settings", func() {
			config.InitDirs()
			So(config.ConfigDir, ShouldEqual, filepath.Join(homeDir, ".config", "stimmtausch"))
			So(config.WorkingDir, ShouldEqual, filepath.Join(homeDir, ".local", "share", "stimmtausch"))
			So(config.LogDir, ShouldEqual, filepath.Join(homeDir, ".local", "log", "stimmtausch"))

			snapCommon := "/home/rose/snap/stimmtausch/common"
			os.Setenv("SNAP_USER_COMMON", snapCommon)
			config.InitDirs()
			So(config.ConfigDir, ShouldEqual, filepath.Join(homeDir, ".config", "stimmtausch"))
			So(config.WorkingDir, ShouldEqual, filepath.Join(snapCommon, "share", "stimmtausch"))
			So(config.LogDir, ShouldEqual, filepath.Join(snapCommon, "log", "stimmtausch"))
		})
	})
}
