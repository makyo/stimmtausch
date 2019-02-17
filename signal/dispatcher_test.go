package signal_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/makyo/stimmtausch/signal"
)

func TestEnvironment(t *testing.T) {

	Convey("When working with an signal", t, func() {

		Convey("One can create a new one", func() {
			e := signal.NewDispatcher()
			So(e, ShouldNotBeNil)
		})

		Convey("One can dispatch and listen for results", func() {
			e := signal.NewDispatcher()
			l1 := make(chan signal.Signal)
			l2 := make(chan signal.Signal)
			e.AddListener("l1", l1)
			e.AddListener("12", l2)
			go e.Dispatch("_", "rose tyler")
			m1 := <-l1
			m2 := <-l2
			So(m1, ShouldResemble, m2)
			So(m1.Name, ShouldEqual, "_")
			So(m1.Payload, ShouldResemble, []string{"rose tyler"})
			So(m1.Err, ShouldBeNil)

			Convey("With an error if the signal doesn't exist", func() {
				go e.Dispatch("bad-wolf", "nonesuch")
				m1 := <-l1
				m2 := <-l2
				So(m1, ShouldResemble, m2)
				So(m1.Err.Error(), ShouldEqual, "unknown macro bad-wolf")
			})
		})
	})
}
