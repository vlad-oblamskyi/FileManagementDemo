package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

type FileManagement struct {
}

func (t *FileManagement) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

func (t *FileManagement) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "put":
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. File body is expected!");
		}
		body := args[0]
		messageDigest := getMessageDigest(body)

		stub.PutState(messageDigest, body)

		return []byte(messageDigest), nil
	default:
		return nil, errors.New("Unsupported operation")
	}
}

func (t *FileManagement) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "getFileByHash":
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. File hash is expected!");
		}
		return stub.GetState(args[0])
	case "checkIfFileExists":
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. File body is expected!");
		}
		body := args[0]
		messageDigest := getMessageDigest(body)
		state, _ := stub.GetState(messageDigest)

		isExist := strconv.FormatBool(state != nil)

		return []byte(isExist), nil
	default:
		return nil, errors.New("Unsupported operation")
	}
}

func getMessageDigest(body string) string {
	byteBody := []byte(body)
	hash := sha256.New()
	hash.Write(byteBody)
	messageDigest := hash.Sum(nil)
	return hex.EncodeToString(messageDigest)
}

func main() {
	err := shim.Start(new(FileManagement))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}