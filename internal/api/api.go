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
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	type responseJSON struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    string    `json:"user_id"`
	}

	w.Header().Set("Content-Type", "application/json")
	params := validateBody{}
	errDecode := decodeJSONBody(r, &params)

	var ResJson any
	var StatusCode = http.StatusOK
	switch {
	case errDecode != nil:
		{
			log.Printf("Error decoding parameters: %s", errDecode)
			StatusCode = http.StatusBadRequest
			ResJson = models.ErrorBody{
				Error: "error decoding JSON",
			}
		}
	case len(params.Body) > 140:
		{
			StatusCode = http.StatusBadRequest
			ResJson = models.ErrorBody{
				Error: "Chirp is too long",
			}
		}
	case params.Body == "":
		{
			StatusCode = http.StatusBadRequest
			ResJson = models.ErrorBody{
				Error: "Chirp must contain characters",
			}
		}
	default:
		{
			userId, err := uuid.Parse(params.UserID)
			if err != nil {
				log.Printf("Invalid user_id format: %v", err)
				StatusCode = http.StatusBadRequest
				ResJson = models.ErrorBody{
					Error: "invalid user_id format",
				}
			} else {
				res, err := api.Db.CreateChirp(r.Context(), database.CreateChirpParams{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Body:      profanityFilter(params.Body),
					UserID:    userId,
				})
				if err != nil {
					log.Printf("Error on database: %v", err)
					StatusCode = http.StatusInternalServerError
					ResJson = models.ErrorBody{
						Error: "database error reported",
					}
				} else {
					StatusCode = http.StatusCreated
					ResJson = responseJSON{
						ID:        res.ID.String(),
						CreatedAt: res.CreatedAt,
						UpdatedAt: res.UpdatedAt,
						Body:      res.Body,
						UserID:    res.UserID.String(),
					}
				}
			}
		}
	}

	w.WriteHeader(StatusCode)
	err := marshalJSON(w, ResJson)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.Write([]byte(`{"error: Something went wrong converting JSON"}`))
		return
	}

}

func GetAllChirps(api *middleware.ApiConfig, w http.ResponseWriter, r *http.Request) {

	type responseJSON struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    string    `json:"user_id"`
	}

	w.Header().Set("Content-Type", "application/json")
	var StatusCode = http.StatusOK

	res, err := api.Db.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Error on database: %v", err)
		StatusCode = http.StatusInternalServerError
		errorRes := models.ErrorBody{
			Error: "database error reported",
		}
		w.WriteHeader(StatusCode)
		err = marshalJSON(w, errorRes)
		if err != nil {
			log.Printf("Error on database: %v", err)
			StatusCode = http.StatusInternalServerError
			errorRes := models.ErrorBody{
				Error: "database error reported",
			}
			w.WriteHeader(StatusCode)
			err = marshalJSON(w, errorRes)
			if err != nil {
				log.Printf("Error marshalling JSON: %s", err)
				w.Write([]byte(`{"error: Something went wrong converting JSON"}`))
				return
			}
			return
		}
	}

	ResJson := []responseJSON{}
	for _, chirp := range res {
		json := responseJSON{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
		}
		ResJson = append(ResJson, json)
	}

	w.WriteHeader(StatusCode)
	err = marshalJSON(w, ResJson)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.Write([]byte(`{"error: Something went wrong converting JSON"}`))
		return
	}
}

