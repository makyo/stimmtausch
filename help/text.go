// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package help

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"

	ansi "github.com/makyo/ansigo"
)

var referRe = regexp.MustCompile("`([^`]+)`")
var longestRe = regexp.MustCompile("\t([^\t]+)\t")
var ansiRe = regexp.MustCompile("\x1b[^m]+m")

var textFuncs = template.FuncMap{
	"synopsis":    textSynopsis,
	"seeAlso":     textSeeAlso,
	"description": textDescription,
}
var textTemplate = template.Must(template.New("text").Funcs(textFuncs).Parse(`{{"\x1b[1m"}}OVERVIEW{{"\x1b[0m"}}

    {{.Overview }}

{{"\x1b[1m"}}SYNOPSIS{{"\x1b[0m"}}

    {{synopsis .Name .Synopsis}}

{{"\x1b[1m"}}DESCRIPTION{{"\x1b[0m"}}

    {{description .Description}}{{seeAlso .SeeAlso}}`))

func textSynopsis(name string, synMap map[string]string) string {
	lines := []string{}
	for cmd, desc := range synMap {
		if len(cmd) > 0 {
			lines = append(lines, fmt.Sprintf("    \t%s %s\t%s", ansi.MaybeApplyOneWithReset("bold", name), ansi.MaybeApplyOneWithReset("underline", cmd), refer(desc)))
		} else {
			lines = append(lines, fmt.Sprintf("    \t%s\t%s", ansi.MaybeApplyOneWithReset("bold", name), refer(desc)))
		}
	}
	return strings.Join(tabSubst(lines), "\n    ")
}

func textSeeAlso(seeAlso string) string {
	if len(seeAlso) != 0 {
		return fmt.Sprintf("\n\n%s\n\n    %s", ansi.MaybeApplyOneWithReset("bold", "SEE ALSO"), refer(seeAlso))
	}
	return ""
}

func textDescription(description string) string {
	topics := []string{}
	for cmd, h := range HelpMessages {
		topics = append(topics, fmt.Sprintf("    \t%s\t%s", ansi.MaybeApplyWithReset("cyan", cmd), h.ShortDesc))
	}
	return strings.Join(tabSubst(topics), "\n    ")
}

func refer(s string) string {
	return string(referRe.ReplaceAll([]byte(s), []byte(ansi.MaybeApplyWithReset("cyan", "$1"))))
}

func tabSubst(lines []string) []string {
	longest := 0
	for _, s := range lines {
		cmd := ansiRe.ReplaceAll(longestRe.FindAllSubmatch([]byte(s), -1)[0][1], []byte(""))
		if len(cmd) > longest {
			longest = len(cmd)
		}
	}
	for i, s := range lines {
		cmd := longestRe.FindAllSubmatch([]byte(s), -1)[0][1]
		lines[i] = string(longestRe.ReplaceAll([]byte(s), []byte(fmt.Sprintf("$1 %s", strings.Repeat(" ", longest-len(ansiRe.ReplaceAll([]byte(cmd), []byte("")))+1)))))
	}
	return lines
}

func RenderText(h Help) string {
	var b strings.Builder
	if err := textTemplate.Execute(&b, h); err != nil {
		log.Errorf("unable to render help: %v", err)
		return ""
	}
	return b.String()
}
