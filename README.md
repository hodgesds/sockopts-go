# sockopts-go
Extended socket options for go, currently only full linux support this adds
support for `SO_REUSEPORT` and `TCP_FASTOPEN` to a `net.Listener`.

# Setup
On linux make sure to check `/proc/sys/net/ipv4/tcp_fastopen` is set to `3` to
allow for fast open.
