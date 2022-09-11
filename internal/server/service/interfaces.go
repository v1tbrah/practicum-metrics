package service

import "time"

type Config interface {
	StorageType() int
	StoreInterval() time.Duration
	StoreFile() string
	Restore() bool
	String() string
}
