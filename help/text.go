// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package help

import (
	"fmt"
	"strings"
	"text/template"

	ansi "github.com/makyo/ansigo"
)

var textFuncs = template.FuncMap{
	"synopsis": textSynopsis,
	"seeAlso":  textSeeAlso,
}
var textTemplate = template.Must(template.New("text").Funcs(textFuncs).Parse(`{{"\x1b[1m"}}OVERVIEW{{"\x1b[0m"}}

    {{.Overview }}

{{"\x1b[1m"}}SYNOPSIS{{"\x1b[0m"}}

    {{synopsis .Name .Synopsis}}
	
{{"\x1b[1m"}}DESCRIPTION{{"\x1b[0m"}}

    {{.Description}}
{{seeAlso .SeeAlso}}`))

func textSynopsis(name string, synMap map[string]string) string {
	lines := []string{}
	for cmd, desc := range synMap {
		if len(cmd) > 0 {
			lines = append(lines, fmt.Sprintf("\t%s %s \t%s", ansi.MaybeApplyOneWithReset("bold", name), ansi.MaybeApplyOneWithReset("underline", cmd), desc))
		} else {
			lines = append(lines, fmt.Sprintf("\t%s \t%s", ansi.MaybeApplyOneWithReset("bold", name), desc))
		}
	}
	return strings.Join(lines, "\n    ")
}

func textSeeAlso(seeAlso string) string {
	if len(seeAlso) != 0 {
		return fmt.Sprintf("%s\n    %s", ansi.MaybeApplyOneWithReset("bold", "SEE ALSO"), seeAlso)
	}
	return ""
}

func RenderText(h Help) string {
	var b strings.Builder
	if err := textTemplate.Execute(&b, h); err != nil {
		log.Errorf("unable to render help: %v", err)
		return ""
	}
	return b.String()
}
