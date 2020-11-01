package nic

import (
	"fmt"
	"net"

	"github.com/cs3238-tsuzu/multipath-proxy/pkg/config"
)

// Peer represents an address to listen on and the destination address
type Peer struct {
	Address *net.UDPAddr
	Server  string
}

// GetPeers resolves peers from config
func GetPeers(peers []config.Peer) ([]*Peer, error) {
	nic2addr, err := getNIC2Addrs()

	if err != nil {
		return nil, fmt.Errorf("failed to get map of nic2addrs: %w", err)
	}

	ret := make([]*Peer, 0, len(peers))

	for i := range peers {
		peer := &peers[i]
		var addr *net.UDPAddr
		var err error

		if peer.Listen != "" {
			addr, err = net.ResolveUDPAddr("", peer.Listen)

			if err != nil {
				return nil, fmt.Errorf("failed to parse '%s': %w", peer.Listen, err)
			}
		}

		if peer.NIC == "" {
			ret = append(ret, &Peer{
				Server:  peer.Server,
				Address: addr,
			})
		} else {
			addrs, ok := nic2addr[peer.NIC]

			if !ok {
				return nil, fmt.Errorf("NIC '%s' is not found", peer.NIC)
			}

			for i := range addrs {
				a := &net.UDPAddr{
					IP: addrs[i],
				}
				if addr != nil {
					a.Port = addr.Port
					a.Zone = addr.Zone
				}

				ret = append(ret, &Peer{
					Server:  peer.Server,
					Address: a,
				})
			}
		}
	}

	return ret, nil
}

func getNIC2Addrs() (map[string][]net.IP, error) {
	nics, err := net.Interfaces()

	if err != nil {
		return nil, fmt.Errorf("failed to list NICs: %w", err)
	}

	nic2addr := make(map[string][]net.IP)
	for i := range nics {
		addrs, err := nics[i].Addrs()

		if err != nil {
			// Ignore the NIC
			continue
		}

		saddrs := make([]net.IP, len(addrs))

		for i := range addrs {
			saddrs[i] = addrs[i].(*net.IPNet).IP
		}

		nic2addr[nics[i].Name] = saddrs
	}

	return nic2addr, nil
}
