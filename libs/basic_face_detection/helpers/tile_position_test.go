package helpers

import (
    "testing"

	"github.com/stretchr/testify/assert"
) 

func TestTilePosition(t *testing.T) {
    result := TilePosition(500.0, 800.0)

    assert.Equal(t, result, 4)
}
