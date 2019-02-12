package config_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/makyo/stimmtausch/config"
)

func TestTriggers(t *testing.T) {
	Convey("When creating triggers", t, func() {

		Convey("They get a default name", func() {
			c := stubConfig()
			So(len(c.Triggers[0].Name), ShouldEqual, 0)
			errs := c.FinalizeAndValidate()
			So(len(errs), ShouldEqual, 0)
			So(len(c.CompiledTriggers[0].Name), ShouldNotEqual, 0)
		})

		Convey("A trigger may only have certain types", func() {
			c := stubConfig()
			c.Triggers = append(c.Triggers, config.Trigger{
				Type:  "bad-wolf",
				Match: "bad-wolf",
			})
			errs := c.FinalizeAndValidate()
			So(len(errs), ShouldEqual, 1)
			So(errs[0].Error(), ShouldEqual, "unknown trigger type bad-wolf")
		})

		Convey("A trigger with no matches is an error", func() {
			c := stubConfig()
			c.Triggers = append(c.Triggers, config.Trigger{
				Type:  "gag",
				Match: "",
			})
			errs := c.FinalizeAndValidate()
			So(len(errs), ShouldEqual, 1)
			So(errs[0].Error(), ShouldStartWith, "no matches for trigger")
		})

		Convey("A trigger with an invalid regexp is an error", func() {
			c := stubConfig()
			c.Triggers = append(c.Triggers, config.Trigger{
				Type:  "gag",
				Match: "*asdf(",
			})
			errs := c.FinalizeAndValidate()
			So(len(errs), ShouldEqual, 1)
			So(errs[0].Error(), ShouldStartWith, "error parsing regexp")
		})
	})

	Convey("When running triggers", t, func() {
		c := stubConfig()
		errs := c.FinalizeAndValidate()
		So(len(errs), ShouldEqual, 0)
		var (
			hilite = c.CompiledTriggers[0]
			gag    = c.CompiledTriggers[1]
			macro  = c.CompiledTriggers[2]
			script = c.CompiledTriggers[3]
		)

		Convey("One can tell whether or not they apply", func() {
			applies, _, _ := hilite.Run("Hello, Rose", c)
			So(applies, ShouldBeTrue)
			applies, _, _ = hilite.Run("bad-wolf", c)
			So(applies, ShouldBeFalse)
		})

		Convey("They can hilite", func() {

			Convey("Once", func() {
				_, line, errs := hilite.Run("Hello, Rose", c)
				So(len(errs), ShouldEqual, 0)
				So(line, ShouldEqual, "Hello, \x1b[1mRose\x1b[22m")
			})

			Convey("More than once", func() {
				_, line, errs := hilite.Run("-I'm the Doctor -Doctor who?", c)
				So(len(errs), ShouldEqual, 0)
				So(line, ShouldEqual, "-I'm \x1b[1mthe Doctor\x1b[22m -\x1b[1mDoctor\x1b[22m who?")
			})
		})

		Convey("They can gag (but what to do about that is on the client)", func() {
			applies, line, errs := gag.Run("bad-wolf", c)
			So(len(errs), ShouldEqual, 0)
			So(line, ShouldEqual, "bad-wolf")
			So(applies, ShouldBeTrue)
		})

		Convey("They can call a macro", func() {
			_, line, errs := macro.Run("Mickey Smith", c)
			So(line, ShouldEqual, "Mickey Smith")

			Convey("NOT IMPLEMENTED YET", func() {
				So(len(errs), ShouldEqual, 1)
				So(errs[0].Error(), ShouldEqual, "not implemented")
			})
		})

		Convey("They can call a script", func() {
			_, line, errs := script.Run("Donna Noble", c)
			So(line, ShouldEqual, "Donna Noble")

			Convey("NOT IMPLEMENTED YET", func() {
				So(len(errs), ShouldEqual, 1)
				So(errs[0].Error(), ShouldEqual, "not implemented")
			})
		})
	})
}
