package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
	"encoding/json"
)


func (t *CounterfeitCC) userExists(stub shim.ChaincodeStubInterface, cn string, role string) bool {
	var prefix string
	switch role {
	case "producer":
		prefix = IndexProducer
	case "pharmacy":
		prefix = IndexPharmacy
	case "reseller":
		prefix = IndexReseller
	default:
		return false
	}

	key, _ := stub.CreateCompositeKey(prefix, []string{cn})

	data, err := stub.GetState(key)
	if err != nil {
		return false
	} else if data == nil {
		return false
	}

	return true
}

func (t *CounterfeitCC) createUser(stub shim.ChaincodeStubInterface, cn string, role string) error {
	var prefix string
	switch role {
	case "producer":
		prefix = IndexProducer
	case "pharmacy":
		prefix = IndexPharmacy
	case "reseller":
		prefix = IndexReseller
	default:
		return errors.New("Unknown user roll")
	}

	key, _ := stub.CreateCompositeKey(prefix, []string{cn})

	user := User{
		Name: cn,
		Role: role,
	}

	data, err := json.Marshal(user)
	if err != nil {
		return errors.New("Error marshaling user object'" + cn + "' with the role '" + role + "': " + err.Error())
	}

	err = stub.PutState(key, []byte(data))

	if err != nil {
		return errors.New("Error creating user '" + cn + "' with the role '" + role + "': " + err.Error())
	}

	return nil
}