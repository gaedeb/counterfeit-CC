package main

type Settings struct {
	Admin        string `json:"admin"`
}

type Carton struct {
	Name			string `json:"name"`
	ProductionDate 	int64 `json:"productionDate"`
	Description		string `json:"description"`
	PackageNum		int64 `json:"packageNum"`
}

type Package struct {
	Sold   			bool `json:"sold"`
	SellDate 		int64 `json:"sellDate"`
}

type User struct {
	Role    		string `json:"role"`
	Name        	string `json:"name"`
}
