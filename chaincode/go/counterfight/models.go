package main

import "time"

type Settings struct {
	Admin        string `json:"admin"`
}

type Carton struct {
	Id 				string `json:"id"`
	Name			string `json:"name"`
	ProductionDate 	time.Time `json:"productionDate"`
	Description		string `json:"description"`
	PackageNum		int64 `json:"packageNum"`
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
