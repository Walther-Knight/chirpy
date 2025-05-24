package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Walther-Knight/chirpy/internal/auth"
	"github.com/Walther-Knight/chirpy/internal/database"
	"github.com/Walther-Knight/chirpy/internal/middleware"
	"github.com/Walther-Knight/chirpy/internal/models"
	"github.com/google/uuid"
)

func decodeJSONBody(r *http.Request, v any) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
}

func marshalJSON(w http.ResponseWriter, v any) error {
	res, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Write(res)
	return nil
}

func writeErrorResponse(w http.ResponseWriter, status int, message string) {
	errorJSON := models.ErrorBody{
		Error: message,
	}

	w.WriteHeader(status)
	err := marshalJSON(w, errorJSON)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.Write([]byte(`{"error: Something went wrong converting JSON"}`))
		return
	}
}

func writeSuccessResponse(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	err := marshalJSON(w, data)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.Write([]byte(`{"error: Something went wrong converting JSON"}`))
		return
	}
}

// removed bool lesson 5:6 and filtered bool var and filtered return
func profanityFilter(s string) string {
	splitString := strings.Split(s, " ")
	var cleanString []string
	for _, word := range splitString {
		switch {
		case strings.ToLower(word) == "kerfuffle":
			cleanString = append(cleanString, "****")
		case strings.ToLower(word) == "sharbert":
			cleanString = append(cleanString, "****")
		case strings.ToLower(word) == "fornax":
			cleanString = append(cleanString, "****")
		default:
			cleanString = append(cleanString, word)
		}
	}
	return strings.Join(cleanString, " ")
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func NewChirp(api *middleware.ApiConfig, w http.ResponseWriter, r *http.Request) {
	type validateBody struct {
		Body string `json:"body"`
	}

	w.Header().Set("Content-Type", "application/json")

	userToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, "unauthorized user, no token")
		return
	}

	UserId, err := auth.ValidateJWT(userToken, api.Token)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, "unauthorized user, invalid token")
		return
	}

	params := validateBody{}
	errDecode := decodeJSONBody(r, &params)

	if errDecode != nil {
		log.Printf("Error decoding parameters: %s", errDecode)
		writeErrorResponse(w, http.StatusBadRequest, "error decoding JSON")
		return
	}

	if len(params.Body) > 140 {
		writeErrorResponse(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	if params.Body == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Chirp must contain characters")
		return
	}

	res, err := api.Db.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      profanityFilter(params.Body),
		UserID:    UserId,
	})
	if err != nil {
		log.Printf("Error on database: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "database error reported")
		return
	}

	ResJson := models.Chirp{
		ID:        res.ID.String(),
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		Body:      res.Body,
		UserID:    res.UserID.String(),
	}

	writeSuccessResponse(w, http.StatusCreated, ResJson)
}

func GetAllChirps(api *middleware.ApiConfig, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	res, err := api.Db.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Error on database: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "database error reported")
		return
	}

	ResJson := []models.Chirp{}
	for _, chirp := range res {
		chirpModel := models.Chirp{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
		}
		ResJson = append(ResJson, chirpModel)
	}

	writeSuccessResponse(w, http.StatusOK, ResJson)
}

func GetChirp(api *middleware.ApiConfig, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	chirpID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		log.Printf("Error on database: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "database error reported")
		return
	}

	res, err := api.Db.GetChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Error on database: %v", err)
			writeErrorResponse(w, http.StatusNotFound, "error: Chirp ID does not exist")
			return
		}
		log.Printf("Error on database: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "database error reported")
		return
	}

	ResJson := models.Chirp{
		ID:        res.ID.String(),
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		Body:      res.Body,
		UserID:    res.UserID.String(),
	}

	writeSuccessResponse(w, http.StatusOK, ResJson)
}

