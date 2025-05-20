package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

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

func profanityFilter(s string) (string, bool) {
	splitString := strings.Split(s, " ")
	var cleanString []string
	filtered := false
	for _, word := range splitString {
		switch {
		case strings.ToLower(word) == "kerfuffle":
			cleanString = append(cleanString, "****")
			filtered = true
		case strings.ToLower(word) == "sharbert":
			cleanString = append(cleanString, "****")
			filtered = true
		case strings.ToLower(word) == "fornax":
			cleanString = append(cleanString, "****")
			filtered = true
		default:
			cleanString = append(cleanString, word)
		}
	}
	return strings.Join(cleanString, " "), filtered
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func Validate(w http.ResponseWriter, r *http.Request) {
	type validateBody struct {
		Body string `json:"body"`
	}

	type errorBody struct {
		Error string `json:"error"`
	}

	type cleanedBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	w.Header().Set("Content-Type", "application/json")
	params := validateBody{}
	errDecode := decodeJSONBody(r, &params)
	if errDecode != nil {
		log.Printf("Error decoding parameters: %s", errDecode)
		w.WriteHeader(http.StatusBadRequest)
		decodeError := errorBody{
			Error: "Something went wrong",
		}
		marshalJSON(w, decodeError)
		return
	}

	var ResJson any
	cleanBody, filter := profanityFilter(params.Body)
	switch {
	case len(params.Body) > 140:
		{
			w.WriteHeader(http.StatusBadRequest)
			ResJson = errorBody{
				Error: "Chirp is too long",
			}
		}
	case params.Body == "":
		{
			w.WriteHeader(http.StatusBadRequest)
			ResJson = errorBody{
				Error: "Something went wrong",
			}
		}
	case filter:
		w.WriteHeader(http.StatusOK)
		ResJson = cleanedBody{
			CleanedBody: cleanBody,
		}
	default:
		{
			w.WriteHeader(http.StatusOK)
			ResJson = cleanedBody{
				CleanedBody: cleanBody,
			}
		}
	}

	err := marshalJSON(w, ResJson)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.Write([]byte(`{"error":"Something went wrong"}`))
		return
	}

}

func NewUser(api *middleware.ApiConfig, w http.ResponseWriter, r *http.Request) {
	type validEmail struct {
		Email string `json:"email"`
	}

	type errorBody struct {
		Error string `json:"error"`
	}

	w.Header().Set("Content-Type", "application/json")
	params := validEmail{}
	errDecode := decodeJSONBody(r, &params)
	if errDecode != nil {
		log.Printf("Error decoding parameters: %s", errDecode)
		w.WriteHeader(http.StatusBadRequest)
		decodeError := errorBody{
			Error: "Something went wrong",
		}
		marshalJSON(w, decodeError)
		return
	}

	var ResJson any
	switch {
	//case for minimum possible email length a@a.aa
	case len(params.Email) < 5:
		{
			w.WriteHeader(http.StatusBadRequest)
			ResJson = errorBody{
				Error: "invalid email submitted",
			}
		}
	case strings.Index(params.Email, "@") < strings.Index(params.Email, ".") && strings.Count(params.Email, "@") == 1:
		res, err := api.Db.CreateUser(r.Context(), database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Email:     params.Email,
		})
		if err != nil {
			log.Printf("Error on database: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"database error reported"}`))
			return
		}
		w.WriteHeader(http.StatusCreated)
		ResJson = models.User{
			ID:        res.ID,
			CreatedAt: res.CreatedAt,
			UpdatedAt: res.UpdatedAt,
			Email:     res.Email,
		}
		log.Printf("User: %s created with ID %v", res.Email, res.ID)
	default:
		w.WriteHeader(http.StatusBadRequest)
		ResJson = errorBody{
			Error: "invalid email submitted",
		}

	}

	marshalJSON(w, ResJson)
}
