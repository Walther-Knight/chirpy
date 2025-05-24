# chirpy
boot.dev course content: Learn HTTP Servers in Go

#Administrative EndPoints: METHOD ENDPOINT APIFUNCTION
##GET /api/healthz api.Health  
    Returns 200 when server is running.  
##GET /admin/metricscfg.HitTotal  
    Returns 200 and number of hits on /app path  
##POST /admin/reset cfg.Reset  
    Returns 200 and resets hit counter on /app  
    **Resets users table in database (requirement for course, this would be not available in an actual system)**  

#Application EndPoints  
##POST /api/login api.UserLogin  
    ```Expects body:
    {
		"password": "text",
		"email": "valid@email.com"
	}```

    Authenticates user password and issues an access token  
    
    Returns 200 and user struct with access token  
    ```type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	IsChirpyRed    bool      `json:"is_chirpy_red"`
	Token          string    `json:"token"`
	RefreshToken   string    `json:"refresh_token"`
}```

##POST /api/refresh api.UpdateAccessToken  
    Expects "Authorization: Bearer" header with valid refresh token  
    
    Returns 200 and a new access token  

##POST /api/revoke api.RevokeRefreshToken  
    Expects "Authorization: Bearer" header with valid refresh token  
    
    Revokes refresh token  
    
    Returns 204 and no body  

##GET /api/chirps/{chirpID} api.GetChirp  
    Expects /api/chirps/{chirpID} where {chirpID} is the UUID for a chirp  
    
    Looks up the specific Chirp and returns 404 if not found  
    
    Returns 200 and chirp struct  
    ```Chirp struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
}```

##DELETE /api/chirps/{id} api.DeleteChirp  
    Expects /api/chirps/{chirpID} where {chirpID} is the UUID for a chirp and a valid access token in "Authorization: Bearer" header  
    
    Validates token and that user is author of chirp  
    Chirp is deleted from the database  
    
    Returns 204 and blank body on success  

##POST /api/chirps api.NewChirp  
    ```Expects valid access token in "Authorization: Bearer" header  
    Expects body:
    {
        "body": "text string for chirp"
    }```
    
    Creates a new chirp and assigns a UUID.  
    Associates chirp with user.  
    Applies profanity filtering to chirp body.  
    
    Returns 201 and Chirp struct  
    ```type Chirp struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
}```

##GET /api/chirps api.GetAllChirps  
    Accepts an optional author_id parameter. Parameter is the UUID of a valid user.  
    Accepts an optional sort parameter. Valid values "asc" or "desc". Defaults to "asc".  

    Chirps are returned in ascending order of created_by field.  
    If author_id is passed returns only chirps associated with that user, otherwise returns all chirps in database.  
    If sort="desc" is passed, chirps will sort in descending order of created_by field.  

    Returns 200 and an array of the Chirps struct  
    ```type Chirp struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
}```

##POST /api/users api.NewUser  
    ```Expects body:
    {
		"password": "text",
		"email": "valid@email.com"
	}```

    Creates a new user with provided password.  
    *User must login to get access token*  
    
    Returns 200 and user struct  
    ```type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	IsChirpyRed    bool      `json:"is_chirpy_red"`
}```

##PUT /api/users api.UpdateUser  
    ```Expects valid access token in "Authorization: Bearer" header  
    Expects body:
    {
		"password": "text",
		"email": "valid@email.com"
	}```

    Updates user password and email.  
    
    Returns 200 and user struct  
    ```type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	IsChirpyRed    bool      `json:"is_chirpy_red"`
}```

##POST /api/polka/webhooks api.UpdateChirpyRed  
    3rd party payment API webhook.  
    ```Expects valid ApiKey header.
    Expects body:
    {
  "event": "user.upgraded",
  "data": {
    "user_id": "valid UUID"
  }
}```

    If a valid UUID is provided and the user.upgraded event is present, user is updated to have Chirpy Red in database.  

    Returns 204 and no body once account is upgraded or if event does not match user.upgraded.  
    Returns 404 if user UUID is invalid.  