func validateEmail(s string) bool {
	if !(strings.Index(s, "@") < strings.LastIndex(s, ".")) {
		return false
	}
	if !(strings.Count(s, "@") == 1) {
		return false
	}
	if !(strings.LastIndex(s, ".") < len(s)-2) {
		return false
	}
	return true
}

func NewUser(api *middleware.ApiConfig, w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	w.Header().Set("Content-Type", "application/json")
	params := reqParams{}
	errDecode := decodeJSONBody(r, &params)

	if errDecode != nil {
		log.Printf("Error decoding parameters: %s", errDecode)
		writeErrorResponse(w, http.StatusBadRequest, "error decoding JSON")
		return
	}
	//case for minimum possible email length a@a.aa
	if len(params.Email) < 5 {
		writeErrorResponse(w, http.StatusBadRequest, "invalid email submitted")
		return
	}
	//basic email formatting checks
	if !validateEmail(params.Email) {
		writeErrorResponse(w, http.StatusBadRequest, "invalid email submitted")
		return
	}

	pwd, errHash := auth.HashPassword(params.Password)
	if errHash != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "error hashing password")
		return
	}

	res, err := api.Db.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          params.Email,
		HashedPassword: pwd,
	})
	if err != nil {
		log.Printf("Error on database: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "database error reported")
		return
	}

	ResJson := models.User{
		ID:          res.ID,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
		Email:       res.Email,
		IsChirpyRed: res.IsChirpyRed.Bool,
	}

	log.Printf("User: %s created with ID %v", res.Email, res.ID)
	writeSuccessResponse(w, http.StatusCreated, ResJson)

}

func UserLogin(api *middleware.ApiConfig, w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	w.Header().Set("Content-Type", "application/json")
	params := reqParams{}
	errDecode := decodeJSONBody(r, &params)

	if errDecode != nil {
		log.Printf("Error decoding parameters: %s", errDecode)
		writeErrorResponse(w, http.StatusBadRequest, "error decoding JSON")
		return
	}

	userInfo, err := api.Db.GetUserPassword(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error on database: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "database error reported")
	}
	err = auth.CheckPasswordHash(userInfo.HashedPassword, params.Password)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, "incorrect email or password")
		return
	}

	ExpiresIn := 1 * time.Hour
	newToken, err := auth.MakeJWT(userInfo.ID, api.Token, ExpiresIn)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "jwt token creation failed")
		return
	}

	newRefreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "refresh token creation failed")
		return
	}
	//update database with token. No err handling on ParseDuration because I'm confident that won't error
	duration, err2 := time.ParseDuration("1440h")
	log.Println(err2)
	err = api.Db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     newRefreshToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userInfo.ID,
		ExpiresAt: time.Now().Add(duration),
	})

	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "database error reported")
		return
	}

	ResJson := models.User{
		ID:           userInfo.ID,
		CreatedAt:    userInfo.CreatedAt,
		UpdatedAt:    userInfo.UpdatedAt,
		Email:        userInfo.Email,
		IsChirpyRed:  userInfo.IsChirpyRed.Bool,
		Token:        newToken,
		RefreshToken: newRefreshToken,
	}
	writeSuccessResponse(w, http.StatusOK, ResJson)
}

func UpdateAccessToken(api *middleware.ApiConfig, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, "unauthorized user, no token")
		return
	}

	tokenDetails, err := api.Db.GetUserFromRefreshToken(r.Context(), tokenString)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Error on database. Token not found: %v", err)
			writeErrorResponse(w, http.StatusUnauthorized, "error: user token does not exist")
			return
		}
		log.Printf("Error on database: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "database error reported")
		return
	}

	if tokenDetails.RevokedAt.Valid {
		writeErrorResponse(w, http.StatusUnauthorized, "unauthorized user, token revoked")
		return
	}

	if time.Now().After(tokenDetails.ExpiresAt) {
		writeErrorResponse(w, http.StatusUnauthorized, "unauthorized user, expired token")
		return
	}

	expiresIn := 1 * time.Hour
	newAccessToken, err := auth.MakeJWT(tokenDetails.UserID, api.Token, expiresIn)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "access token creation failed")
		return
	}
	ResJson := models.Token{
		Token: newAccessToken,
	}
	writeSuccessResponse(w, http.StatusOK, ResJson)
}

