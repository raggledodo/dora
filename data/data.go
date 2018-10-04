package data

import "github.com/raggledodo/dora/proto"

type Filter struct {
	Names      []string
	NcaseLimit uint32
}

func ListReqToFilter(req *proto.ListRequest) *Filter {
	return &Filter{
		Names:      req.Names,
		NcaseLimit: req.NcaseLimit,
	}
}

type Database interface {
	ListTestcases(filter *Filter) (map[string]*proto.GeneratedTest, error)
	AddTestcase(tname string, gcase *proto.GeneratedCase) error
	RemoveTestcases(tnames []string) error
}
