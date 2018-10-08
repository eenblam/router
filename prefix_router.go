package router

type PrefixRouter struct {
	PrefixTree *PrefixTree
	Routes     map[IPv4]([]*PrefixTree)
}

func (p *PrefixRouter) Add(r Route) {
	p.PrefixTree.Add(r)
}

func (p *PrefixRouter) Drop(r Route) {
	p.PrefixTree.Drop(r)
}

func (p *PrefixRouter) DropAllTo(r IPv4) {
	nodes, found := p.Routes[r]
	if !found {
		return
	}
	for _, node := range nodes {
		if *node.Route == r {
			node.Route = nil
		}
	}
}

func (p *PrefixRouter) Get(ipv4 IPv4) *IPv4 {
	return p.PrefixTree.Get(ipv4)
}
