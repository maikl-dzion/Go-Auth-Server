package model

type User struct {
	Id       int `json:"id, omitempty"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Accessable string `json:"accessable"`
	Tab_number string `json:"tab_number"`
}


type ChangePasswd struct {
	Login    string `json:"login"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}



type AccIbis struct {
	ClientLogin     string `json:"client_login"`
	ClientPassword  string `json:"client_password"`
	ServerLogin     string `json:"server_login"`
	ServerPassword  string `json:"server_password"`
	PresharedKey    string `json:"preshared_key"`
}


type IbisSession struct {
	Account int32
	Uuid    string
	Status  int
}


type AccSession struct {
	AccountId  int32
	Uuid       string
	ClosedAt   string
	OpenedAt   string
}


var BackendAccount int32 = 1