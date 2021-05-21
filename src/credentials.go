package main

type databaseCredentials struct {
	DatabaseName string `json:"dbname"`
	Engine       string `json:"engine"`
	Host         string `json:"host"`
	Instance     string `json:"instance"`
	Password     string `json:"password"`
	Port         string `json:"port"`
	Username     string `json:"username"`
}
