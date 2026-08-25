package handler

import (
	"encoding/json"
	"net/http"
)

// Struktura recepta
type Recipe struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
	Steps       []string `json:"steps"`
}

// Glavna Go funkcija za API
func Handler(w http.ResponseWriter, r *http.Request) {
	// CORS i JSON zaglavlja
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Lažna baza podataka s receptima
	recipes := []Recipe{
		{
			Slug:        "domaci-cevapi",
			Title:       "Sočni recept za domaće ćevape",
			Category:    "Glavna jela",
			Description: "Nauči kako pripremiti savršene, mekane ćevape s domaćim somunom i lukom.",
			Ingredients: []string{
				"500g mljevene junetine",
				"1 žličica sode bikarbone",
				"Sol i papar po ukusu",
				"1 dcl mineralne vode",
			},
			Steps: []string{
				"Meso dobro izmiješajte sa sodom bikarbonom, solju, paprom i mineralnom vodom.",
				"Ostavite smjesu u hladnjaku minimalno 3 sata, najbolje preko noći.",
				"Oblikujte ćevape pomoću kalupa ili ruku.",
				"Pecite na roštilju ili dobro zagrijanoj tavi s malo ulja dok ne dobiju lijepu boju.",
			},
		},
		{
			Slug:        "palacinke-s-cokoladom",
			Title:       "Fluffy palačinke s domaćom čokoladom",
			Category:    "Slastice",
			Description: "Meke i prozračne palačinke koje uspijevaju baš svaki put.",
			Ingredients: []string{
				"2 jaja",
				"300ml brašna",
				"200ml mlijeka i 100ml mineralne vode",
				"Prstohvat soli i žličica šećera",
			},
			Steps: []string{
				"Izmutite jaja, dodajte mlijeko i mineralnu vodu.",
				"Postupno dodavajte brašno uz miješanje da nema grudica.",
				"Pecite na vrućoj tavi s malo ulja s obje strane.",
			},
		},
	}

	// Provjeravamo traži li se specifičan recept preko slug-a (npr. /api/post?slug=domaci-cevapi)
	querySlug := r.URL.Query().Get("slug")

	if querySlug != "" {
		for _, recipe := range recipes {
			if recipe.Slug == querySlug {
				json.NewEncoder(w).Encode(recipe)
				return
			}
		}
		// Ako recept nije pronađen
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Recept nije pronađen"})
		return
	}

	// Ako nema slug-a, vraća listu svih recepata
	json.NewEncoder(w).Encode(recipes)
}
