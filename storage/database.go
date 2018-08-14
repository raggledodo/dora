package storage

import (
	"github.com/raggledodo/dora/proto"
)

type Database interface {
	AddTestResult(string, *proto.TestOutput) error
	ListTestcases(read func(string)error) error
	GetTestResults(string) (chan *proto.TestOutput, error)
}
