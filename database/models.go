package database

type User struct {
	Id       string `json:"id"`
	UserName string `json:"username"`
}

type Session struct {
	Id        string `json:"id"`
	End       string `json:"end"`
	IP        string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	User      string `json:"user"`
}
