package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

	type validRes struct {
		Valid bool `json:"valid"`
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
