package storage

import (
	"github.com/raggledodo/dora/proto"
)

type Database interface {
	ListTestcases(read func(string)error) error
	GetTestResults(string) (chan *proto.GeneratedCase, error)

	AddTestResult(string, *proto.GeneratedCase) error
	RemoveTestResult(string) error
}
