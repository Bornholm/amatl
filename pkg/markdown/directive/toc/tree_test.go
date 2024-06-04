package toc

import (
	"testing"
)

func TestTree(t *testing.T) {
	trunk := NewBranch[*tocItem]()
	level1 := NewBranch[*tocItem]()

	level2 := NewBranch[*tocItem]()
	level2bis := NewBranch[*tocItem]()

	trunk.Add(level1)
	level1.Add(level2)
	level1.Add(level2bis)
}
