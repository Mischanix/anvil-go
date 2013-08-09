package anvil

import (
  "compress/gzip"
  "compress/zlib"
  "errors"
  "fmt"
  "io"
)

const (
  noCompression = iota
  gzipCompression
  zlibCompression
)

func decompressor(compression int, rd io.Reader) (io.ReadCloser, error) {
  if compression == gzipCompression {
    return gzip.NewReader(rd)
  } else if compression == zlibCompression {
    return zlib.NewReader(rd)
  } else {
    return nil, errors.New(fmt.Sprintf("invalid compression type: %d", compression))
  }
}

func compressor(compression int, w io.Writer) (io.WriteCloser, error) {
  if compression == gzipCompression {
    return gzip.NewWriter(w), nil
  } else if compression == zlibCompression {
    return zlib.NewWriterLevel(w, 9)
  } else {
    return nil, errors.New("invalid compression type")
  }
}
