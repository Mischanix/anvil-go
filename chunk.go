package anvil

import (
  "bytes"
  "encoding/binary"
  "github.com/Mischanix/nbt-go"
  "io"
)

type Chunk struct {
  X, Z int
  Data nbt.TagCompound
}

func (c *Chunk) read(rd io.Reader) error {
  var f fileChunkHeader
  binary.Read(rd, binary.BigEndian, &f)

  decompressed, err := decompressor(int(f.Compression), rd)
  if err != nil {
    return err
  }
  data, err := nbt.Load(decompressed)
  decompressed.Close()
  if err != nil {
    return err
  }
  c.Data = data
  return nil
}

func (c *Chunk) write(w io.Writer) (int, error) {
  var f fileChunkHeader
  buf := &bytes.Buffer{}
  f.Compression = zlibCompression
  compressed, err := compressor(zlibCompression, buf)
  if err != nil {
    return 0, err
  }
  err = c.Data.Save(compressed)
  compressed.Close()
  if err != nil {
    return 0, err
  }
  b := buf.Bytes()
  f.Length = uint32(len(b))
  err = binary.Write(w, binary.BigEndian, &f)
  if err != nil {
    return 0, err
  }
  _, err = w.Write(b)
  if err != nil {
    return 0, err
  }
  return 5 + int(f.Length), nil
}
