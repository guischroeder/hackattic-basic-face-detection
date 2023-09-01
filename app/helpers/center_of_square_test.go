package helpers

import (
    "testing"

	"github.com/stretchr/testify/assert"
) 

func TestCenterOfSquare(t *testing.T) {
    result := CenterOfSquare(10.0, 10.0)

    assert.Equal(t, result, 15.0)
}
