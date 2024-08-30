# File API

Endpoints for uploading, viewing and downloading files and folders.

**Group base URL:** `/api/share`

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
### `GET /api/share/:id` 
- Authentication is optional: **Bearer \<TOKEN>**
- Example response:
	```
	{
		"status":  200,
		"message": "OK",
		"data":    File | Folder + Items
	}
	```

### `POST /api/share/create_share` 
- Authentication required: **Bearer \<TOKEN>**
- Request body:
	```
	{
		"object": string,     // File or Folder ID
		"recipients": string  // an array of User IDs or empty
	}
	```
- Example response:
	```
	{
		"status":  200,
		"message": "OK",
		"data":    Share
	}
	```

### `POST /api/share/remove_share`
- Authentication required: **Bearer \<TOKEN>**
- Request body:
	```
	{
		"id": string
	}
	```
- Example response:
	```
	{
		"status":  200,
		"message": "OK",
		"data":    Share
	}
	```

### `POST /api/share/update_share`
- Authentication required: **Bearer \<TOKEN>**
- Request body:
	```
	{
		"id": string,
		"recipients": []string  // an array of User IDs or empty
	}
	```
- Example response:
	```
	{
		"status":  200,
		"message": "OK",
		"data":    Share
	}
	```

### `GET /api/share/shared_with_me` 
- Authentication required: **Bearer \<TOKEN>**
- Example response:
	```
	{
		"status":  200,
		"message": "OK",
		"data":    []Share
	}
	```