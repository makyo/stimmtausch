package macro_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/makyo/stimmtausch/macro"
)

func TestBuiltins(t *testing.T) {

	Convey("Environment builtins", t, func() {
		e := macro.New()
		listener := make(chan macro.MacroResult)
		e.AddListener(listener)

		Convey("There are passthrough builtins", func() {
			passthroughTests := []string{"_", "connect", "disconnect"}
			for _, m := range passthroughTests {
				go e.Dispatch(m, "rose tyler")
				result := <-listener
				So(result.Name, ShouldEqual, m)
				So(result.Results, ShouldResemble, []string{"rose", "tyler"})
				So(result.Err, ShouldBeNil)
			}
		})

		Convey("fg builtin", func() {
			go e.Dispatch("fg", "<")
			result := <-listener
			So(result.Name, ShouldEqual, "fg")
			So(result.Results, ShouldResemble, []string{"rotate", "-1"})
			So(result.Err, ShouldBeNil)

			go e.Dispatch("fg", ">")
			result = <-listener
			So(result.Name, ShouldEqual, "fg")
			So(result.Results, ShouldResemble, []string{"rotate", "1"})
			So(result.Err, ShouldBeNil)

			go e.Dispatch("fg", "rose_tyler")
			result = <-listener
			So(result.Name, ShouldEqual, "fg")
			So(result.Results, ShouldResemble, []string{"switch", "rose_tyler"})
			So(result.Err, ShouldBeNil)

			go e.Dispatch("fg", "bad wolf")
			result = <-listener
			So(result.Name, ShouldEqual, "fg")
			So(result.Results, ShouldResemble, []string{})
			So(result.Err.Error(), ShouldEqual, `received wrong number of args, wanted 1, got 2 (["bad" "wolf"])`)
		})
	})
}
