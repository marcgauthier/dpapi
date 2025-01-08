package dpapi

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

var verboseLogging bool // Global variable to control verbose logging

// Start checks if the OS is Windows, sets the verbose logging flag, and returns an error if the OS is not supported.
func Start(verbose bool) error {
	// Set the global verbose logging flag
	verboseLogging = verbose

	if runtime.GOOS != "windows" {
		if verboseLogging {
			log.Println("Unsupported OS detected. This package only works on Windows.")
		}
		return fmt.Errorf("this package only supports Windows")
	}

	if verboseLogging {
		log.Println("Initialization successful. Running on Windows.")
	}
	return nil
}

// EncryptDPAPI encrypts the data using DPAPI in the machine context.
// It returns the encrypted data and any error encountered.
func EncryptDPAPI(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("input data cannot be empty")
	}

	if verboseLogging {
		log.Println("Encrypting data using DPAPI...")
	}
	desc := windows.StringToUTF16Ptr("") // Optional description for the encryption
	inBlob := windows.DataBlob{
		Size: uint32(len(data)),
		Data: &data[0],
	}
	var outBlob windows.DataBlob

	err := windows.CryptProtectData(
		&inBlob,
		desc,
		nil,
		0,
		nil,
		windows.CRYPTPROTECT_LOCAL_MACHINE,
		&outBlob,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(outBlob.Data)))

	// Create a copy of the encrypted data
	encrypted := make([]byte, outBlob.Size)
	copy(encrypted, unsafe.Slice(outBlob.Data, outBlob.Size))

	if verboseLogging {
		log.Println("Data encryption successful.")
	}
	return encrypted, nil
}

// DecryptDPAPI decrypts the data using DPAPI in the machine context.
// It returns the decrypted data and any error encountered.
func DecryptDPAPI(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("input data cannot be empty")
	}

	if verboseLogging {
		log.Println("Decrypting data using DPAPI...")
	}

	inBlob := windows.DataBlob{
		Size: uint32(len(data)),
		Data: &data[0],
	}
	var outBlob windows.DataBlob
	var desc *uint16

	err := windows.CryptUnprotectData(
		&inBlob,
		&desc,
		nil,
		0,
		nil,
		windows.CRYPTPROTECT_LOCAL_MACHINE,
		&outBlob,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(outBlob.Data)))
	if desc != nil {
		defer windows.LocalFree(windows.Handle(unsafe.Pointer(desc)))
	}

	// Create a copy of the decrypted data
	decrypted := make([]byte, outBlob.Size)
	copy(decrypted, unsafe.Slice(outBlob.Data, outBlob.Size))

	if verboseLogging {
		log.Println("Data decryption successful.")
	}
	return decrypted, nil
}

// EncryptFile encrypts the contents of a file and writes the encrypted data to the same file or a new file.
// If outputFile is an empty string, the inputFile is overwritten.
func EncryptFile(inputFile, outputFile string) error {
	if verboseLogging {
		log.Printf("Reading file: %s\n", inputFile)
	}
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	if verboseLogging {
		log.Println("File read successful. Encrypting contents...")
	}
	encryptedData, err := EncryptDPAPI(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt file data: %w", err)
	}

	if outputFile == "" {
		outputFile = inputFile
	}

	if verboseLogging {
		log.Printf("Writing encrypted data to file: %s\n", outputFile)
	}
	err = os.WriteFile(outputFile, encryptedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write encrypted data to file: %w", err)
	}

	if verboseLogging {
		log.Println("File encryption successful.")
	}
	return nil
}

// DecryptFile decrypts the contents of a file and writes the decrypted data to the same file or a new file.
// If outputFile is an empty string, the inputFile is overwritten.
func DecryptFile(inputFile, outputFile string) error {
	if verboseLogging {
		log.Printf("Reading file: %s\n", inputFile)
	}
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	if verboseLogging {
		log.Println("File read successful. Decrypting contents...")
	}
	decryptedData, err := DecryptDPAPI(data)
	if err != nil {
		return fmt.Errorf("failed to decrypt file data: %w", err)
	}

	if outputFile == "" {
		outputFile = inputFile
	}

	if verboseLogging {
		log.Printf("Writing decrypted data to file: %s\n", outputFile)
	}
	err = os.WriteFile(outputFile, decryptedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write decrypted data to file: %w", err)
	}

	if verboseLogging {
		log.Println("File decryption successful.")
	}
	return nil
}
