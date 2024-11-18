package model

import "time"

// Metadata represents metadata related to a file (BLOB)
type Metadata struct {
	MetadataID       string `gorm:"primaryKey"`
	BlobID           string
	KeyID            string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	EncryptionAlg    string
	ContentType      string
	Size             int
	UserID           string
	CryptographicKey CryptographicKey `gorm:"foreignKey:KeyID"`
	Blob             Blob             `gorm:"foreignKey:BlobID"`
}
