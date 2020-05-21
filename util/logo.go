package util

import (
	ansi "github.com/makyo/ansigo"
)

func Logo() string {
	return ansi.MaybeApply("111", "  _____ _   _                     _                       _     \n") +
		ansi.MaybeApply("111", " /  ___| | (_)                   | |                     | |    \n") +
		ansi.MaybeApply("212", " \\ `--.| |_ _ _ __ ___  _ __ ___ | |_ __ _ _   _ ___  ___| |__  \n") +
		ansi.MaybeApply("253", "  `--. \\ __| | '_ ` _ \\| '_ ` _ \\| __/ _` | | | / __|/ __| '_ \\ \n") +
		ansi.MaybeApply("212", " /\\__/ / |_| | | | | | | | | | | | || (_| | |_| \\__ \\ (__| | | |\n") +
		ansi.MaybeApply("111", " \\____/ \\__|_|_| |_| |_|_| |_| |_|\\__\\__,_|\\__,_|___/\\___|_| |_|\n")
}
