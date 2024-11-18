package model

import "time"

// Blob represents metadata on the actual blob being stored
type Blob struct {
	BlobID           string `gorm:"primaryKey"`
	BlobStoragePath  string
	UploadTime       time.Time
	UserID           string
	BlobName         string
	BlobSize         int
	BlobType         string
	Metadata         []Metadata
	CryptographicKey CryptographicKey `gorm:"foreignKey:KeyID"`
}
