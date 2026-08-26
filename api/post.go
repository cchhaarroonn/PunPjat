package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

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

	// Dohvaćamo varijable
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SECRET_KEY")

	// Ako slučajno varijable nisu postavljene na Vercelu, spriječit ćemo pad i vratiti jasnu poruku
	if supabaseURL == "" || supabaseKey == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Nedostaju Supabase environment varijable na Vercelu!"})
		return
	}

	querySlug := r.URL.Query().Get("slug")

	apiURL := fmt.Sprintf("%s/rest/v1/recipes?select=*", supabaseURL)
	if querySlug != "" {
		apiURL = fmt.Sprintf("%s&slug=eq.%s", apiURL, querySlug)
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Greška pri kreiranju zahtjeva"})
		return
	}

	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", supabaseKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Greška pri spajanju na bazu"})
		return
	}
	defer resp.Body.Close()

	// Provjeravamo je li Supabase vratio grešku
	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(map[string]string{"error": "Supabase je odbio zahtjev, provjeri ključeve"})
		return
	}

	var recipes []Recipe
	if err := json.NewDecoder(resp.Body).Decode(&recipes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Greška pri dekodiranju podataka"})
		return
	}

	if querySlug != "" {
		if len(recipes) == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Recept nije pronađen"})
			return
		}
		json.NewEncoder(w).Encode(recipes[0])
		return
	}

	json.NewEncoder(w).Encode(recipes)
}
