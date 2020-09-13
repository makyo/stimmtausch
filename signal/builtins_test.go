package signal_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/makyo/stimmtausch/signal"
)

func TestBuiltins(t *testing.T) {

	Convey("Environment builtins", t, func() {
		e := signal.NewDispatcher()
		listener := make(chan signal.Signal)
		e.AddListener("l", listener)

		Convey("There are passthrough builtins", func() {
			passthroughTests := []string{"_", "connect"}
			for _, m := range passthroughTests {
				go e.Dispatch(m, "rose tyler")
				result := <-listener
				So(result.Name, ShouldEqual, m)
				So(result.Payload, ShouldResemble, []string{"rose tyler"})
				So(result.Err, ShouldBeNil)
			}
		})

		Convey("fg builtin", func() {
			go e.Dispatch("fg", "<")
			result := <-listener
			So(result.Name, ShouldEqual, "fg")
			So(result.Payload, ShouldResemble, []string{"rotate", "-1"})
			So(result.Err, ShouldBeNil)

			go e.Dispatch("fg", ">")
			result = <-listener
			So(result.Name, ShouldEqual, "fg")
			So(result.Payload, ShouldResemble, []string{"rotate", "1"})
			So(result.Err, ShouldBeNil)

			go e.Dispatch("fg", "rose_tyler")
			result = <-listener
			So(result.Name, ShouldEqual, "fg")
			So(result.Payload, ShouldResemble, []string{"switch", "rose_tyler"})
			So(result.Err, ShouldBeNil)
		})
	})
}
