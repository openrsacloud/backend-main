# File API

Endpoints for uploading, viewing and downloading files and folders.

**Group base URL:** `/api/files`

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
### `POST /api/files/upload` 
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

### `GET /api/files/get_folder/:id?`
- Authentication required: **Bearer \<TOKEN>**
- id parameter is optional, if not ptovided, it returns the user's home directory
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

### `GET /api/files/get_file/:id`
- Authentication is optional: **Bearer \<TOKEN>**
- id parameter required