---
frontend:
---

1) sign data:
// Import ethers.js library
const { ethers } = require('ethers');

// Connect to MetaMask provider
const provider = new ethers.providers.Web3Provider(window.ethereum);

// Function to sign data
async function signData(data) {
  try {
    // Get the user's Ethereum address from MetaMask
    const signer = provider.getSigner();
    const address = await signer.getAddress();

    // Sign the data
    const signature = await signer.signMessage(data);

    // Return the signed data and the user's Ethereum address
    return { signature, address };
  } catch (error) {
    console.error('Failed to sign the data:', error);
    return null;
  }
}

// Usage example
const dataToSign = 'Hello, World!';

signData(dataToSign)
  .then((signedData) => {
    if (signedData) {
      console.log('Signed Data:', signedData);
      // Send the signed data and user's Ethereum address to the backend
      // via an HTTP request
    }
  })
  .catch((error) => {
    console.error('Error while signing data:', error);
  });



2) verify it:
// Import ethers.js library
const { ethers } = require('ethers');

// Function to verify the signed data
async function verifySignedData(data, signature, address) {
  try {
    // Recover the signer's address from the signature
    const recoveredAddress = ethers.utils.verifyMessage(data, signature);

    // Compare the recovered address with the provided address
    if (recoveredAddress.toLowerCase() === address.toLowerCase()) {
      // Signature is valid and matches the provided address
      return true;
    } else {
      // Signature does not match the provided address
      return false;
    }
  } catch (error) {
    console.error('Failed to verify the signed data:', error);
    return false;
  }
}

// Usage example
const receivedData = {
  data: 'Hello, World!',
  signature: '<signature>',
  address: '<user-address>',
};

verifySignedData(receivedData.data, receivedData.signature, receivedData.address)
  .then((isValid) => {
    if (isValid) {
      console.log('Signature is valid!');
      // Process the verified data on the backend
    } else {
      console.log('Signature is invalid!');
      // Handle the invalid signature case
    }
  })
  .catch((error) => {
    console.error('Error while verifying signed data:', error);
  });




---
backend:
---

package main

import (
    "fmt"
    "log"
    "strings"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"
)

// Function to verify the signed data
func verifySignedData(data string, signature string, address string) bool {
    // Remove the '0x' prefix from the address if present
    address = strings.TrimPrefix(address, "0x")

    // Decode the signature and data
    sigBytes := common.FromHex(signature)
    dataHash := crypto.Keccak256Hash([]byte(data))

    // Recover the public key from the signature
    pubKey, err := crypto.SigToPub(dataHash.Bytes(), sigBytes)
    if err != nil {
        log.Println("Failed to recover public key from signature:", err)
        return false
    }

    // Generate the Ethereum address from the recovered public key
    recoveredAddress := crypto.PubkeyToAddress(*pubKey).Hex()

    // Compare the recovered address with the provided address
    return strings.ToLower(recoveredAddress) == strings.ToLower(address)
}

func main() {
    receivedData := struct {
        Data      string `json:"data"`
        Signature string `json:"signature"`
        Address   string `json:"address"`
    }{
        Data:      "Hello, World!",
        Signature: "<signature>",
        Address:   "<user-address>",
    }

    isValid := verifySignedData(receivedData.Data, receivedData.Signature, receivedData.Address)
    if isValid {
        fmt.Println("Signature is valid!")
        // Process the verified data on the backend
    } else {
        fmt.Println("Signature is invalid!")
        // Handle the invalid signature case
    }
}


---

import (
    "database/sql"
    "log"
    "strings"

    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, err := sql.Open("sqlite3", "path/to/database.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Define the values to be inserted
    values := []string{"Value1", "Value2", "Value3", "Value4"} // Example values, you can modify this as per your data

    // Prepare the insert statement with placeholders for the values
    stmt, err := db.Prepare("INSERT INTO your_table (column1) VALUES " + buildValuePlaceholders(len(values)))
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()

    // Build the slice of values as arguments for the Exec() method
    args := make([]interface{}, len(values))
    for i, value := range values {
        args[i] = value
    }

    // Execute the insert statement with the values
    _, err = stmt.Exec(args...)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Data inserted successfully")
}

// Helper function to build the value placeholders for the SQL statement
func buildValuePlaceholders(count int) string {
    placeholders := make([]string, count)
    for i := range placeholders {
        placeholders[i] = "(?)"
    }
    return strings.Join(placeholders, ", ")
}
