# File API

Endpoints for uploading, viewing and downloading files and folders.

**Group base URL:** `/api/files`

## Models
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

## Error
The response structure is the same for all errors. 
Example response:
```
	{
		"status":  401,
		"message": "Unauthorized",
		"info":    "Invalid token"
	}
```

## Endpoints
 - `POST /api/files/upload` 
 	- Authentication required: **Bearer \<TOKEN>**
	- Request form data:
	 ```
		------
		Content-Disposition: form-data; name="name" // OPTIONAL

		filename123.gif
		------
		Content-Disposition: form-data; name="parent" // OPTIONAL

		folders:ID
		------
		Content-Disposition: form-data; name="file"; filename="cat.gif"
		Content-Type: image/gif

		<contents of file here>
		--------
	 ```
	- Example response:
	 ```
		{
			"status":  200,
			"message": "OK"
		}
	 ```

 - `GET /api/files/folder/:id?`
	- Authentication required: **Bearer \<TOKEN>**
	- id parameter is optional (ex. folders:ID), if not ptovided, it returns the user's home directory
	- Example response:
	 ```
		{
			"status":  200,
			"message": "OK",
			"data": {
				?"parent":   Folder,
				"files":   []File,
				"folders": []Folder
			}
		}
	 ```

