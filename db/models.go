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

type File struct {
	Id       string                 `json:"id"`
	Parent   string                 `json:"parent"`
	Name     string                 `json:"name"`
	Owner    string                 `json:"owner"`
	Size     uint64                 `json:"size"`
	Created  string                 `json:"created"`
	Modified string                 `json:"modified"`
	Type     string                 `json:"type"`
	Metadata map[string]interface{} `json:"metadata"`
}

type Folder struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Parent   string `json:"parent"`
	Owner    string `json:"owner"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
}
