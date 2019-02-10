# Contributing to Stimmtausch

Please abide by the [code of conduct](code-of-conduct.md).

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
* You should run with `ST_ENV=DEV go run main.go --log-level TRACE <world> 2>log.out`. That will use the development config files and show logging at trace level, which will be dumped into `log.out` (which you can `tail -F`).
