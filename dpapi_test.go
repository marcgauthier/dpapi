package dpapi

import (
	"bytes"
	"os"
	"testing"
)

func TestStartAndEnvironment(t *testing.T) {
	if err := Start(true); err != nil {
		t.Skip("Skipping test as the environment is not Windows: ", err)
	}
}

func TestEncryptDecryptDPAPI(t *testing.T) {
	// 1. Test with non-empty input.
	originalData := []byte("Hello DPAPI!")
	encryptedData, errEncrypt := EncryptDPAPI(originalData)
	if errEncrypt != nil {
		t.Fatalf("failed to encrypt data: %v", errEncrypt)
	}
	if len(encryptedData) == 0 {
		t.Fatalf("expected encrypted data to be non-empty")
	}
	if bytes.Equal(originalData, encryptedData) {
		t.Fatalf("expected encrypted data to differ from original data")
	}

	decryptedData, errDecrypt := DecryptDPAPI(encryptedData)
	if errDecrypt != nil {
		t.Fatalf("failed to decrypt data: %v", errDecrypt)
	}
	if !bytes.Equal(originalData, decryptedData) {
		t.Errorf("decrypted data does not match original.\nGot:  %v\nWant: %v", decryptedData, originalData)
	}

	// 2. Test with empty input.
	_, errEncrypt = EncryptDPAPI([]byte{})
	if errEncrypt == nil {
		t.Error("expected an error when encrypting empty data, but got nil")
	}

	_, errDecrypt = DecryptDPAPI([]byte{})
	if errDecrypt == nil {
		t.Error("expected an error when decrypting empty data, but got nil")
	}
}

func TestEncryptDecryptFile(t *testing.T) {
	// Create a temporary file for testing.
	inputFile, err := os.CreateTemp("", "testfile*.txt")
	if err != nil {
		t.Fatalf("failed to create temporary input file: %v", err)
	}
	defer os.Remove(inputFile.Name())

	// Write some test data to the file.
	originalData := []byte("Hello File Encryption!")
	if _, err := inputFile.Write(originalData); err != nil {
		t.Fatalf("failed to write to temporary input file: %v", err)
	}
	inputFile.Close()

	// Encrypt the file and write to a new output file.
	outputFile := inputFile.Name() + ".enc"
	if err := EncryptFile(inputFile.Name(), outputFile); err != nil {
		t.Fatalf("failed to encrypt file: %v", err)
	}
	defer os.Remove(outputFile)

	// Verify that the output file exists and is non-empty.
	encryptedData, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read encrypted file: %v", err)
	}
	if len(encryptedData) == 0 {
		t.Fatal("expected encrypted file to be non-empty")
	}
	if bytes.Equal(originalData, encryptedData) {
		t.Fatal("expected encrypted file data to differ from original data")
	}

	// Decrypt the file back to the original input file.
	if err := DecryptFile(outputFile, inputFile.Name()); err != nil {
		t.Fatalf("failed to decrypt file: %v", err)
	}

	// Verify that the decrypted file matches the original data.
	decryptedData, err := os.ReadFile(inputFile.Name())
	if err != nil {
		t.Fatalf("failed to read decrypted file: %v", err)
	}
	if !bytes.Equal(originalData, decryptedData) {
		t.Errorf("decrypted file data does not match original.\nGot:  %s\nWant: %s", decryptedData, originalData)
	}
}

func TestEncryptDecryptInvalidCases(t *testing.T) {
	// Test decryption of invalid data
	invalidData := []byte("this is not encrypted")
	if _, err := DecryptDPAPI(invalidData); err == nil {
		t.Error("expected error when decrypting invalid data, but got nil")
	}

	// Test encryption and decryption with non-existent file
	nonExistentFile := "nonexistentfile.txt"
	if err := EncryptFile(nonExistentFile, ""); err == nil {
		t.Error("expected error when encrypting non-existent file, but got nil")
	}
	if err := DecryptFile(nonExistentFile, ""); err == nil {
		t.Error("expected error when decrypting non-existent file, but got nil")
	}
}

func TestEncryptDecryptFileWithVariousOutputs(t *testing.T) {
	// Create a temporary file for testing
	inputFile, err := os.CreateTemp("", "testfile*.txt")
	if err != nil {
		t.Fatalf("failed to create temporary input file: %v", err)
	}
	defer os.Remove(inputFile.Name())

	// Write some test data to the file
	originalData := []byte("Hello File Encryption!")
	if _, err := inputFile.Write(originalData); err != nil {
		t.Fatalf("failed to write to temporary input file: %v", err)
	}
	inputFile.Close()

	// Test encryption and decryption with a specified output file
	outputFile := inputFile.Name() + ".enc"
	if err := EncryptFile(inputFile.Name(), outputFile); err != nil {
		t.Fatalf("failed to encrypt file: %v", err)
	}
	defer os.Remove(outputFile)

	encryptedData, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read encrypted file: %v", err)
	}
	if bytes.Equal(originalData, encryptedData) {
		t.Fatal("expected encrypted file data to differ from original data")
	}

	if err := DecryptFile(outputFile, inputFile.Name()); err != nil {
		t.Fatalf("failed to decrypt file: %v", err)
	}

	decryptedData, err := os.ReadFile(inputFile.Name())
	if err != nil {
		t.Fatalf("failed to read decrypted file: %v", err)
	}
	if !bytes.Equal(originalData, decryptedData) {
		t.Errorf("decrypted file data does not match original.\nGot: %s\nWant: %s", decryptedData, originalData)
	}

	// Test encryption and decryption with an empty output file (overwrite behavior)
	if err := EncryptFile(inputFile.Name(), ""); err != nil {
		t.Fatalf("failed to encrypt file with empty outputFile: %v", err)
	}
	if err := DecryptFile(inputFile.Name(), ""); err != nil {
		t.Fatalf("failed to decrypt file with empty outputFile: %v", err)
	}

	// Verify the overwritten file matches original data after decryption
	decryptedData, err = os.ReadFile(inputFile.Name())
	if err != nil {
		t.Fatalf("failed to read decrypted file: %v", err)
	}
	if !bytes.Equal(originalData, decryptedData) {
		t.Errorf("decrypted file data does not match original after overwrite.\nGot: %s\nWant: %s", decryptedData, originalData)
	}
}
