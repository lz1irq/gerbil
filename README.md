# gerbil

A basic CLI for checking whether an IP address is considered blocked by a DNS RBL service.

## Usage

```term
% go build && ./gerbil -ip 180.214.236.4 -rbl.file rbl.list                                                                   :(
rbl=zen.spamhaus.org, ip=180.214.236.4, blocked=true, reason=https://www.spamhaus.org/query/ip/180.214.236.4 https://www.spamhaus.org/sbl/query/SBLCSS 
rbl=dnsbl.uni-sofia.bg, ip=180.214.236.4, blocked=false, reason=
```