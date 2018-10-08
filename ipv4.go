package router

import (
	"fmt"
	"net"
)

type IPv4 [4]byte

func IPv4FromString(address string) (*IPv4, error) {
	addr := net.ParseIP(address)
	if addr == nil || addr.To4() == nil {
		return nil, fmt.Errorf("Address %s is not an IPv4 address", address)
	}
	ipv4 := IPv4{}
	for i := 0; i < 4; i++ {
		ipv4[i] = addr[i]
	}
	return &ipv4, nil
}

// IsMask returns true if the address is a valid subnet mask, false otherwise.
func (i *IPv4) IsMask() bool {
	ones := true
	var j uint
	for _, v := range i {
		if ones {
			if v != 255 {
				ones = false
				// Validate byte as having only leading ones, if any
				for j = 0; j < 8; j++ {
					shifted := (v << j) & 255
					bitIsZero := 128 != (shifted & 128)
					shiftedZero := 0 == shifted
					if bitIsZero {
						if shiftedZero {
							// Done with this byte
							break
						} else {
							// Found a zero bit before zeroing the byte
							return false
						}
					}
				}
			}
		} else {
			if v != 0 {
				return false
			}
		}
	}
	return true
}

func (i *IPv4) MaskWithPrefix(prefix uint8) *IPv4 {
	mask, err := MaskFromPrefix(prefix)
	if err != nil {
		return nil
	}
	return i.MaskWith(*mask)
}

func (i *IPv4) MaskWith(mask IPv4) *IPv4 {
	if !mask.IsMask() {
		return nil
	}
	masked := &IPv4{}
	for j := 0; j < 4; j++ {
		masked[j] = i[j] & mask[j]
	}
	return masked
}

func MaskFromPrefix(prefix uint8) (*IPv4, error) {
	if prefix > 32 {
		return nil, fmt.Errorf("Prefix %d is greater than 32", prefix)
	}
	mask := &IPv4{}
	numFull := prefix / 8
	var i uint8
	for i = 0; i < numFull; i++ {
		mask[i] = 255
	}
	// Set value of last non-zero byte
	if numFull < 4 {
		rem := prefix % 8
		// 256 - (2 ^ (8-rem))
		power := 8 - rem
		partial := 1
		for i = 0; i < power; i++ {
			partial *= 2
		}
		mask[numFull] = byte(256 - partial)
	}
	return mask, nil
}
