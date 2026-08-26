package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	ImageURL    string   `json:"image_url"`
	PrepTime    string   `json:"prep_time"`
	Difficulty  string   `json:"difficulty"`
	Author      string   `json:"author"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// CORS i Content-Type zaglavlja moraju biti prva stvar
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SECRET_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Nedostaju environment varijable na Vercelu"})
		return
	}

	querySlug := r.URL.Query().Get("slug")

	u, err := url.Parse(fmt.Sprintf("%s/rest/v1/recipes", supabaseURL))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Nevaljan Supabase URL"})
		return
	}

	q := u.Query()
	q.Set("select", "*")
	if querySlug != "" {
		q.Set("slug", fmt.Sprintf("eq.%s", querySlug))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
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

	var recipes []Recipe
	if err := json.NewDecoder(resp.Body).Decode(&recipes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Greška pri dekodiranju podataka iz baze"})
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
