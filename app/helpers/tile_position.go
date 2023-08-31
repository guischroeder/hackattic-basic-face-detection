package helpers

import "math"

func TilePosition(value float64, imgDimensionSize float64) int {
  numberOfTiles := 8
  tileSize := imgDimensionSize / float64(numberOfTiles)

  position := math.Floor(value / tileSize)

  if v := int(value) % 100; v == 0 {
    position = position - 1
  }

  return int(position)
}
