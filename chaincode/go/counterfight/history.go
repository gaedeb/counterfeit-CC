package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"

)

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
