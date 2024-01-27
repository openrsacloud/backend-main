# Auth API

Endpoints used for user authentication and session management.

**Group base URL:** `/api/auth`

## Models
 - User
	```
		{
			id       string, // "users:ID"
			username string
		}
	```

 - Session
	```
		{
			id         string, // "sessions:ID"
			end        string, // ISO_8601 in UTC
			ip_address string,
			user_agent string,
			user       string  // "user:ID"
		}
	```

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
 - `POST /api/auth/login` 
	- Request body:
	 ```
	 	{
			"username": string,
			"password": string
		}
	 ```
	- Example response:
	 ```
		{
			"status":  200,
			"message": "OK",
			"data": {
				"user":  User,
				"token": <HS256_TOKEN>,
			}
		}
	 ```

 - `POST /api/auth/clear_sessions`
	- Authentication required: **Bearer \<TOKEN>**
	- Example response:
	 ```
		{
			"status":  200,
			"message": "OK"
		}
	 ```

 - `GET /api/auth/get_sessions`
	- Authentication required: **Bearer \<TOKEN>**
	- Example response:
	 ```
		{
			"status":  200,
			"message": "OK",
			"data":    []Session
		}
	 ```

 - `POST /api/auth/remove_session`
	- Authentication required: **Bearer \<TOKEN>**
	- Request body:
	 ```
	 	{
			"session_id": string
		}
	 ```
	- Example response:
	 ```
		{
			"status":  200,
			"message": "OK"
		}
	 ```

 - `GET /api/auth/get_user`
	- Authentication required: **Bearer \<TOKEN>**
	- Example response:
	 ```
		{
			"status":  200,
			"message": "OK",
			"data":    []User
		}
	 ```