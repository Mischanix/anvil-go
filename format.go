package anvil

const chunksPerRegion = 32 * 32

type fileHeader struct {
  Locations  [chunksPerRegion]fileLocation
  Timestamps [chunksPerRegion]uint32
}

type fileLocation struct {
  Sector      [3]uint8
  SectorCount uint8
}

func (l *fileLocation) sector() uint {
  return uint(l.Sector[2]) | uint(l.Sector[1])<<8 | uint(l.Sector[0])<<16
}

func (l *fileLocation) offset() int {
  return int(4 * 1024 * l.sector())
}

func (l *fileLocation) setOffset(s int) {
  l.Sector[0] = uint8(s & 0xff0000 >> 16)
  l.Sector[1] = uint8(s & 0xff00 >> 8)
  l.Sector[2] = uint8(s & 0xff)
}

func (l *fileLocation) setSectorCount(dataLength int) {
  l.SectorCount = uint8(sectorCount(dataLength))
}

func sectorCount(dataLength int) int {
  return 1 + (5+dataLength)/(4*1024)
}

type fileChunkHeader struct {
  Length      uint32 // completely useless
  Compression uint8  // slightly useless
}
