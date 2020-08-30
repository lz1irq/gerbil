# gerbil

A basic CLI for checking whether an IP address is considered blocked by a DNS RBL service.

## Usage

```term
% ./gerbil 180.214.236.4 zen.spamhaus.org
Original IP: 180.214.236.4
DNS RBL domain: zen.spamhaus.org
queryDomain : 4.236.214.180.zen.spamhaus.org
IP 180.214.236.4 blocked: https://www.spamhaus.org/query/ip/180.214.236.4 https://www.spamhaus.org/sbl/query/SBLCSS 
```