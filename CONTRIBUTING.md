# Contributing to Stimmtausch

Please abide by the [code of conduct](code-of-conduct.md).

## Some notes

* Comments are good, but logging is often better. Use `log.Tracef` liberally.
* Exit codes:
    * 0: success
    * 1: misconfigured, non-recoverable, user should fix it
    * 2: code failed, non-recoverable, devs should fix it
* You should run with `ST_ENV=DEV go run main.go --log-level TRACE <world> 2>log.out`. That will use the development config files and show logging at trace level, which will be dumped into `log.out` (which you can `tail -F`).
