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
	// Dodajemo keširanje da Vercel pamti odgovor i isporučuje ga trenutno!
	w.Header().Set("Cache-Control", "public, s-maxage=60, stale-while-revalidate=30")

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SECRET_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Nedostaju varijable"})
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
		return
	}

	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", supabaseKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var recipes []Recipe
	if err := json.NewDecoder(resp.Body).Decode(&recipes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if querySlug != "" {
		if len(recipes) == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Nema recepta"})
			return
		}
		json.NewEncoder(w).Encode(recipes[0])
		return
	}

	json.NewEncoder(w).Encode(recipes)
}
