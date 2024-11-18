package model

import "time"

// Blob represents metadata on the actual file being stored
type Blob struct {
	BlobID           string `gorm:"primaryKey"`
	BlobStoragePath  string
	UploadTime       time.Time
	UserID           string
	FileName         string
	FileSize         int
	FileType         string
	Metadata         []Metadata
	CryptographicKey CryptographicKey `gorm:"foreignKey:KeyID"`
}
