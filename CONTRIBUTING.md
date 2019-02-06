# Contributing to Stimmtausch

Please abide by the [code of conduct](code-of-conduct.md).

## Some notes

* Comments are good, but logging is often better. Use `log.Tracef` liberally.
* Exit codes:
    * 0: success
    * 1: misconfigured, non-recoverable, user should fix it
    * 2: code failed, non-recoverable, devs should fix it
