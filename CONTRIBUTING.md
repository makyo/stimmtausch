# Contributing to Stimmtausch

So you want to contribute to Stimmtausch! Yay!

* Please abide by the [code of conduct](code-of-conduct.md).
* Consider picking something from the [current milestone][currms]/[current project][currproj]

## Some notes

* Comments are good, but logging is often better. Use `log.Tracef` liberally.
* On that note, logging levels:
    * `CRITICAL`: we cannot continue; the program will exit as gracefully as it can.
    * `ERROR`: we cannot continue as requested, but won't quit; maybe you can fix it?
    * `WARNING`: we may not continue as expected, but things may still work; be careful.
    * `INFO`: just letting you know.
    * `DEBUG`: we are doing this thing that's probably important to devs, but not to users.
    * `TRACE`: we got to this point in execution, maybe have some info.
* Exit codes:
    * 0: success
    * 1: misconfigured, non-recoverable, user should fix it
    * 2: code failed, non-recoverable, devs should fix it

## Running

* Run with `go run cmd/st/main.go` - you can watch the logs with `tail -F ~/.local/log/stimmtausch/stimmtausch.log`
* You'll need to create a config file for your worlds, currently. This can live in `~/.strc` or `~/.config/stimmtausch/<whatever>.st.yaml`. You might want to set `stimmtausch.client.syslog.level` (in YAML-ese) to `DEBUG` or even `TRACE`

[currms]: https://github.com/makyo/stimmtausch/milestone/2
[currproj]: https://github.com/makyo/stimmtausch/projects/2
