package main

import (
	"encoding/json"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"time"
	"fmt"
	"errors"
	"crypto/x509"
	"encoding/pem"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/protos/msp"
	"math/rand"
)

type CounterfeitCC struct {
}

type Settings struct {
	Admin        string `json:"admin"`
}

type Carton struct {
	Id 				string `json:"id"`
	Name			string `json:"name"`
	ProductionDate 	time.Time `json:"productionDate"`
	Description		string `json:"description"`
	PackageNum		int `json:"packageNum"`
	Producer		string `json:"producer"`
	Owner			string `json:"owner"`
}

type Package struct {
	Id 				string `json:"id"`
	Sold   			bool `json:"sold"`
	SellDate 		time.Time `json:"sellDate"`
}

type User struct {
	Role    		string `json:"role"`
	Name        	string `json:"name"`
}

type CreateCartonResponse struct {
	Carton			Carton `json:"carton"`
	PackageList		[]Package `json:"packages"`
}

type CartonRef struct {
	CartonId    	string `json:"cartonId"`
	Buyer        	string `json:"buyer"`
}

type PackageRef struct {
	CartonId    	string `json:"cartonId"`
	PackageId    	string `json:"packageId"`
}

type HistoryEntry struct {
	Owner 		string `json:"owner"`
	TxId 		string `json:"txId"`
	Timestamp 	int64 `json:"timeStamp"`
}

type PackageHistoryResponse struct {
	Carton Carton `json:"carton"`
	Package Package `json:"package"`
	OwnerHistory []HistoryEntry `json:"ownerHistory"`
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
		return t.sellCarton(stub, args)
	case "sellPackage":
		return t.sellPackage(stub, args)
	case "getPackageHistory":
		return t.getPackageHistory(stub, args)
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
	carton.ProductionDate = time.Now()

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

	sellCarton := CartonRef{}
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

func (t *CounterfeitCC) sellPackage(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("expected 1 argument")
	}

	caller, err := CallerCN(stub)
	if err != nil {
		return shim.Error("Error extracting user identity")
	}

	sellPackage := PackageRef{}
	err = json.Unmarshal([]byte(args[0]), &sellPackage)
	if err != nil {
		return shim.Error("Error parsing sellPackage request json")
	}

	carton, err := t.getCarton(stub, sellPackage.CartonId)
	if err != nil {
		return shim.Error(err.Error())
	}

	if carton.Owner != caller {
		return shim.Error("Carton doesn't belong to you!")
	}

	err = t.markPackageSold(stub, sellPackage.CartonId, sellPackage.PackageId)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *CounterfeitCC) getPackageHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("expected 1 argument")
	}

	packageRef := PackageRef{}
	err := json.Unmarshal([]byte(args[0]), &packageRef)
	if err != nil {
		return shim.Error("Error parsing sellPackage request json")
	}

	carton, err := t.getCarton(stub, packageRef.CartonId)
	if err != nil {
		return shim.Error(err.Error())
	}

	pckg, err := t.getPackage(stub, packageRef.CartonId, packageRef.PackageId)
	if err != nil {
		return shim.Error(err.Error())
	}

	key, _ := stub.CreateCompositeKey(IndexCartons, []string{packageRef.CartonId})
	history, err := t.getHistory(stub, key)
	if err != nil {
		return shim.Error(err.Error())
	}

	response := PackageHistoryResponse{
		Carton: carton,
		Package: pckg,
		OwnerHistory: history,
	}

	data, _ := json.Marshal(response)
	if err != nil {
		return shim.Error("Error generating package history response")
	}

	return shim.Success(data)
}

// ------------------------------------------------------------------
func (t *CounterfeitCC) getHistory(stub shim.ChaincodeStubInterface, key string) ([]HistoryEntry, error) {

	historyIer, err := stub.GetHistoryForKey(key)

	if err != nil {
		return nil, err
	}

	var history []HistoryEntry = []HistoryEntry{}
	for i := 0; historyIer.HasNext(); i++ {
		modification, err := historyIer.Next()
		if err != nil {
			return nil, err
		}

		historyEntry := &HistoryEntry{
			TxId: modification.TxId,
			Timestamp: modification.Timestamp.Seconds,
			Owner:   string(modification.Value),
		}

		history = append(history, *historyEntry)
	}

	return history, nil
}
// ------------------------------------------------------------------
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
	for i := 0; i < carton.PackageNum; i++ {

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
		return Carton{}, errors.New("Error getting carton: " + err.Error())
	}  else if data == nil {
		return Carton{}, errors.New("No Carton for " + cartonId)
	}

	carton := Carton{}
	err = json.Unmarshal(data, &carton)
	if err != nil {
		return Carton{}, errors.New("Error parsing carton json: " + err.Error())
	}

	return carton, nil
}

func (t *CounterfeitCC) getPackage(stub shim.ChaincodeStubInterface, cartonId string, packageId string) (Package, error) {
	key, _ := stub.CreateCompositeKey(IndexPackage, []string{cartonId, packageId})
	data, err := stub.GetState(key)
	if err != nil {
		return Package{}, errors.New("Error getting package: " + err.Error())
	}  else if data == nil {
		return Package{}, errors.New("No package for " + cartonId + ":" + packageId)
	}

	pckg := Package{}
	err = json.Unmarshal(data, &pckg)
	if err != nil {
		return Package{}, errors.New("Error parsing package json: " + err.Error())
	}

	return pckg, nil
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

func (t *CounterfeitCC) markPackageSold(stub shim.ChaincodeStubInterface, cartonId string, packageId string) error {
	pckg, err := t.getPackage(stub, cartonId, packageId)
	if err != nil {
		return err
	}

	pckg.Sold = true
	pckg.SellDate = time.Now()

	key, _ := stub.CreateCompositeKey(IndexPackage, []string{cartonId, packageId})

	data, err := json.Marshal(pckg)
	if err != nil {
		return errors.New("Error marshaling package object: " + err.Error())
	}
	err = stub.PutState(key, data)

	if err != nil {
		return errors.New("Error storing package: " + err.Error())
	}

	return nil
}
// ------------------------------------------------------------------
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

// ------------------------------------------------------------------
func parsePEM(certPEM string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil, errors.New("Failed to parse PEM certificate")
	}

	return x509.ParseCertificate(block.Bytes)
}

// extracts CN from an x509 certificate
func CNFromX509(certPEM string) (string, error) {
	cert, err := parsePEM(certPEM)
	if err != nil {
		return "", errors.New("Failed to parse certificate: " + err.Error())
	}
	return cert.Subject.CommonName, nil
}

// extracts CN from caller of a chaincode function
func CallerCN(stub shim.ChaincodeStubInterface) (string, error) {
	data, _ := stub.GetCreator()
	serializedId := msp.SerializedIdentity{}
	err := proto.Unmarshal(data, &serializedId)
	if err != nil {
		return "", errors.New("Could not unmarshal Creator")
	}

	cn, err := CNFromX509(string(serializedId.IdBytes))
	if err != nil {
		return "", err
	}
	return cn, nil
}

func uintToString(num uint64) (string) {
	return strconv.FormatUint(num, 10)
}

func uint64Random() uint64 {
	return uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
}

// ------------------------------------------------------------------
func main() {
	err := shim.Start(new(CounterfeitCC))
	if err != nil {
		fmt.Printf("Error starting CounterfeitCC: %s", err)
	}
}