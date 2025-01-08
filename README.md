# DPAPI Package

The `dpapi` package provides a simple interface for encrypting and decrypting data and files using the Windows Data Protection API (DPAPI). This package is specifically designed for Windows environments and utilizes the `golang.org/x/sys/windows` package to interact with the underlying Windows DPAPI functionality.

## Features

- Encrypt and decrypt data in memory.
- Encrypt and decrypt files.
- Supports verbose logging for detailed debugging.
- Simple and intuitive API for seamless integration.

## Requirements

- Windows OS (the package does not support other operating systems).
- Go 1.18 or later.

## Installation

```sh
go get -u <your-repo-path>/dpapi
```

## Initialization

Before using the `dpapi` package, you must initialize it with the `Start` function. This ensures compatibility with the environment.

```go
if err := dpapi.Start(true); err != nil {
    log.Fatalf("Failed to initialize dpapi: %v", err)
}
```

## Functions

### `Start(verbose bool) error`

Initializes the package and checks if the environment is Windows. Logs detailed messages if verbose logging is enabled.

- **Parameters:**
  - `verbose`: If `true`, enables detailed logging for debugging purposes.
- **Returns:**
  - `error`: An error if the environment is not supported.

Example:
```go
if err := dpapi.Start(true); err != nil {
    log.Fatalf("Failed to initialize dpapi: %v", err)
}
```

### `EncryptDPAPI(data []byte) ([]byte, error)`

Encrypts the given data using DPAPI.

- **Parameters:**
  - `data`: A `[]byte` slice containing the data to be encrypted. Must not be empty.
- **Returns:**
  - `[]byte`: The encrypted data.
  - `error`: An error if encryption fails.

Example:
```go
data := []byte("Hello, World!")
encryptedData, err := dpapi.EncryptDPAPI(data)
if err != nil {
    log.Fatalf("Failed to encrypt data: %v", err)
}
```

### `DecryptDPAPI(data []byte) ([]byte, error)`

Decrypts the given data that was previously encrypted with DPAPI.

- **Parameters:**
  - `data`: A `[]byte` slice containing the encrypted data. Must not be empty.
- **Returns:**
  - `[]byte`: The decrypted data.
  - `error`: An error if decryption fails.

Example:
```go
decryptedData, err := dpapi.DecryptDPAPI(encryptedData)
if err != nil {
    log.Fatalf("Failed to decrypt data: %v", err)
}
```

### `EncryptFile(inputFile, outputFile string) error`

Encrypts the contents of a file and writes the encrypted data to the specified output file. If `outputFile` is empty, the input file is overwritten.

- **Parameters:**
  - `inputFile`: The path to the file to encrypt.
  - `outputFile`: The path to the file where encrypted data should be written. If empty, the input file is overwritten.
- **Returns:**
  - `error`: An error if encryption fails.

Example:
```go
if err := dpapi.EncryptFile("example.txt", "example.enc"); err != nil {
    log.Fatalf("Failed to encrypt file: %v", err)
}
```

### `DecryptFile(inputFile, outputFile string) error`

Decrypts the contents of a file and writes the decrypted data to the specified output file. If `outputFile` is empty, the input file is overwritten.

- **Parameters:**
  - `inputFile`: The path to the file to decrypt.
  - `outputFile`: The path to the file where decrypted data should be written. If empty, the input file is overwritten.
- **Returns:**
  - `error`: An error if decryption fails.

Example:
```go
if err := dpapi.DecryptFile("example.enc", "example.txt"); err != nil {
    log.Fatalf("Failed to decrypt file: %v", err)
}
```

## Examples

### Encrypt and Decrypt Data

```go
package main

import (
	"fmt"
	"log"

	"example/dpapi" // Replace with the actual import path of your dpapi package
)

func main() {
	if err := dpapi.Start(true); err != nil {
		log.Fatalf("Failed to initialize dpapi: %v", err)
	}

	data := []byte("Hello, DPAPI!")
	encryptedData, err := dpapi.EncryptDPAPI(data)
	if err != nil {
		log.Fatalf("Failed to encrypt data: %v", err)
	}
	fmt.Println("Encrypted data:", encryptedData)

	decryptedData, err := dpapi.DecryptDPAPI(encryptedData)
	if err != nil {
		log.Fatalf("Failed to decrypt data: %v", err)
	}
	fmt.Println("Decrypted data:", string(decryptedData))
}
```

### Encrypt and Decrypt Files

```go
package main

import (
	"fmt"
	"log"
	"os"

	"example/dpapi" // Replace with the actual import path of your dpapi package
)

func main() {
	if err := dpapi.Start(true); err != nil {
		log.Fatalf("Failed to initialize dpapi: %v", err)
	}

	inputFile := "example.txt"
	outputFile := "example.enc"

	// Create a sample input file
	if err := os.WriteFile(inputFile, []byte("Hello, File Encryption!"), 0644); err != nil {
		log.Fatalf("Failed to write to input file: %v", err)
	}
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	// Encrypt the file
	if err := dpapi.EncryptFile(inputFile, outputFile); err != nil {
		log.Fatalf("Failed to encrypt file: %v", err)
	}
	fmt.Printf("File encrypted: %s -> %s\n", inputFile, outputFile)

	// Decrypt the file
	if err := dpapi.DecryptFile(outputFile, inputFile); err != nil {
		log.Fatalf("Failed to decrypt file: %v", err)
	}
	fmt.Printf("File decrypted: %s -> %s\n", outputFile, inputFile)
}
```

## Logging

Verbose logging can be enabled by passing `true` to the `Start` function. This will provide detailed output about the operations being performed.

## License

This package is released under the MIT License. See [LICENSE](LICENSE) for details.

