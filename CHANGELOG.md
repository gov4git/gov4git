# v1.1.10

## Changes

* the organizer can approve join requests by replying to github issues
* the organizer can send directives to the governance system by writing github issues
* add license apache/mit

# v1.1.9

## Context

* by default, go-git uses an external git binary for git file urls
* gov4git uses git file urls only when caching is enabled (in the config)

## Changes

* use go-git native git file url handling on linux, darwin and windows
* use cache during unit tests only on linux and darwin (not windows)
* use cache during api tests on linux, darwin and windows

## Notes

* caching does not pass ci unit tests on windows, but it does pass ci api tests on windows
* there is no explanation for this discrepancy currently
* users can use caching on windows at their own risk
