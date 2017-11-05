package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
	"encoding/json"
)


func (t *CounterfeitCC) createCarton(stub shim.ChaincodeStubInterface, id string, carton Carton) error {

	key, _ := stub.CreateCompositeKey(IndexCartons, []string{id})

	data, err := json.Marshal(carton)
	if err != nil {
		return errors.New("Error marshaling carton object'" + id + "': " + err.Error())
	}

	err = stub.PutState(key, []byte(data))

	if err != nil {
		return errors.New("Error creating user '" + id + "': " + err.Error())
	}

	return nil
}