package handlers

import (
	"groupie-tracker-visualizations/internal/clients"
	"groupie-tracker-visualizations/internal/models"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type PageData struct {
	Artist             models.Artist
	Relations          models.Relation
	Dates              models.Date
	Locations          models.Location
	FormattedLocations []string
	FormattedDates     []string
	FormattedRelations map[string][]string
}

type IndexPageData struct {
	Artists   []models.Artist
	HeroImage string
}

type Errormes struct {
	Message string
}

var (
	heroImages = []string{
		"/static/1.jpg",
		"/static/2.jpg",
	}
	apiClient *clients.APIClient
)

func init() {
	apiClient = clients.NewAPIClient()
}

func findArtistByID(artists []models.Artist, id int) (models.Artist, bool) {
	for _, artist := range artists {
		if artist.ID == id {
			return artist, true
		}
	}
	return models.Artist{}, false
}

func findLocationByID(locations []models.Location, id int) models.Location {
	for _, location := range locations {
		if location.ID == id {
			return location
		}
	}
	return models.Location{}
}

func findDateByID(dates []models.Date, id int) models.Date {
	for _, date := range dates {
		if date.ID == id {
			return date
		}
	}
	return models.Date{}
}

func findRelationByID(relations []models.Relation, id int) models.Relation {
	for _, relation := range relations {
		if relation.ID == id {
			return relation
		}
	}
	return models.Relation{}
}

func titleWord(word string) string {
	if word == "" {
		return word
	}
	lower := strings.ToLower(word)
	return strings.ToUpper(string(lower[0])) + lower[1:]
}

func formatLocationString(raw string) string {
	if raw == "" {
		return raw
	}

	parts := strings.Split(raw, "-")
	for i, part := range parts {
		part = strings.ReplaceAll(part, "_", " ")
		words := strings.Fields(part)
		for j, word := range words {
			words[j] = titleWord(word)
		}
		parts[i] = strings.Join(words, " ")
	}
	return strings.Join(parts, " - ")
}

func formatLocationList(location models.Location) []string {
	formatted := make([]string, 0, len(location.Locations))
	for _, loc := range location.Locations {
		formatted = append(formatted, formatLocationString(loc))
	}
	return formatted
}

func formatDatesList(dates models.Date) []string {
	formatted := make([]string, 0, len(dates.Dates))
	for _, date := range dates.Dates {
		formatted = append(formatted, strings.TrimPrefix(date, "*"))
	}
	return formatted
}

func formatRelationsMap(relation models.Relation) map[string][]string {
	formatted := make(map[string][]string, len(relation.DatesLocations))
	for location, dates := range relation.DatesLocations {
		formatted[formatLocationString(location)] = dates
	}
	return formatted
}

func HealthHandler(c *fiber.Ctx) error {
	return c.SendString("OK")
}

func HomeHandler(c *fiber.Ctx) error {
	if c.Path() != "/" {
		log.Printf("Page Not Found")
		return showErrorPage(c, fiber.StatusNotFound, "404: Page Not Found")
	}
	if c.Method() != fiber.MethodGet {
		log.Printf("Method Not Allowed")
		return showErrorPage(c, fiber.StatusMethodNotAllowed, "405: Method Not Allowed")
	}

	artists, err := apiClient.FetchArtists()
	if err != nil {
		log.Printf("Error fetching artists: %v", err)
		return showErrorPage(c, fiber.StatusInternalServerError, "500: Internal Server Error")
	}

	data := IndexPageData{
		Artists:   artists,
		HeroImage: heroImages[rand.Intn(len(heroImages))],
	}

	return c.Render("index", data)
}

func ArtistHandler(c *fiber.Ctx) error {
	id := c.Query("id")
	if c.Method() != fiber.MethodGet {
		return showErrorPage(c, fiber.StatusMethodNotAllowed, "405: Method not allowed")
	}
	if id == "" {
		return showErrorPage(c, fiber.StatusBadRequest, "400: Bad Request")
	}
	artistID, err := strconv.Atoi(id)
	if err != nil || artistID <= 0 {
		return showErrorPage(c, fiber.StatusBadRequest, "400: Invalid artist id")
	}

	artists, err := apiClient.FetchArtists()
	if err != nil {
		log.Printf("Error fetching artists: %v", err)
		return showErrorPage(c, fiber.StatusInternalServerError, "500: Internal Server Error")
	}

	artist, found := findArtistByID(artists, artistID)
	if !found {
		return showErrorPage(c, fiber.StatusNotFound, "404: Artist Not Found")
	}

	locationsList, err := apiClient.FetchLocations()
	if err != nil {
		log.Printf("Error fetching locations: %v", err)
		return showErrorPage(c, fiber.StatusInternalServerError, "500: Internal Server Error")
	}

	datesList, err := apiClient.FetchDates()
	if err != nil {
		log.Printf("Error fetching dates: %v", err)
		return showErrorPage(c, fiber.StatusInternalServerError, "500: Internal Server Error")
	}

	relationsList, err := apiClient.FetchRelation()
	if err != nil {
		log.Printf("Error fetching relations: %v", err)
		return showErrorPage(c, fiber.StatusInternalServerError, "500: Internal Server Error")
	}

	artistLocations := findLocationByID(locationsList, artistID)
	artistDates := findDateByID(datesList, artistID)
	artistRelation := findRelationByID(relationsList, artistID)

	pageData := PageData{
		Artist:             artist,
		Locations:          artistLocations,
		Dates:              artistDates,
		Relations:          artistRelation,
		FormattedLocations: formatLocationList(artistLocations),
		FormattedDates:     formatDatesList(artistDates),
		FormattedRelations: formatRelationsMap(artistRelation),
	}

	if c.Get("X-Requested-With") == "XMLHttpRequest" {
		return c.JSON(pageData)
	}

	return c.Render("artist", pageData)
}

func showErrorPage(c *fiber.Ctx, statusCode int, message string) error {
	c.Status(statusCode)
	errmess := Errormes{Message: message}
	return c.Render("error", errmess)
}
