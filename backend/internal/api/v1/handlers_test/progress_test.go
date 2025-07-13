package handlers_test

// TODO: Fix authentication for this test
/*
func TestUpdateUserProgress(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	// Create a user for testing
	user := models.User{
		ID:       uuid.New(),
		Email:    "progress@test.com",
		Provider: "local",
		Name:     "Progress User",
		Role:     "user",
		IsActive: true,
	}
	err := userRepo.Create(context.Background(), &user)
	assert.NoError(t, err)

	// Create words for testing
	words := []models.Word{
		{
			ID:           uuid.New(),
			Word:         "apple",
			CEFRLevel:    "a1",
			PartOfSpeech: "noun",
			Translation:  "яблоко",
			Context:      "I ate an apple",
			AudioURL:     "http://example.com/audio1.mp3",
		},
		{
			ID:           uuid.New(),
			Word:         "banana",
			CEFRLevel:    "a2",
			PartOfSpeech: "noun",
			Translation:  "банан",
			Context:      "He peeled a banana",
			AudioURL:     "http://example.com/audio2.mp3",
		},
	}

	for _, w := range words {
		err := wordRepo.Create(context.Background(), &w)
		assert.NoError(t, err)
	}

	progress := []map[string]interface{}{
		{
			"word":             "apple",
			"translation":      "яблоко",
			"learned_at":       "2024-01-15T10:30:00Z",
			"confidence_score": 80,
			"cnt_reviewed":     2,
		},
		{
			"word":             "banana",
			"translation":      "банан",
			"learned_at":       "2024-01-16T11:00:00Z",
			"confidence_score": 90,
			"cnt_reviewed":     1,
		},
	}

	e.POST("/progress").
		WithJSON(progress).
		Expect().
		Status(http.StatusOK)
}
*/
