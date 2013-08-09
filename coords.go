package anvil

import (
  "fmt"
)

// takes absolute chunk coords x, z and returns their index in a region file's
// header
func chunkRegionIndex(x, z int) int {
  x %= 32
  z %= 32
  if x < 0 {
    x = 31 - -x
  }
  if z < 0 {
    z = 31 - -z
  }
  return z*32 + x
}

// gets the filename (*.mca) that contains chunk with coords x, z
func chunkRegionFileName(x, z int) string {
  return fmt.Sprintf("r.%d.%d.mca", x>>5, z>>5)
}

// BlockChunkCoords converts block coordinates to chunk and section coordinates
func BlockChunkCoords(x, y, z int) (chunkX, sectionY, chunkZ int) {
  chunkX = x >> 4
  chunkZ = z >> 4
  sectionY = y >> 4
  return chunkX, sectionY, chunkZ
}
