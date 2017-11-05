package main

import (
	"encoding/json"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type CounterfeitCC struct {
}

const KeySettings = "__settings"
const IndexProducer = "cn~producer"
const IndexPharmacy = "cn~pharmacy"
const IndexReseller = "cn~reseller"
const IndexCartons = "cn~carton"
const IndexPackage = "cn~package"

func (t *CounterfeitCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()


	if function != "init" {
		return shim.Error("Expected 'init' function.")
	}

	if len(args) != 1 {
		return shim.Error("Expected 1 argument, but got " + strconv.Itoa(len(args)))
	}

	// get token data from JSON
	settings := Settings{}
	err := json.Unmarshal([]byte(args[0]), &settings)

	if err != nil {
		return shim.Error("Error parsing settings json")
	}

	err = stub.PutState(KeySettings, []byte(args[0]))
	if err != nil {
		return shim.Error("Error saving token data")
	}

	return shim.Success(nil)
}

func (t *CounterfeitCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	// call routing
	switch function {
	case "info":
		info, _ := stub.GetState(KeySettings)
		return shim.Success(info)
	case "createUser":
		return t.registerUser(stub, args)
	case "createCarton":
		return t.registerCarton(stub, args)
	case "sellCarton":
		return t.registerCarton(stub, args)
	default:
		return shim.Error("Incorrect function name: " + function)
	}
}

func (t *CounterfeitCC) getSettings(stub shim.ChaincodeStubInterface) (Settings, error) {
	settingsByteArr, err := stub.GetState(KeySettings)
	if err != nil {
		return Settings{}, err
	}

	settings := Settings{}
	err = json.Unmarshal(settingsByteArr, &settings)
	if err != nil {
		return Settings{}, err
	}

	return settings, nil
}

func (t *CounterfeitCC) registerUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("expected 1 argument")
	}

	caller, err := CallerCN(stub)
	if err != nil {
		return shim.Error("Error extracting user identity")
	}

	err = t.createUser(stub, caller, args[0])

	if err != nil {
		return shim.Error("Error creating user '" + caller + "'")
	}

	return shim.Success(nil)
}

func (t *CounterfeitCC) registerCarton(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("expected 1 argument")
	}

	caller, err := CallerCN(stub)
	if err != nil {
		return shim.Error("Error extracting user identity")
	}

	carton := Carton{}
	err = json.Unmarshal([]byte(args[0]), &carton)
	if err != nil {
		return shim.Error("Error parsing carton json")
	}

	carton.Producer = caller
	carton.Id = uintToString(uint64Random())

	packages, err := t.createCarton(stub, carton.Id, carton)

	response := CreateCartonResponse{
		Carton: carton,
		PackageList: *packages,
	}

	data, _ := json.Marshal(response)
	if err != nil {
		return shim.Error("Error generating response")
	}

	return shim.Success(data)
}


func (t *CounterfeitCC) sellCarton(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("expected 1 argument")
	}

	caller, err := CallerCN(stub)
	if err != nil {
		return shim.Error("Error extracting user identity")
	}

	sellCarton := SellCarton{}
	err = json.Unmarshal([]byte(args[0]), &sellCarton)
	if err != nil {
		return shim.Error("Error parsing sellCarton request json")
	}

	carton, err := t.getCarton(stub, sellCarton.CartonId)
	if err != nil {
		return shim.Error(err.Error())
	}

	if carton.Owner != caller {
		return shim.Error("Carton doesn't belong to you!")
	}

	err = t.updateCartonOwner(stub, sellCarton.CartonId, sellCarton.Buyer)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}