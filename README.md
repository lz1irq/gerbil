# gerbil

A basic CLI for checking whether an IP address is considered blocked by a DNS RBL service.

## Usage

```term
go build && ./gerbil -rbl.file rbl.list -host 180.214.236.4
OK: 180.214.236.4 at zen.spamhaus.org, dnsbl.uni-sofia.bg,
```