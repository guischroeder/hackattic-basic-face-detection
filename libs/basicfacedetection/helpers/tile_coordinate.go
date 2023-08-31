package helpers

import "math"

func TileCoordinate(value float64) int {
  position := math.Floor(value / 100)

  if v := int(value) % 100; v == 0 {
    position = position - 1
  }

  return int(position)
}
