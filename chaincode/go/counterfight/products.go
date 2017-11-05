package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
	"encoding/json"
)


func (t *CounterfeitCC) createCarton(stub shim.ChaincodeStubInterface, id string, carton Carton) (*[]Package, error) {

	key, _ := stub.CreateCompositeKey(IndexCartons, []string{id})

	data, err := json.Marshal(carton)
	if err != nil {
		return nil, errors.New("Error marshaling carton object'" + id + "': " + err.Error())
	}

	err = stub.PutState(key, []byte(data))

	if err != nil {
		return nil, errors.New("Error creating user '" + id + "': " + err.Error())
	}

	var result []Package = []Package{}
	for i := 0; carton.PackageNum; i++ {

		pckg := Package{
			Id: uintToString(uint64Random()),
			Sold: false,
		}

		err := t.createPackage(stub, id, pckg.Id, pckg)
		if err != nil {
			return nil, err
		}

		result = append(result, pckg)
	}

	return &result, nil
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

func (t *CounterfeitCC) getCarton(stub shim.ChaincodeStubInterface, cartonId string) (Carton, error) {
	key, _ := stub.CreateCompositeKey(IndexCartons, []string{cartonId})
	data, err := stub.GetState(key)
	if err != nil {
		return nil, errors.New("Error getting carton: " + err.Error())
	}  else if data == nil {
		return nil, errors.New("No Carton for " + cartonId)
	}

	carton := Carton{}
	err = json.Unmarshal(data, &carton)
	if err != nil {
		return nil, errors.New("Error parsing carton json: " + err.Error())
	}

	return carton, nil
}

func (t *CounterfeitCC) updateCartonOwner(stub shim.ChaincodeStubInterface, cartonId string, newOwner string) error {
	carton, err := t.getCarton(stub, cartonId)
	if err != nil {
		return err
	}

	key, _ := stub.CreateCompositeKey(IndexCartons, []string{cartonId})

	carton.Owner = newOwner

	data, err := json.Marshal(carton)
	if err != nil {
		return errors.New("Error marshaling carton object: " + err.Error())
	}
	err = stub.PutState(key, data)

	if err != nil {
		return errors.New("Error storing carton: " + err.Error())
	}

	return nil
}