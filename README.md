# DNSync
DNSync is a simple meta-DNS service helping to provide name servers with new zones sent in via NOTIFY from
remote servers. It will listen to NOTYIFY DNS queries, operate on them as necessary, and reply accordingly.

## Supported DNS servers
Currently only the BIND name server is supported.

## Installation
Since this is a Go application, deployment is rather easy:

1. Checkout the project somewhere in your GOPATH.
2. Compile the project for the target architecture you need.
3. Copy the resulting file `dnsync` and an adjusted version of its config file `dnsync.json` to where ever
    you need the tool.
4. Start dnsync or configure it as a system service.

## Todo
What still needs to be done:

* More nameserver implementations
* Actually dismiss data when it does not come from a known source
