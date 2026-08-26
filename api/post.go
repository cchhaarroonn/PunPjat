package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Struktura recepta prilagođena PostgreSQL tablici
type Recipe struct {
	ID          int      `json:"id"`
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
	Steps       []string `json:"steps"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Supabase podaci iz environment varijabli (ili direktno upisani za test)
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SECRET_KEY")

	// Fallback ako zaboraviš postaviti na Vercelu (za lakši test)
	if supabaseURL == "" {
		supabaseURL = "https://ldvpfbkhehtrbbwfklmx.supabase.co"
	}

	querySlug := r.URL.Query().Get("slug")

	// Gradimo URL prema Supabase REST API-ju za tablicu "recipes"
	apiURL := fmt.Sprintf("%s/rest/v1/recipes?select=*", supabaseURL)
	if querySlug != "" {
		apiURL = fmt.Sprintf("%s&slug=eq.%s", apiURL, querySlug)
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Greška kod kreiranja zahtjeva prema bazi"})
		return
	}

	// Postavljamo Supabase zaglavlja
	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", supabaseKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ne mogu se spojiti na bazu"})
		return
	}
	defer resp.Body.Close()

	var recipes []Recipe
	if err := json.NewDecoder(resp.Body).Decode(&recipes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Greška kod čitanja podataka iz baze"})
		return
	}

	// Ako je tražen specifičan slug
	if querySlug != "" {
		if len(recipes) == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Recept nije pronađen"})
			return
		}
		json.NewEncoder(w).Encode(recipes[0])
		return
	}

	// Vrati sve recepte
	json.NewEncoder(w).Encode(recipes)
}
