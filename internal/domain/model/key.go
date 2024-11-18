package model

import "time"

// CryptographicKey represents the encryption key entity
type CryptographicKey struct {
	KeyID     string `gorm:"primaryKey"`
	KeyValue  string
	KeyType   string
	CreatedAt time.Time
	ExpiresAt time.Time
	UserID    string `gorm:"index"`
	Metadata  []Metadata
	Blobs     []Blob
}
