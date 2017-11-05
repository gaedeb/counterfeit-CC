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

type SellCarton struct {
	CartonId    	string `json:"cartonId"`
	Buyer        	string `json:"buyer"`
}

type SellPackage struct {
	CartonId    	string `json:"cartonId"`
	PackageId    	string `json:"packageId"`
}
