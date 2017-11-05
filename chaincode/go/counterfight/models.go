package main

type Settings struct {
	Admin        string `json:"admin"`
}

type Carton struct {
	Id 				string `json:"id"`
	Name			string `json:"name"`
	ProductionDate 	int64 `json:"productionDate"`
	Description		string `json:"description"`
	PackageNum		int64 `json:"packageNum"`
	Producer		string `json:"producer"`
	Owner			string `json:"owner"`
}

type Package struct {
	Id 				string `json:"id"`
	Sold   			bool `json:"sold"`
	SellDate 		int64 `json:"sellDate"`
}

type User struct {
	Role    		string `json:"role"`
	Name        	string `json:"name"`
}

type CreateCartonResponse struct {
	Carton			Carton `json:"carton"`
	PackageList		[]Package `json:"packages"`
}
