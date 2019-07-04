package main

import log "github.com/sirupsen/logrus"

const (
  DefaultHeaderSize     = 1024 * 4
  DefaultPageSize       = 1024 * 4
  DefaultPageCount      = 1024
  DefaultMaxKeySize     = 256
  DefaultPageHeaderSize = 64
)

func init() {
  log.SetFormatter(&log.TextFormatter{
    FullTimestamp: true,
  })
  // log.SetReportCaller(true)
}

// Paged Config
type Config struct {
  headerSize     int16 // header size
  pageSize       int32 // page size
  pageCount      int64 // init page count
  maxKeySize     int16 // max key size
  pageHeaderSize int8  // page header size
}

// Default paged config
func DefaultConfig() Config {
  return Config{
    headerSize:     DefaultHeaderSize,
    pageSize:       DefaultPageSize,
    pageCount:      DefaultPageCount,
    maxKeySize:     DefaultMaxKeySize,
    pageHeaderSize: DefaultPageHeaderSize,
  }
}
