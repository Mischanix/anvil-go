// Package anvil provides functions for accessing and saving nbt chunk data in
// the Minecraft anvil format.
package anvil

// New loads a minecraft level with path.  The path should be the region folder
// containing r.*.*.mca files.
func New(path string) *Level {
  return &Level{path, make(map[string]*region)}
}

type Level struct {
  path    string
  regions map[string]*region // filename -> region
}

// Chunk returns the chunk at x, z
func (l *Level) Chunk(x, z int) (*Chunk, error) {
  region, err := l.getOrLoadRegion(x, z)
  if err != nil {
    return nil, err
  }
  return region.chunkAt(x, z)
}

// SetChunk sets the chunk at x, z to c.  TODO: NewChunk(x, z int) *Chunk for
// a chunk with NBT skeleton in Data.
func (l *Level) SetChunk(x, z int, c *Chunk) error {
  region, err := l.getOrLoadRegion(x, z)
  if err != nil {
    return err
  }
  return region.saveChunk(x, z, c)
}

// ResetChunk frees the sectors in the corresponding header for the chunk at x,
// z.  Only the region header is modified.
func (l *Level) ResetChunk(x, z int) error {
  region, err := l.getOrLoadRegion(x, z)
  if err != nil {
    return err
  }
  return region.resetChunk(x, z)
}

// getOrLoadRegion gets the region file for the chunk at x, z, loading the file
// if necessary.
func (l *Level) getOrLoadRegion(x, z int) (*region, error) {
  regionFileName := chunkRegionFileName(x, z)
  if region, ok := l.regions[regionFileName]; ok {
    return region, nil
  } else {
    region, err := loadRegion(l.path + regionFileName)
    if err != nil {
      return nil, err
    }
    l.regions[regionFileName] = region
    return region, nil
  }
}
