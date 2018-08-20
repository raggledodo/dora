package storage

import (
	"github.com/raggledodo/dora/proto"
)

type Database interface {
	ListTestcases(read func(string)error) error
	GetTestResults(string) (chan *proto.TestOutput, error)

	AddTestResult(string, *proto.TestOutput) error
	RemoveTestResult(string) error
}
