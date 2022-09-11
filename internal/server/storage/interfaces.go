package storage

type Config interface {
	StorageType() int
	PgConnString() string
	StoreFile() string
	String() string
}
