package anvil

import (
  "bytes"
  "encoding/binary"
  "errors"
  "os"
  "time"
)

type region struct {
  f *os.File
  h *fileHeader
}

func loadRegion(path string) (*region, error) {
  var result region
  f, err := os.OpenFile(path, os.O_RDWR, 0)
  if err != nil {
    return nil, err
  }
  result.f = f
  result.h = &fileHeader{}
  binary.Read(result.f, binary.BigEndian, result.h)
  return &result, nil
}

func (r *region) chunkAt(x, z int) (*Chunk, error) {
  if r.f == nil {
    return nil, errors.New("file for region is nil")
  }
  loc := r.h.Locations[chunkRegionIndex(x, z)]
  if _, err := r.f.Seek(int64(loc.offset()), 0); err != nil {
    return nil, err
  }
  c := &Chunk{x, z, nil}
  if err := c.read(r.f); err != nil {
    return nil, err
  } else {
    return c, nil
  }
}

func (r *region) saveChunk(x, z int, c *Chunk) error {
  loc := r.h.Locations[chunkRegionIndex(x, z)]
  buf := new(bytes.Buffer)
  l, err := c.write(buf)
  if err != nil {
    return err
  }
  sectorCount := sectorCount(l)
  if sectorCount <= int(loc.SectorCount) {
    loc.setSectorCount(l)
    if _, err = r.f.Seek(int64(loc.offset()), 0); err != nil {
      return err
    }
    if _, err = r.f.Write(buf.Bytes()); err != nil {
      return err
    }
  } else {
    allocations := r.allocations()
    for i := 0; i < int(loc.SectorCount); i++ {
      allocations[i+int(loc.sector())] = -1
    }
    targetSector := -1
    for i, exists := range allocations {
      if exists <= 0 {
        valid := true
        for j := 1; j < sectorCount; j++ {
          if allocations[j+i] > 0 {
            valid = false
            break
          }
        }
        if valid {
          targetSector = i
          break
        }
      }
    }
    if targetSector < 0 {
      targetSector = len(allocations)
      loc.setOffset(targetSector)
      buf.Write(make([]byte, sectorCount*4*1024-l))
    } else {
      loc.setOffset(targetSector)
    }

    if _, err = r.f.Seek(int64(4*1024*targetSector), 0); err != nil {
      return err
    }
    if _, err = r.f.Write(buf.Bytes()); err != nil {
      return err
    }
    loc.setSectorCount(l)
  }
  r.h.Locations[chunkRegionIndex(x, z)] = loc
  r.updateTimestamp(x, z)
  return r.writeHeader()
}

func (r *region) allocations() []int {
  allocationList := []int{2, 2}
  for _, s := range r.h.Locations {
    for i := 0; i < int(s.SectorCount); i++ {
      sector := int(s.sector()) + i
      if sector >= len(allocationList) {
        allocationList = append(
          allocationList,
          make([]int, 1+sector-len(allocationList))...,
        )
      }
      allocationList[sector] += 1
    }
  }
  return allocationList
}

func (r *region) resetChunk(x, z int) error {
  loc := &r.h.Locations[chunkRegionIndex(x, z)]
  loc.Sector = [3]uint8{0, 0, 0}
  loc.SectorCount = 0
  r.h.Timestamps[chunkRegionIndex(x, z)] = 0
  return r.writeHeader()
}

func (r *region) updateTimestamp(x, z int) {
  r.h.Timestamps[chunkRegionIndex(x, z)] = uint32(time.Now().Unix())
}

func (r *region) writeHeader() error {
  if _, err := r.f.Seek(0, 0); err != nil {
    return err
  }
  if err := binary.Write(r.f, binary.BigEndian, r.h); err != nil {
    return err
  }
  return nil
}
