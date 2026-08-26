package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
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
		json.NewEncoder(w).Encode(map[string]string{"error": "Nedostaju environment varijable"})
		return
	}

	querySlug := r.URL.Query().Get("slug")

	u, err := url.Parse(fmt.Sprintf("%s/rest/v1/recipes", supabaseURL))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Nevaljan URL baze"})
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
		json.NewEncoder(w).Encode(map[string]string{"error": "Greška zahtjeva"})
		return
	}

	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", supabaseKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Greška spajanja na Supabase: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Greška čitanja odgovora"})
		return
	}

	// Supabase ne vraća niz kad nešto pođe po zlu (kriv apikey, RLS,
	// nepostojeća tablica...) nego JSON objekt s porukom greške.
	// Prije se to tiho pokušavalo dekodirati kao []Recipe i puklo bi
	// s generičnom "Greška dekodiranja" bez pravog razloga.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{
			"error":            "Supabase je vratio grešku",
			"supabase_status":  fmt.Sprintf("%d", resp.StatusCode),
			"supabase_message": string(body),
		})
		return
	}

	var recipes []Recipe
	if err := json.Unmarshal(body, &recipes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Greška dekodiranja: " + err.Error()})
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
