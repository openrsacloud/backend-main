# Admin API

Endpoints used for administration and server management.

**Group base URL:** `/api/admin`

## Error
The response structure is the same for all errors. 
Example response:
```
{
	"status":  401,
	"message": "Unauthorized",
	"info":    "The username or password is incorrect"
}
```

## Endpoints
### `POST /api/admin/create_user` 
- Authentication required: **Bearer \<TOKEN>**
- Request body:
	```
	{
		"username": string,
		"password": string,
		"admin":    boolean
	}
	```
- OK response:
	```
	{
		"status":  200,
		"message": "OK",
		"data":    User
	}
	```

### `GET /api/admin/get_users` 
- Authentication required: **Bearer \<TOKEN>**
- Gets 50 user records, use the `page` query property to get more
- Request body:
	```
	{
		"username": string,
		"password": string,
		"admin":    boolean
	}
	```
- OK response:
	```
	{
		"status":  200,
		"message": "OK",
		"data":    []User
	}
	```