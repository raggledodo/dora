package data

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/raggledodo/dora/proto"
)

type TimeRange struct {
	After time.Time
	Until time.Time
}

type Filter struct {
	TestNames   []string
	TestBetween *TimeRange
}

func (tr TimeRange) IsBetween(ts *timestamp.Timestamp) bool {
	if ts == nil {
		return false
	}
	t := time.Unix(ts.Seconds, int64(ts.Nanos)).UTC()
	return tr.After.Before(t) && tr.Until.After(t)
}

func ListReqToFilter(req *proto.ListRequest) (*Filter, error) {
	var between *TimeRange
	if req.TestsAfter != nil && req.TestsUntil != nil {
		tAfter := req.TestsAfter
		tUntil := req.TestsUntil
		after := time.Unix(tAfter.Seconds, int64(tAfter.Nanos)).UTC()
		until := time.Unix(tUntil.Seconds, int64(tUntil.Nanos)).UTC()
		if after.After(until) {
			return nil, fmt.Errorf(
				"tests_after of value (%s) must be tests_until (%s)",
				after, until)
		}
		between = &TimeRange{After: after, Until: until}
	}
	return &Filter{TestNames: req.TestNames, TestBetween: between}, nil
}

type Database interface {
	ListTestcases(filter *Filter) (map[string]*proto.GeneratedTest, error)
	AddTestcase(tname string, gcase *proto.GeneratedCase) error
	RemoveTestcases(tnames []string) error
}
