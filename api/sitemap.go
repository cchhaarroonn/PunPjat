package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const siteURL = "https://pun-pjat.vercel.app"

type sitemapRecipe struct {
	Slug string `json:"slug"`
}

func Sitemap(w http.ResponseWriter, r *http.Request) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SECRET_KEY")

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, s-maxage=3600, stale-while-revalidate=600")

	if supabaseURL == "" || supabaseKey == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	apiURL := fmt.Sprintf("%s/rest/v1/recipes?select=slug", supabaseURL)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", supabaseKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var recipes []sitemapRecipe
	if err := json.NewDecoder(resp.Body).Decode(&recipes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	sb.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")

	sb.WriteString(fmt.Sprintf("  <url><loc>%s/</loc><priority>1.0</priority></url>\n", siteURL))

	for _, rec := range recipes {
		sb.WriteString(fmt.Sprintf(
			"  <url><loc>%s/post.html?slug=%s</loc><priority>0.8</priority></url>\n",
			siteURL, rec.Slug,
		))
	}

	sb.WriteString(`</urlset>`)

	w.Write([]byte(sb.String()))
}
