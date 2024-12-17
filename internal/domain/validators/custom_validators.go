package validators

import "github.com/go-playground/validator/v10"

func KeySizeValidation(fl validator.FieldLevel) bool {
	algorithm := fl.Parent().FieldByName("Algorithm").String()
	keySize := fl.Field().Uint()

	switch algorithm {
	case "AES":
		// AES key sizes should be 128, 192, or 256
		return keySize == 128 || keySize == 192 || keySize == 256
	case "RSA":
		// RSA key sizes should be at least 512
		return keySize >= 512 && keySize <= 4096
	case "EC":
		// EC key sizes can be 256, 384, or 521
		return keySize == 256 || keySize == 384 || keySize == 521
	default:
		return false
	}
}
