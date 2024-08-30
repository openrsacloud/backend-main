# Models

 - User
	```
		{
			id       string,
			username string,
			admin    boolean
		}
	```

 - Session
	```
		{
			id         string, // "sessions:ID"
			end        string, // ISO_8601 in UTC
			ip_address string,
			user_agent string,
			user       string
		}
	```

 - Share
	```
		{
			"id":         string
			"owner":      string    // user ID
			"recipients": []string  // an array of user IDs or an empty array if public
			"object":     string    // a file or folder ID
		}
	```

 - File
	```
		{
			id       string, // "files:ID"
			name     string,
			parent   string, // "folders:ID"
			owner    string, // "user:ID"
			created  string, // ISO_8601 in UTC
			modified string, // ISO_8601 in UTC
			type     string, // MIME type (ex. image/gif)
			size     int64,  // in megabytes
			metadata object
		}
	```

 - Folder
	```
		{
			id       string, // "folders:ID"
			name     string,
			parent   string, // "folders:ID"
			owner    string, // "user:ID"
			created  string, // ISO_8601 in UTC
			modified string, // ISO_8601 in UTC
		}
	```