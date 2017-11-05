package main

import (
	"encoding/json"
	// "errors"
	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"./mock"
	"./testdata"
	"testing"
)

var settings = Settings{
	Admin:        "testUser",
}

func initToken(t *testing.T) *mock.FullMockStub {
	loyalty := &CounterfeitCC{}

	stub := mock.NewFullMockStub("counterfeit", loyalty)
	stub.MockCreator("default", testdata.TestUser1Cert)

	tokenBytes, _ := json.Marshal(settings)
	res := stub.MockInit("1", util.ToChaincodeArgs("init", string(tokenBytes)))
	if res.Status != shim.OK {
		t.Error("Counterfeit cc init failed: " + res.Message)
	}

	st := Settings{Admin: "testUser"}
	stBytes, _ := json.Marshal(st)
	infoRes := stub.MockInvoke("1", util.ToChaincodeArgs("info", string(stBytes)))
	settings := Settings{}
	err := json.Unmarshal(infoRes.Payload, &settings)

	if (err != nil) {
		t.Error("Could not get info")
	}

	if settings.Admin != st.Admin {
		t.Error("Chaincode admin name is wrong")
	}

	return stub
}

func TestInitToken(t *testing.T) {
	initToken(t)
}

