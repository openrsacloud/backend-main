package db

type User struct {
	Id       string `json:"id"`
	UserName string `json:"username"`
	Admin    bool   `json:"admin"`
}

type Session struct {
	Id        string `json:"id"`
	End       string `json:"end"`
	IP        string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	User      string `json:"user"`
}
