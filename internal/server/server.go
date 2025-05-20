package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Walther-Knight/chirpy/internal/middleware"
)

func Start() error {
	newMux := http.NewServeMux()
	httpSrv := &http.Server{
		Handler: newMux,
		Addr:    ":8080",
	}
	var hitCount middleware.ApiConfig
	log.Println("Starting handlers...")
	//admin functions
	newMux.HandleFunc("GET /api/healthz", health)
	newMux.HandleFunc("GET /admin/metrics", hitCount.HitTotal)
	newMux.HandleFunc("POST /admin/reset", hitCount.HitReset)
	//application functions
	newMux.Handle("/app/", hitCount.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./static")))))
	newMux.HandleFunc("POST /api/validate_chirp", validate)

	log.Printf("Starting http server on %s\n", httpSrv.Addr)
	return httpSrv.ListenAndServe()
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func validate(w http.ResponseWriter, r *http.Request) {
	type validateBody struct {
		Body string `json:"body"`
	}

	type errorBody struct {
		Error string `json:"error"`
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
	default:
		{
			w.WriteHeader(http.StatusOK)
			ResJson = validRes{
				Valid: true,
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
