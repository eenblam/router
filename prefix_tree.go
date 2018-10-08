package router

import (
	"fmt"
	"log"
	"strings"
)

// A PrefixTree is a binary trie providing efficient lookups of IPv4 addresses.
type PrefixTree struct {
	One  *PrefixTree
	Zero *PrefixTree
	// For now, assume only one route per target
	Route *IPv4
}

func NewPrefixTree() *PrefixTree {
	return &PrefixTree{nil, nil, nil}
}

// Add inserts a new route into the trie.
func (p *PrefixTree) Add(r Route) {
	var pastByteBits, i, j uint8
	current := p
	for i = 0; i < 4; i++ {
		pastByteBits = i * 8
		for j = 0; j < 8; j++ {
			if r.Prefix == (j + pastByteBits) {
				// Done, store route here and exit
				if current.Route != nil {
					log.Printf("Dropping route %v for %v", current.Route, r.To)
				}
				current.Route = r.To
				return
			}
			// Iterate over byte bits from left
			if 1 == 1&(r.Masked[i]>>(7-j)) {
				// Step deeper as 1
				if current.One == nil {
					current.One = NewPrefixTree()
				}
				current = current.One
			} else {
				// Step deeper as 0
				if current.Zero == nil {
					current.Zero = NewPrefixTree()
				}
				current = current.Zero
			}
		}
	}
}

// Drop removes a route from the PrefixTree router.
// In this case, the drop is soft:
// it only removes a route from a subnode,
// and it does not prune subnodes that no longer lead to a route.
func (p *PrefixTree) Drop(r Route) {
	var pastByteBits, i, j uint8
	current := p
	for i = 0; i < 4; i++ {
		pastByteBits = i * 8
		for j = 0; j < 8; j++ {
			if r.Prefix == j+pastByteBits {
				// Found it. Drop route, whether it exists or not.
				current.Route = nil
				return
			}
			// Iterate over byte bits from left
			if 1 == (r.Masked[i] >> (7 - j)) {
				// Step deeper as 1
				if current.One == nil {
					// Don't have it anyway
					return
				}
				current = current.One
			} else {
				// Step deeper as 0
				if current.Zero == nil {
					// Don't have it anyway
					return
				}
				current = current.Zero
			}
		}
	}
}

// Get retrieves a route from the table, returning a pointer to the
// destination if found. Pointer is nil if not found.
func (p *PrefixTree) Get(ipv4 IPv4) *IPv4 {
	var lastBest *IPv4
	var b byte
	current := p
	var i, j uint8
	for i = 0; i < 4; i++ {
		b = ipv4[i]
		for j = 0; j < 8; j++ {
			if current.Route != nil {
				lastBest = current.Route
			}
			// Iterate over byte bits from left
			if 1 == 1&(b>>(7-j)) {
				// Step deeper as 1
				if current.One == nil {
					return lastBest
				}
				current = current.One
			} else {
				// Step deeper as 0
				if current.Zero == nil {
					return lastBest
				}
				current = current.Zero
			}
		}
	}
	return lastBest
}

// String provides a simple debug string for the trie.
func (p *PrefixTree) String() string {
	zs, os, rs := "", "", ""
	if p.Zero != nil {
		zs = fmt.Sprintf(" 0: {\n%s\n}, ", p.Zero.string("  "))
	}
	if p.One != nil {
		os = fmt.Sprintf(" 1: {\n%s\n}, ", p.One.string("  "))
	}
	if p.Route != nil {
		rs = fmt.Sprintf(" R: %d.%d.%d.%d", p.Route[0], p.Route[1], p.Route[2], p.Route[3])
	}
	almost := strings.Trim(strings.Join([]string{zs, os, rs}, "\n"), "\n")
	return fmt.Sprintf("{\n%s\n}", almost)
}

// string is used internally by String for recursive printing of subnodes
// in the trie.
func (p *PrefixTree) string(pad string) string {
	zs, os, rs := "", "", ""
	nextPad := fmt.Sprintf("%s  ", pad)
	if p.Zero != nil {
		zs = fmt.Sprintf("%s0: {\n%s\n%s},", pad, p.Zero.string(nextPad), pad)
	}
	if p.One != nil {
		os = fmt.Sprintf("%s1: {\n%s\n%s},", pad, p.One.string(nextPad), pad)
	}
	if p.Route != nil {
		rs = fmt.Sprintf("%sR: %d.%d.%d.%d", pad, p.Route[0], p.Route[1], p.Route[2], p.Route[3])
	}
	return strings.Trim(strings.Join([]string{zs, os, rs}, "\n"), "\n")
}
