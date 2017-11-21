package sockopts

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// SO_REUSEPORT is the socket option to reuse socket port.
const SO_REUSEPORT int = 0x0F

// TCP_FASTOPEN is the socket option to open a TCP fast.
const TCP_FASTOPEN int = 0x17

// LISTEN_BACKLOG is the socket listen backlog.
const LISTEN_BACKLOG int = 23

func allowFastOpen() error {
	b, err := ioutil.ReadFile("/proc/sys/net/ipv4/tcp_fack")
	if err != nil {
		return err
	}
	allowed, err = strconv.Atoi(strings.Replace(string(b), "\n", "", -1))
	if err != nil {
		return err
	}

	if allowed != 3 {
		return fmt.Errorf("set /proc/sys/net/ipv4/tcp_fastopen to 3")
	}

	return nil
}
