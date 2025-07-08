package handlers_test

// import (
// 	"fmt"
// 	"net/http"
// 	"testing"

// 	"github.com/gavv/httpexpect/v2"
// )

// func TestUpdateUserProgress(t *testing.T) {
// 	setupTest(t)

// 	e := httpexpect.WithConfig(httpexpect.Config{
// 		BaseURL:  testServer.URL,
// 		Client:   testServer.Client(),
// 		Reporter: httpexpect.NewAssertReporter(t),
// 	})

// 	req := map[string]interface{}{
// 		"name":          "Prog user",
// 		"email":         "progress@test.com",
// 		"provider":      "local",
// 		"google_id":     "",
// 		"password_hash": "some_hash",
// 		"role":          "user",
// 		"is_active":     true,
// 	}

// 	userResp := e.POST("/users").
// 		WithJSON(req).
// 		Expect().
// 		Status(http.StatusCreated).
// 		JSON().Object()

// 	createdUserID := userResp.Value("id").String().Raw()

// 	authHeader := fmt.Sprintf("Bearer %s", createdUserID)

// 	words := []map[string]interface{}{
// 		{
// 			"word":           "apple",
// 			"cefr_level":     "A1",
// 			"part_of_speech": "noun",
// 			"translation":    "яблоко",
// 			"context":        "I ate an apple",
// 			"audio_url":      "http://example.com/audio1.mp3",
// 		},
// 		{
// 			"word":           "banana",
// 			"cefr_level":     "A2",
// 			"part_of_speech": "noun",
// 			"translation":    "банан",
// 			"context":        "He peeled a banana",
// 			"audio_url":      "http://example.com/audio2.mp3",
// 		},
// 	}

// 	for _, w := range words {
// 		e.POST("/words").
// 			WithHeader("Authorization", authHeader).
// 			WithJSON(w).
// 			Expect().
// 			Status(http.StatusCreated)
// 	}

// 	progress := []map[string]interface{}{
// 		{
// 			"word":             "apple",
// 			"learned_at":       "2024-01-15T10:30:00Z",
// 			"confidence_score": 80,
// 			"cnt_reviewed":     2,
// 		},
// 		{
// 			"word":             "banana",
// 			"learned_at":       "2024-01-16T11:00:00Z",
// 			"confidence_score": 90,
// 			"cnt_reviewed":     1,
// 		},
// 	}

// 	e.POST("/progress").
// 		WithHeader("Authorization", authHeader).
// 		WithJSON(progress).
// 		Expect().
// 		Status(http.StatusOK)
// }
