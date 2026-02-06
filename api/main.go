package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/AaronLieb/octagon/brackets"
	"github.com/AaronLieb/octagon/cache"
	"github.com/AaronLieb/octagon/characters"
	"github.com/AaronLieb/octagon/conflicts"
	"github.com/AaronLieb/octagon/seeding"
	"github.com/AaronLieb/octagon/startgg"
	"github.com/AaronLieb/octagon/tournament"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Attendee struct {
	ID        string `json:"id"`
	GamerTag  string `json:"gamerTag"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PlayerID  string `json:"playerId"`
}

type SeedRequest struct {
	Tournament string `json:"tournament"`
	Redemption bool   `json:"redemption"`
}

type PublishRequest struct {
	Tournament string       `json:"tournament"`
	Redemption bool         `json:"redemption"`
	Players    []SeedResult `json:"players"`
}

type ConflictRequest struct {
	Player1    string  `json:"player1"`
	Player2    string  `json:"player2"`
	Reason     string  `json:"reason"`
	Priority   int     `json:"priority"`
	Expiration *string `json:"expiration"`
}

type ConflictResponse struct {
	Player1    string `json:"player1"`
	Player2    string `json:"player2"`
	Reason     string `json:"reason"`
	Priority   int    `json:"priority"`
	Expiration string `json:"expiration"`
}

type SetResponse struct {
	ID       int    `json:"id"`
	Player1  Player `json:"player1"`
	Player2  Player `json:"player2"`
	Round    string `json:"round"`
	Entrant1 int    `json:"entrant1"`
	Entrant2 int    `json:"entrant2"`
	P1Char   string `json:"p1Char,omitempty"`
	P2Char   string `json:"p2Char,omitempty"`
}

type Player struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type GameResult struct {
	Winner int    `json:"winner"`
	P1Char string `json:"p1Char"`
	P2Char string `json:"p2Char"`
}

type ReportSetRequest struct {
	SetID int          `json:"setId"`
	Games []GameResult `json:"games"`
}

type SeedResult struct {
	Name   string  `json:"name"`
	Rating float64 `json:"rating"`
	Seed   int     `json:"seed"`
	ID     string  `json:"id"`
}

func main() {
	// Load environment variables
	godotenv.Load()

	// Initialize cache
	cache.Open()

	r := gin.Default()

	// Enable CORS for React frontend
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/api/attendees", getAttendees)
	r.POST("/api/seed", runSeeding)
	r.POST("/api/seed/publish", publishSeeding)
	r.GET("/api/conflicts", getConflicts)
	r.POST("/api/conflicts", addConflict)
	r.DELETE("/api/conflicts/:index", deleteConflict)
	r.GET("/api/sets", getSets)
	r.GET("/api/sets/ready", getReadySets)
	r.POST("/api/sets/report", reportSet)
	r.GET("/api/characters", getCharacters)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}

func getAttendees(c *gin.Context) {
	tournamentSlug := c.DefaultQuery("tournament", "octagon")

	resp, err := startgg.GetParticipants(context.Background(), tournamentSlug)
	if err != nil {
		fmt.Printf("Error getting participants: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tournament := resp.Tournament
	participants := tournament.Participants.GetNodes()

	fmt.Printf("Tournament: %s, Participants count: %d\n", tournament.Name, len(participants))

	attendees := make([]Attendee, len(participants))
	for i, participant := range participants {
		var playerID string
		switch v := participant.Player.Id.(type) {
		case float64:
			playerID = strconv.FormatFloat(v, 'f', 0, 64)
		case int:
			playerID = strconv.Itoa(v)
		default:
			playerID = fmt.Sprintf("%v", v)
		}

		attendees[i] = Attendee{
			ID:        fmt.Sprintf("%v", participant.Id),
			GamerTag:  participant.GamerTag,
			FirstName: participant.ContactInfo.NameFirst,
			LastName:  participant.ContactInfo.NameLast,
			PlayerID:  playerID,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"tournament": tournament.Name,
		"attendees":  attendees,
	})
}

func runSeeding(c *gin.Context) {
	var req SeedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tournamentSlug := req.Tournament
	if tournamentSlug == "" {
		tournamentSlug = "octagon"
	}

	players, err := seeding.GenerateSeeding(context.Background(), tournamentSlug, req.Redemption, nil)
	if err != nil {
		fmt.Printf("Error generating seeding: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fullTournamentSlug, err := startgg.GetTournamentSlug(context.Background(), tournamentSlug)
	if err != nil {
		fmt.Printf("Error getting tournament slug: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	event := startgg.EventUltimateSingles
	if req.Redemption {
		event = startgg.EventRedemptionBracket
	}

	results := make([]SeedResult, len(players))
	for i, player := range players {
		results[i] = SeedResult{
			Name:   player.Name,
			Rating: player.Rating,
			Seed:   i + 1,
			ID:     startgg.ToString(player.ID),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"tournament": fullTournamentSlug,
		"event":      event,
		"seeds":      results,
	})
}

func publishSeeding(c *gin.Context) {
	var req PublishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tournamentSlug := req.Tournament
	if tournamentSlug == "" {
		tournamentSlug = "octagon"
	}

	fullTournamentSlug, err := startgg.GetTournamentSlug(context.Background(), tournamentSlug)
	if err != nil {
		fmt.Printf("Error getting tournament slug: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	event := startgg.EventUltimateSingles
	if req.Redemption {
		event = startgg.EventRedemptionBracket
	}

	slug := fmt.Sprintf(startgg.EventSlugFormat, fullTournamentSlug, event)

	// Convert SeedResult back to brackets.Player for publishing
	players := make([]brackets.Player, len(req.Players))
	for i, seedResult := range req.Players {
		players[i] = brackets.Player{
			Name:   seedResult.Name,
			ID:     startgg.ToID(seedResult.ID),
			Rating: seedResult.Rating,
		}
	}

	err = seeding.PublishSeeding(context.Background(), slug, players)
	if err != nil {
		fmt.Printf("Error publishing seeds: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Seeding published successfully",
	})
}

func getConflicts(c *gin.Context) {
	conflictsList := conflicts.GetConflicts([]string{})

	conflictResponses := make([]ConflictResponse, len(conflictsList))
	for i, conflict := range conflictsList {
		player1 := ""
		player2 := ""
		if len(conflict.Players) >= 2 {
			player1 = conflict.Players[0].Name
			player2 = conflict.Players[1].Name
		}

		expiration := ""
		if conflict.Expiration != nil {
			expiration = conflict.Expiration.Format("2006-01-02 15:04")
		}

		conflictResponses[i] = ConflictResponse{
			Player1:    player1,
			Player2:    player2,
			Reason:     conflict.Reason,
			Priority:   conflict.Priority,
			Expiration: expiration,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"conflicts": conflictResponses,
	})
}

func addConflict(c *gin.Context) {
	var req ConflictRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var expiration *time.Time
	if req.Expiration != nil && *req.Expiration != "" {
		parsed, err := time.Parse("2006-01-02T15:04", *req.Expiration)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiration format"})
			return
		}
		expiration = &parsed
	}

	newConflict := conflicts.Conflict{
		Players: []conflicts.Player{
			{Name: req.Player1},
			{Name: req.Player2},
		},
		Priority:   req.Priority,
		Reason:     req.Reason,
		Expiration: expiration,
	}

	err := conflicts.SaveConflict(newConflict)
	if err != nil {
		fmt.Printf("Error saving conflict: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Conflict added successfully",
	})
}

func deleteConflict(c *gin.Context) {
	indexStr := c.Param("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid index"})
		return
	}

	conflictsList := conflicts.GetConflicts([]string{})
	if index < 0 || index >= len(conflictsList) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Index out of range"})
		return
	}

	// Remove the conflict at the specified index
	updatedConflicts := append(conflictsList[:index], conflictsList[index+1:]...)

	err = conflicts.WriteConflictsFile(updatedConflicts)
	if err != nil {
		fmt.Printf("Error writing conflicts file: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Conflict deleted successfully",
	})
}

func getSets(c *gin.Context) {
	tournamentSlug := c.DefaultQuery("tournament", "octagon")
	includeCompleted := c.DefaultQuery("includeCompleted", "false") == "true"

	fullTournamentSlug, err := startgg.GetTournamentSlug(context.Background(), tournamentSlug)
	if err != nil {
		fmt.Printf("Error getting tournament slug: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch sets from both events
	events := []string{startgg.EventUltimateSingles, startgg.EventRedemptionBracket}
	var allSets []tournament.Set

	for _, event := range events {
		eventSlug := fmt.Sprintf(startgg.EventSlugFormat, fullTournamentSlug, event)
		sets, err := tournament.FetchReportableSets(context.Background(), eventSlug, includeCompleted)
		if err != nil {
			// Log error but continue with other events
			fmt.Printf("Error fetching sets for %s: %v\n", event, err)
			continue
		}
		allSets = append(allSets, sets...)
	}

	setResponses := make([]SetResponse, len(allSets))
	for i, set := range allSets {
		p1Char, _ := cache.GetPlayerCharacter(set.Player1.ID)
		p2Char, _ := cache.GetPlayerCharacter(set.Player2.ID)

		setResponses[i] = SetResponse{
			ID:       set.ID,
			Player1:  Player{Name: set.Player1.Name, ID: set.Player1.ID},
			Player2:  Player{Name: set.Player2.Name, ID: set.Player2.ID},
			Round:    set.Round,
			Entrant1: set.Entrant1,
			Entrant2: set.Entrant2,
			P1Char:   p1Char,
			P2Char:   p2Char,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"sets": setResponses,
	})
}

func getReadySets(c *gin.Context) {
	tournamentSlug := c.DefaultQuery("tournament", "octagon")

	fullTournamentSlug, err := startgg.GetTournamentSlug(context.Background(), tournamentSlug)
	if err != nil {
		fmt.Printf("Error getting tournament slug: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch sets with state 1 (ready to call) from both events
	events := []string{startgg.EventUltimateSingles, startgg.EventRedemptionBracket}
	var allSets []tournament.Set

	for _, event := range events {
		eventSlug := fmt.Sprintf(startgg.EventSlugFormat, fullTournamentSlug, event)
		sets, err := tournament.FetchReportableSets(context.Background(), eventSlug, false)
		if err != nil {
			fmt.Printf("Error fetching sets for %s: %v\n", event, err)
			continue
		}
		// Filter for state 1 only
		for _, set := range sets {
			if set.State == 1 {
				allSets = append(allSets, set)
			}
		}
	}

	setResponses := make([]SetResponse, len(allSets))
	for i, set := range allSets {
		p1Char, _ := cache.GetPlayerCharacter(set.Player1.ID)
		p2Char, _ := cache.GetPlayerCharacter(set.Player2.ID)

		setResponses[i] = SetResponse{
			ID:       set.ID,
			Player1:  Player{Name: set.Player1.Name, ID: set.Player1.ID},
			Player2:  Player{Name: set.Player2.Name, ID: set.Player2.ID},
			Round:    set.Round,
			Entrant1: set.Entrant1,
			Entrant2: set.Entrant2,
			P1Char:   p1Char,
			P2Char:   p2Char,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"sets": setResponses,
	})
}

func reportSet(c *gin.Context) {
	var req ReportSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to tournament.GameResult format and get character IDs
	gameResults := make([]tournament.GameResult, len(req.Games))
	for i, game := range req.Games {
		p1CharID := 0
		p2CharID := 0

		if game.P1Char != "" {
			if id, ok := characters.GetCharacterID(game.P1Char); ok {
				p1CharID = id
			}
		}

		if game.P2Char != "" {
			if id, ok := characters.GetCharacterID(game.P2Char); ok {
				p2CharID = id
			}
		}

		gameResults[i] = tournament.GameResult{
			Winner:   game.Winner,
			P1Char:   game.P1Char,
			P2Char:   game.P2Char,
			P1CharID: p1CharID,
			P2CharID: p2CharID,
		}
	}

	// Validate the set
	if err := tournament.ValidateSetScore(gameResults); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the set details to report
	tournamentSlug := c.DefaultQuery("tournament", "octagon")
	redemption := c.DefaultQuery("redemption", "false") == "true"

	fullTournamentSlug, err := startgg.GetTournamentSlug(context.Background(), tournamentSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	event := startgg.EventUltimateSingles
	if redemption {
		event = startgg.EventRedemptionBracket
	}

	eventSlug := fmt.Sprintf(startgg.EventSlugFormat, fullTournamentSlug, event)
	sets, err := tournament.FetchReportableSets(context.Background(), eventSlug, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Find the set
	var targetSet *tournament.Set
	for _, set := range sets {
		if set.ID == req.SetID {
			targetSet = &set
			break
		}
	}

	if targetSet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Set not found"})
		return
	}

	// Report the set
	err = tournament.ReportSet(context.Background(), *targetSet, gameResults)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Cache the character choices for each player
	for _, game := range req.Games {
		if game.P1Char != "" {
			if err := cache.SetPlayerCharacter(targetSet.Player1.ID, game.P1Char); err != nil {
				fmt.Printf("Failed to cache P1 character: %v\n", err)
			}
		}
		if game.P2Char != "" {
			if err := cache.SetPlayerCharacter(targetSet.Player2.ID, game.P2Char); err != nil {
				fmt.Printf("Failed to cache P2 character: %v\n", err)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Set reported successfully",
	})
}

func getCharacters(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=3600")
	c.JSON(http.StatusOK, gin.H{
		"characters": characters.Characters,
	})
}
