// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package help_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/makyo/stimmtausch/help"
)

func TestText(t *testing.T) {
	Convey("One can render help as ANSIfied text", t, func() {

		h := help.Help{
			Name:      "rose",
			ShortDesc: "the Doctor's companion",
			Synopsis: map[string]string{
				"Portrayed by": "Billie Piper",
				"Affiliation":  "9th/10th Doctor",
			},
			Overview:    "Rose Tyler was a companion of The Doctor.",
			Description: "Rose Tyler is a fictional character in the British science fiction television series Doctor Who. She was created by series producer Russell T Davies and portrayed by Billie Piper.",
			SeeAlso:     "Mickey Smith",
		}
		expected := "\x1b[1mOVERVIEW\x1b[0m\n\n    Rose Tyler was a companion of The Doctor.\n\n\x1b[1mSYNOPSIS\x1b[0m\n\n    \t\x1b[1mrose\x1b[0m \x1b[4mPortrayed by\x1b[0m \tBillie Piper\n    \t\x1b[1mrose\x1b[0m \x1b[4mAffiliation\x1b[0m \t9th/10th Doctor\n\n\x1b[1mDESCRIPTION\x1b[0m\n\n    Rose Tyler is a fictional character in the British science fiction television series Doctor Who. She was created by series producer Russell T Davies and portrayed by Billie Piper.\n\n\x1b[1mSEE ALSO\x1b[0m\n\n    Mickey Smith"
		So(help.RenderText(h), ShouldEqual, expected)
	})
}
