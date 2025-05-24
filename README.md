# chirpy
boot.dev course content: Learn HTTP Servers in Go

#Administrative EndPoints: METHOD ENDPOINT APIFUNCTION
##GET /api/healthz api.Health\n
    Returns 200 when server is running.\n
##GET /admin/metricscfg.HitTotal\n
    Returns 200 and number of hits on /app path\n
##POST /admin/reset cfg.Reset\n
    Returns 200 and resets hit counter on /app\n
    **Resets users table in database (requirement for course, this would be not available in an actual system)**\n

#Application EndPoints\n
##POST /api/login api.UserLogin\n
    ```Expects body:
    {
		"password": "text",
		"email": "valid@email.com"
	}```

    Authenticates user password and issues an access token\n
    
    Returns 200 and user struct with access token\n
    ```type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	IsChirpyRed    bool      `json:"is_chirpy_red"`
	Token          string    `json:"token"`
	RefreshToken   string    `json:"refresh_token"`
}```

##POST /api/refresh api.UpdateAccessToken\n
    Expects "Authorization: Bearer" header with valid refresh token\n
    
    Returns 200 and a new access token\n

##POST /api/revoke api.RevokeRefreshToken\n
    Expects "Authorization: Bearer" header with valid refresh token\n
    
    Revokes refresh token\n
    
    Returns 204 and no body\n

##GET /api/chirps/{chirpID} api.GetChirp\n
    Expects /api/chirps/{chirpID} where {chirpID} is the UUID for a chirp\n
    
    Looks up the specific Chirp and returns 404 if not found\n
    
    Returns 200 and chirp struct\n
    ```Chirp struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
}```

##DELETE /api/chirps/{id} api.DeleteChirp\n
    Expects /api/chirps/{chirpID} where {chirpID} is the UUID for a chirp and a valid access token in "Authorization: Bearer" header\n
    
    Validates token and that user is author of chirp\n
    Chirp is deleted from the database\n
    
    Returns 204 and blank body on success\n

##POST /api/chirps api.NewChirp\n
    ```Expects valid access token in "Authorization: Bearer" header\n
    Expects body:
    {
        "body": "text string for chirp"
    }```
    
    Creates a new chirp and assigns a UUID.\n
    Associates chirp with user.\n
    Applies profanity filtering to chirp body.\n
    
    Returns 201 and Chirp struct\n
    ```type Chirp struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
}```

##GET /api/chirps api.GetAllChirps\n
    Accepts an optional author_id parameter. Parameter is the UUID of a valid user.\n
    Accepts an optional sort parameter. Valid values "asc" or "desc". Defaults to "asc".\n

    Chirps are returned in ascending order of created_by field.\n
    If author_id is passed returns only chirps associated with that user, otherwise returns all chirps in database.\n
    If sort="desc" is passed, chirps will sort in descending order of created_by field.\n

    Returns 200 and an array of the Chirps struct\n
    ```type Chirp struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
}```

##POST /api/users api.NewUser\n
    ```Expects body:
    {
		"password": "text",
		"email": "valid@email.com"
	}```

    Creates a new user with provided password.\n
    *User must login to get access token*\n
    
    Returns 200 and user struct\n
    ```type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	IsChirpyRed    bool      `json:"is_chirpy_red"`
}```

##PUT /api/users api.UpdateUser\n
    ```Expects valid access token in "Authorization: Bearer" header\n
    Expects body:
    {
		"password": "text",
		"email": "valid@email.com"
	}```

    Updates user password and email.\n
    
    Returns 200 and user struct\n
    ```type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	IsChirpyRed    bool      `json:"is_chirpy_red"`
}```

##POST /api/polka/webhooks api.UpdateChirpyRed\n
    3rd party payment API webhook.\n
    ```Expects valid ApiKey header.
    Expects body:
    {
  "event": "user.upgraded",
  "data": {
    "user_id": "valid UUID"
  }
}```

    If a valid UUID is provided and the user.upgraded event is present, user is updated to have Chirpy Red in database.\n

    Returns 204 and no body once account is upgraded or if event does not match user.upgraded.\n
    Returns 404 if user UUID is invalid.\n


