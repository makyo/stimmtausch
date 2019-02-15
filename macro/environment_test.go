package macro_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/makyo/stimmtausch/macro"
)

func TestEnvironment(t *testing.T) {

	Convey("When working with an macro", t, func() {

		Convey("One can create a new one", func() {
			e := macro.New()
			So(e, ShouldNotBeNil)
		})

		Convey("One can dispatch and listen for results", func() {
			e := macro.New()
			l1 := make(chan macro.MacroResult)
			l2 := make(chan macro.MacroResult)
			e.AddListener(l1)
			e.AddListener(l2)
			go e.Dispatch("_", "rose tyler")
			m1 := <-l1
			m2 := <-l2
			So(m1, ShouldResemble, m2)
			So(m1.Name, ShouldEqual, "_")
			So(m1.Results, ShouldResemble, []string{"rose", "tyler"})
			So(m1.Err, ShouldBeNil)

			Convey("With an error if the macro doesn't exist", func() {
				go e.Dispatch("bad-wolf", "nonesuch")
				m1 := <-l1
				m2 := <-l2
				So(m1, ShouldResemble, m2)
				So(m1.Err.Error(), ShouldEqual, "unknown macro bad-wolf")
			})
		})
	})
}
