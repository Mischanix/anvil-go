package anvil

import (
  "fmt"
  "math/rand"
  "testing"
)

func ExampleNew(t *testing.T) {
  l := New("world/region/")
  c, err := l.Chunk(0, 0)
  if err != nil {
    t.Errorf(err.Error())
  }

  sections := c.Data.Path("Level/Sections").List()
  for i := int32(0); i < sections.Length(); i++ {
    section := sections.At(i).Compound()
    section.Set("Add", make([]int8, 2048))
    section.Set("Data", make([]int8, 2048))
    blocks := make([]int8, 4096)
    for i, _ := range blocks {
      blocks[i] = int8(rand.Intn(6))
    }
    section.Set("Blocks", blocks)
  }

  r, _ := l.regions["r.0.0.mca"]
  fmt.Println(r.saveChunk(0, 0, c))
}