func GetChirp(api *middleware.ApiConfig, w http.ResponseWriter, r *http.Request) {

	type responseJSON struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    string    `json:"user_id"`
	}

	w.Header().Set("Content-Type", "application/json")
	var StatusCode = http.StatusOK

	chirpID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		log.Printf("chirpID: %v", chirpID)
		log.Printf("Error on database: %v", err)
		StatusCode = http.StatusInternalServerError
		errorRes := models.ErrorBody{
			Error: "database error reported",
		}
		w.WriteHeader(StatusCode)
		err = marshalJSON(w, errorRes)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.Write([]byte(`{"error: Something went wrong converting JSON"}`))
			return
		}
		return
	}

	res, err := api.Db.GetChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Error on database: %v", err)
			StatusCode = http.StatusNotFound
			w.WriteHeader(StatusCode)
			w.Write([]byte(`{"error: Chirp ID does not exist"}`))
			return
		}
		log.Printf("Error on database: %v", err)
		StatusCode = http.StatusInternalServerError
		errorRes := models.ErrorBody{
			Error: "database error reported",
		}
		w.WriteHeader(StatusCode)
		err = marshalJSON(w, errorRes)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.Write([]byte(`{"error: Something went wrong converting JSON"}`))
			return
		}
		return
	}

	ResJson := responseJSON{
		ID:        res.ID.String(),
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		Body:      res.Body,
		UserID:    res.UserID.String(),
	}

	w.WriteHeader(StatusCode)
	err = marshalJSON(w, ResJson)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.Write([]byte(`{"error: Something went wrong converting JSON"}`))
		return
	}
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

	var ResJson any
	var StatusCode = http.StatusOK
	switch {
	case errDecode != nil:
		{
			log.Printf("Error decoding parameters: %s", errDecode)
			StatusCode = http.StatusBadRequest
			ResJson = models.ErrorBody{
				Error: "error decoding JSON",
			}
		}
	//case for minimum possible email length a@a.aa
	case len(params.Email) < 5:
		{
			StatusCode = http.StatusBadRequest
			ResJson = models.ErrorBody{
				Error: "invalid email submitted",
			}
		}
	//basic email formatting checks
	case !validateEmail(params.Email):
		{
			StatusCode = http.StatusBadRequest
			ResJson = models.ErrorBody{
				Error: "invalid email submitted",
			}
		}
	default:
		pwd, errHash := auth.HashPassword(params.Password)
		if errHash != nil {
			ResJson = models.ErrorBody{
				Error: "error hashing password",
			}
		} else {
			res, err := api.Db.CreateUser(r.Context(), database.CreateUserParams{
				ID:             uuid.New(),
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
				Email:          params.Email,
				HashedPassword: pwd,
			})
			if err != nil {
				log.Printf("Error on database: %v", err)
				StatusCode = http.StatusInternalServerError
				ResJson = models.ErrorBody{
					Error: "database error reported",
				}
			} else {
				StatusCode = http.StatusCreated
				ResJson = models.User{
					ID:        res.ID,
					CreatedAt: res.CreatedAt,
					UpdatedAt: res.UpdatedAt,
					Email:     res.Email,
				}
				log.Printf("User: %s created with ID %v", res.Email, res.ID)
			}
		}
	}

	w.WriteHeader(StatusCode)
	err := marshalJSON(w, ResJson)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.Write([]byte(`{"error: Something went wrong converting JSON"}`))
		return
	}
}

func UserLogin(api *middleware.ApiConfig, w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	w.Header().Set("Content-Type", "application/json")
	params := reqParams{}
	errDecode := decodeJSONBody(r, &params)

	var ResJson any
	var StatusCode = http.StatusOK
	switch {
	case errDecode != nil:
		{
			log.Printf("Error decoding parameters: %s", errDecode)
			StatusCode = http.StatusBadRequest
			ResJson = models.ErrorBody{
				Error: "error decoding JSON",
			}
		}
	default:
		userInfo, err := api.Db.GetUserPassword(r.Context(), params.Email)
		if err != nil {
			log.Printf("Error on database: %v", err)
			StatusCode = http.StatusInternalServerError
			ResJson = models.ErrorBody{
				Error: "database error reported",
			}
		} else {
			err = auth.CheckPasswordHash(userInfo.HashedPassword, params.Password)
			if err != nil {
				StatusCode = http.StatusUnauthorized
				ResJson = models.ErrorBody{
					Error: "Incorrect email or password",
				}
			} else {
				StatusCode = http.StatusOK
				ResJson = models.User{
					ID:        userInfo.ID,
					CreatedAt: userInfo.CreatedAt,
					UpdatedAt: userInfo.UpdatedAt,
					Email:     userInfo.Email,
				}
			}
		}
	}

	w.WriteHeader(StatusCode)
	err := marshalJSON(w, ResJson)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.Write([]byte(`{"error: Something went wrong converting JSON"}`))
		return
	}
}
