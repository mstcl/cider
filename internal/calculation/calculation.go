package calculation

import (
	"fmt"
	"math"
	"net/netip"
)

const (
	a   = 0b11111111000000000000000000000000
	b   = 0b00000000111111110000000000000000
	c   = 0b00000000000000001111111100000000
	d   = 0b00000000000000000000000011111111
	all = 0b11111111111111111111111111111111
)

var masks = []uint32{a, b, c, d}

type Results struct {
	Netmask   *netip.Addr
	Gateway   *netip.Addr
	Broadcast *netip.Addr
	Count     int
}

func GetResults(prefix netip.Prefix) (*Results, error) {
	base := prefix.Masked().Addr()

	addrOctets, err := prefix.Addr().MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshal binary: %v", err)
	}

	broadcastOctets := getBroadcast(prefix, addrOctets)

	broadcastAddr := netip.Addr{}
	if err := broadcastAddr.UnmarshalBinary(broadcastOctets); err != nil {
		return nil, fmt.Errorf("unmarshal binary: %v", err)
	}

	mask := getMask(prefix)

	netmaskOctets := getNetmask(prefix)
	netmaskAddr := netip.Addr{}
	if err := netmaskAddr.UnmarshalBinary(netmaskOctets); err != nil {
		return nil, fmt.Errorf("unmarshal binary: %v", err)
	}

	return &Results{
		Netmask:   &netmaskAddr,
		Count:     int(mask + 1),
		Gateway:   &base,
		Broadcast: &broadcastAddr,
	}, nil
}

func getNetmask(prefix netip.Prefix) []byte {
	mask := getMask(prefix)
	remainder := all ^ mask

	netmaskOctets := make([]byte, 4)
	for i := 0; i < 4; i++ {
		l := (remainder & masks[i])
		netmaskOctets[i] = byte(l >> ((3 - i) * 8))
	}

	return netmaskOctets
}

func getBroadcast(prefix netip.Prefix, addrOctets []byte) []byte {
	mask := getMask(prefix)
	broadcastOctets := make([]byte, 4)
	for i, m := range masks {
		maskOctet := uint32(mask&m) >> ((3 - i) * 8)
		broadcastOctets[i] = byte(maskOctet | uint32(addrOctets[i]))
	}

	return broadcastOctets
}

func getMask(prefix netip.Prefix) uint32 {
	endBits := 32 - prefix.Bits()

	mask := uint32(0)
	for i := 0; i < endBits; i++ {
		mask += uint32(math.Pow(2, float64(i)))
	}

	return mask
}
