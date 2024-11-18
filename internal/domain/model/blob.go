package model

import "time"

// Blob represents metadata on the actual blob being stored
type Blob struct {
	BlobID              string `gorm:"primaryKey"`
	BlobStoragePath     string
	UploadTime          time.Time
	UserID              string
	BlobName            string
	BlobSize            int
	BlobType            string
	EncryptionAlgorithm string
	HashAlgorithm       string
	IsEncrypted         bool
	IsSigned            bool
	CryptographicKey    CryptographicKey `gorm:"foreignKey:KeyID"`
}