func RevokeRefreshToken(api *middleware.ApiConfig, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, "no token found or incorrect token string")
		return
	}

	err = api.Db.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
		Token:     tokenString,
	})
	if err != nil {
		log.Printf("Error on database: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "database error reported")
		return
	}
	userID, err := api.Db.GetUserFromRefreshToken(r.Context(), tokenString)
	if err != nil {
		log.Printf("Token revoked successfully, but couldn't retrieve user ID for logging: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "database error reported")
		return
	}
	log.Printf("Token for user '%v' revoked", userID.UserID)
	writeSuccessResponse(w, http.StatusNoContent, "")
}

func UpdateUser(api *middleware.ApiConfig, w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	w.Header().Set("Content-Type", "application/json")

	params := reqParams{}
	errDecode := decodeJSONBody(r, &params)

	if errDecode != nil {
		log.Printf("Error decoding parameters: %s", errDecode)
		writeErrorResponse(w, http.StatusBadRequest, "error decoding JSON")
		return
	}

	if params.Password == "" || params.Email == "" {
		writeErrorResponse(w, http.StatusBadRequest, "invalid request both password and email must be updated")
		return
	}

	newHash, err := auth.HashPassword(params.Password)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "error hashing password")
		return
	}

	res, err := api.Db.UpdateUserPasswordEmail(r.Context(), database.UpdateUserPasswordEmailParams{
		HashedPassword: newHash,
		Email:          params.Email,
		UpdatedAt:      time.Now(),
	})

	if err != nil {
		log.Printf("Error on database: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "database error reported")
		return
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, "missing token")
		return
	}

	_, err = auth.ValidateJWT(tokenString, api.Token)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, "invalid token")
		return
	}

	resJson := models.User{
		ID:          res.ID,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
		Email:       res.Email,
		IsChirpyRed: res.IsChirpyRed.Bool,
	}

	writeSuccessResponse(w, http.StatusOK, resJson)
}

func DeleteChirp(api *middleware.ApiConfig, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, "missing token")
		return
	}

	userID, err := auth.ValidateJWT(tokenString, api.Token)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, "invalid token")
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		log.Printf("Error on database: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "database error reported")
		return
	}

	chirpDetails, err := api.Db.GetChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Error on database: %v", err)
			writeErrorResponse(w, http.StatusNotFound, "error: Chirp ID does not exist")
			return
		}
		log.Printf("Error on database: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "database error reported")
		return
	}

	if chirpDetails.UserID != userID {
		writeErrorResponse(w, http.StatusForbidden, "user not authorized to delete chirp")
		return
	}

	err = api.Db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		log.Printf("Error on database: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "database error reported")
		return
	}

	writeSuccessResponse(w, http.StatusNoContent, "")

}

func UpdateChirpyRed(api *middleware.ApiConfig, w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	w.Header().Set("Content-Type", "application/json")
	params := reqParams{}
	errDecode := decodeJSONBody(r, &params)

	if errDecode != nil {
		log.Printf("Error decoding parameters: %s", errDecode)
		writeErrorResponse(w, http.StatusBadRequest, "error decoding JSON")
		return
	}

	if params.Event != "user.upgraded" {
		writeErrorResponse(w, http.StatusNoContent, "")
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		log.Printf("Error on database: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "database error reported")
		return
	}

	err = api.Db.UpdateChirpyRed(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Error on database: %v", err)
			writeErrorResponse(w, http.StatusNotFound, "error: Chirp ID does not exist")
			return
		}
		log.Printf("Error on database: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "database error reported")
		return
	}

	log.Printf("User %s upgraded to Chirpy Red.", userID)
	writeSuccessResponse(w, http.StatusNoContent, "")
}
