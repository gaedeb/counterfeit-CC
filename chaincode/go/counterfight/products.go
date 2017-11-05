package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
	"encoding/json"
)


func (t *CounterfeitCC) createCarton(stub shim.ChaincodeStubInterface, id string, carton Carton) ([]*Package, error) {

	key, _ := stub.CreateCompositeKey(IndexCartons, []string{id})

	data, err := json.Marshal(carton)
	if err != nil {
		return nil, errors.New("Error marshaling carton object'" + id + "': " + err.Error())
	}

	err = stub.PutState(key, []byte(data))

	if err != nil {
		return nil, errors.New("Error creating user '" + id + "': " + err.Error())
	}

	var result []*Package = []*Package{}
	for i := 0; carton.PackageNum; i++ {

		pckg := Package{
			Sold: false,
		}

		packageId := uintToString(uint64Random())

		err := t.createPackage(stub, id, packageId, pckg)
		if err != nil {
			return nil, err
		}

		result = append(result, &pckg)
	}

	return result, nil
}

func (t *CounterfeitCC) createPackage(stub shim.ChaincodeStubInterface, cartonId string, id string, pckg Package) error {

	key, _ := stub.CreateCompositeKey(IndexPackage, []string{cartonId, id})

	data, err := json.Marshal(pckg)
	if err != nil {
		return errors.New("Error marshaling package object'" + id + "': " + err.Error())
	}

	err = stub.PutState(key, []byte(data))

	if err != nil {
		return errors.New("Error creating package '" + id + "': " + err.Error())
	}

	return nil
}