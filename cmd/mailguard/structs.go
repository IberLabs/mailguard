package main

type Config struct {
	Smtp 	struct {
		Host     	string `json:"host"`
		Port	 	string `json:"port"`
	} `json:"smtp"`
	Imap 	struct {
		Host     	string `json:"host"`
		Port 		string `json:"port"`
	} `json:"imap"`
	Auth 	struct {
		Username	string `json:"username"`
		Password 	string `json:"password"`
	} `json:"auth"`
	Account struct{
		Sender		string `json:"sender"`
	}
}
