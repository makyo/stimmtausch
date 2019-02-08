// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package main

import (
	"github.com/makyo/st/cmd"
)

func main() {
	cmd.GenMarkdownDocs()
	cmd.GenManPages()
}
