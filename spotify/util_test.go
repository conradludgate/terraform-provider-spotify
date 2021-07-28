package spotify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBatches(t *testing.T) {
	ranges := batches(521, 100)
	assert.Equal(t, []Range{
		{Start: 0, End: 100},
		{Start: 100, End: 200},
		{Start: 200, End: 300},
		{Start: 300, End: 400},
		{Start: 400, End: 500},
		{Start: 500, End: 521},
	}, ranges)
}

func TestBatchesFull(t *testing.T) {
	ranges := batches(100, 100)
	assert.Equal(t, []Range{
		{Start: 0, End: 100},
	}, ranges)
}

func TestBatchesPartial(t *testing.T) {
	ranges := batches(50, 100)
	assert.Equal(t, []Range{
		{Start: 0, End: 50},
	}, ranges)
}

func TestBatchesEmpty(t *testing.T) {
	ranges := batches(0, 100)
	assert.Empty(t, ranges)
}
