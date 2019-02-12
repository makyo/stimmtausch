package ui_test

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/makyo/stimmtausch/ui"
)

func TestHistory(t *testing.T) {
	Convey("One can create a new empty buffer", t, func() {
		h := ui.NewHistory(100)
		So(h, ShouldNotBeNil)
		So(h.Size(), ShouldEqual, 0)
	})

	Convey("When using a history buffer", t, func() {

		Convey("One can write to it", func() {
			h := ui.NewHistory(100)
			l, err := fmt.Fprint(h, "rose")
			So(l, ShouldEqual, 4)
			So(err, ShouldBeNil)
		})

		Convey("One can get its size", func() {
			h := ui.NewHistory(100)
			So(h.Size(), ShouldEqual, 0)

			Convey("Which increases as you write to it", func() {
				fmt.Fprint(h, "tyler")
				So(h.Size(), ShouldEqual, 1)
			})
		})

		Convey("One can read from it", func() {
			h := ui.NewHistory(100)
			fmt.Fprint(h, "Rose Tyler")
			fmt.Fprint(h, "Mickey Smith")
			fmt.Fprint(h, "Donna Noble")

			Convey("One can get the current line", func() {
				So(h.Current().Text, ShouldEqual, "Donna Noble")
			})

			Convey("One can move back in time", func() {
				So(h.Current().Text, ShouldEqual, "Donna Noble")

				Convey("Moving back returns the current line to make working with input buffers make more sense", func() {
					So(h.Back().Text, ShouldEqual, "Donna Noble")

					Convey("And that moves the current line", func() {
						So(h.Current().Text, ShouldEqual, "Mickey Smith")
					})
				})

				So(h.Last().Text, ShouldEqual, "Donna Noble")

				Convey("One can move forward in time", func() {
					So(h.Back().Text, ShouldEqual, "Donna Noble")
					So(h.Current().Text, ShouldEqual, "Mickey Smith")
					So(h.Forward().Text, ShouldEqual, "Donna Noble")

					Convey("And that moves the current line", func() {
						So(h.Current().Text, ShouldEqual, "Donna Noble")
					})
				})
			})
		})

		Convey("It collects date stamps", func() {
			h := ui.NewHistory(100)
			fmt.Fprint(h, "Rose Tyler")
			rt := h.Current()
			fmt.Fprint(h, "Mickey Smith")
			ms := h.Current()
			fmt.Fprint(h, "Donna Noble")
			dn := h.Current()

			So(rt.Timestamp.Before(ms.Timestamp), ShouldBeTrue)
			So(rt.Timestamp.Before(dn.Timestamp), ShouldBeTrue)
			So(ms.Timestamp.Before(dn.Timestamp), ShouldBeTrue)
		})
	})
}
