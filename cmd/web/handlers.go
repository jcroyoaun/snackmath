package main

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/jcroyoaun/snackmath/assets"
	"github.com/jcroyoaun/snackmath/internal/response"
)

type FoodItem struct {
	Name         string
	PortionSize  float64
	Calories     float64
	Protein      float64
	Carbs        float64
	Fat          float64
	CalsPer100g  float64
	ProtPer100g  float64
	CarbsPer100g float64
	FatPer100g   float64
}

type CompareResult struct {
	ItemA      FoodItem
	ItemB      FoodItem
	Winner     string
	WinnerItem string
	HasMacros  bool
}

func (app *application) serviceWorker(w http.ResponseWriter, r *http.Request) {
	data, err := assets.EmbeddedFiles.ReadFile("static/sw.js")
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/javascript")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(data)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	q := r.URL.Query()
	hasParams := q.Get("grams_a") != "" && q.Get("grams_b") != "" &&
		q.Get("cal_a") != "" && q.Get("cal_b") != ""

	if hasParams {
		itemA := parseFoodItem(q, "a")
		itemB := parseFoodItem(q, "b")

		if itemA.PortionSize > 0 && itemB.PortionSize > 0 {
			normalize(&itemA)
			normalize(&itemB)

			result := buildResult(itemA, itemB)
			data["Result"] = result
		}
	}

	err := response.Page(w, http.StatusOK, data, "pages/home.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func normalize(item *FoodItem) {
	item.CalsPer100g = round1((item.Calories / item.PortionSize) * 100)
	item.ProtPer100g = round1((item.Protein / item.PortionSize) * 100)
	item.CarbsPer100g = round1((item.Carbs / item.PortionSize) * 100)
	item.FatPer100g = round1((item.Fat / item.PortionSize) * 100)
}

func buildResult(itemA, itemB FoodItem) CompareResult {
	var winner string
	if itemA.CalsPer100g < itemB.CalsPer100g {
		winner = "A"
	} else if itemB.CalsPer100g < itemA.CalsPer100g {
		winner = "B"
	} else {
		winner = "TIE"
	}

	var winnerItem string
	switch winner {
	case "A":
		winnerItem = itemA.Name
	case "B":
		winnerItem = itemB.Name
	default:
		winnerItem = "It's a tie"
	}

	hasMacros := itemA.Protein > 0 || itemB.Protein > 0 ||
		itemA.Carbs > 0 || itemB.Carbs > 0 ||
		itemA.Fat > 0 || itemB.Fat > 0

	return CompareResult{
		ItemA:      itemA,
		ItemB:      itemB,
		Winner:     winner,
		WinnerItem: winnerItem,
		HasMacros:  hasMacros,
	}
}

type queryGetter interface {
	Get(key string) string
}

func parseFoodItem(q queryGetter, prefix string) FoodItem {
	name := strings.TrimSpace(q.Get("name_" + prefix))
	grams := parseNumber(q.Get("grams_" + prefix))
	calories := parseNumber(q.Get("cal_" + prefix))
	protein := parseNumber(q.Get("prot_" + prefix))
	carbs := parseNumber(q.Get("carbs_" + prefix))
	fat := parseNumber(q.Get("fat_" + prefix))

	if name == "" {
		if prefix == "a" {
			name = "Item A"
		} else {
			name = "Item B"
		}
	}

	return FoodItem{
		Name:        name,
		PortionSize: grams,
		Calories:    calories,
		Protein:     protein,
		Carbs:       carbs,
		Fat:         fat,
	}
}

func parseNumber(s string) float64 {
	s = strings.TrimSpace(s)
	end := len(s)
	for end > 0 && !unicode.IsDigit(rune(s[end-1])) && s[end-1] != '.' {
		end--
	}
	if end == 0 {
		return 0
	}
	start := 0
	for start < end && !unicode.IsDigit(rune(s[start])) && s[start] != '.' {
		start++
	}
	f, _ := strconv.ParseFloat(s[start:end], 64)
	return f
}

func round1(f float64) float64 {
	return math.Round(f*10) / 10
}
