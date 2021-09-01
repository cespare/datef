# datef

Parse or print timestamps.

Because I'm tired of reading the man page every damn time I use date(1).

## Install

    go install github.com/cespare/datef@latest

## Examples

Print the current time (defaults to RFC3339):

    $ datef
    2015-09-24T17:57:01-07:00

Print the current time as seconds since epoch:

    $ datef -o unix
    1443142623

Parse a timestamp (defaults to unix):

    $ datef 1443142623
    2015-09-24T17:57:03-07:00

Parse an RFC3339 timestamp:

    $ datef -i RFC3339 -o unix 2015-09-24T17:57:03-07:00
    1443142623

Use your own crazy format ([reference](https://golang.org/pkg/time/)):

    $ datef -o 'Jan 1, 2006 at 3:04pm'
    Sep 9, 2015 at 6:10pm

Parse multiple timestamps at once:

    $ datef 1443022623 1443092423 1443142223
    2015-09-23T08:37:03-07:00
    2015-09-24T04:00:23-07:00
    2015-09-24T17:50:23-07:00

Take input from stdin:

    $ datef -i unixms -
    1443022623319
    2015-09-23T08:37:03-07:00
    1443092423411
    2015-09-24T04:00:23-07:00
    1443142223370
    2015-09-24T17:50:23-07:00
    ^D

See `datef -h` for documentation.
