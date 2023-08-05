func main() {
    receivedData := struct {
        Data      map[string]interface{} `json:"data"`
        Signature string                 `json:"signature"`
        Address   string                 `json:"address"`
    }{
        Data:      map[string]interface{}{"name": "John Doe", "email": "john@example.com", "age": 30, "nonce": time.Now().Unix()},
        Signature: "<signature>",
        Address:   "<user-address>",
    }

    isValid := verifySignedDataWithNonce(receivedData.Data, receivedData.Signature, receivedData.Address, receivedData.Data["nonce"].(int64))
    if isValid {
        fmt.Println("Signature is valid!")
        // Process the verified data on the backend
    } else {
        fmt.Println("Signature is invalid!")
        // Handle the invalid signature case
    }
}
