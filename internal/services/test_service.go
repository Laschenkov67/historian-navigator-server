package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/historian/backend/internal/models"
	"github.com/historian/backend/internal/repository"
)

// TestService handles tests.
type TestService struct {
	repo      *repository.TestRepo
	analytics *repository.AnalyticsRepo
	kafka     *repository.KafkaProducer
}

// NewTestService constructor.
func NewTestService(r *repository.TestRepo, a *repository.AnalyticsRepo, k *repository.KafkaProducer) *TestService {
	return &TestService{repo: r, analytics: a, kafka: k}
}

// List all tests.
func (s *TestService) List(ctx context.Context) ([]models.Test, error) {
	return s.repo.List(ctx)
}

// Get returns test (without correct answers).
func (s *TestService) Get(ctx context.Context, id string) (*models.Test, error) {
	t, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	for i := range t.Questions {
		t.Questions[i].Correct = 0
	}
	return t, nil
}

// SubmitRequest incoming data.
type SubmitRequest struct {
	Answers      []models.Answer `json:"answers"`
	TimeTakenSec int             `json:"time_taken_sec"`
}

// Submit checks answers and stores result.
func (s *TestService) Submit(ctx context.Context, testID, userID string, req SubmitRequest) (*models.TestResult, error) {
	test, err := s.repo.Get(ctx, testID)
	if err != nil {
		return nil, err
	}

	correctMap := make(map[string]int, len(test.Questions))
	pointsMap := make(map[string]int, len(test.Questions))
	maxScore := 0
	for _, q := range test.Questions {
		correctMap[q.ID] = q.Correct
		pointsMap[q.ID] = q.Points
		maxScore += q.Points
	}

	score, correct, incorrect := 0, 0, 0
	for i := range req.Answers {
		a := &req.Answers[i]
		if c, ok := correctMap[a.QuestionID]; ok && c == a.Selected {
			a.IsCorrect = true
			score += pointsMap[a.QuestionID]
			correct++
		} else {
			incorrect++
		}
	}

	pct := 0.0
	if maxScore > 0 {
		pct = float64(score) / float64(maxScore) * 100
	}

	result := &models.TestResult{
		ID:           uuid.NewString(),
		TestID:       testID,
		UserID:       userID,
		Score:        score,
		MaxScore:     maxScore,
		Percentage:   pct,
		Correct:      correct,
		Incorrect:    incorrect,
		TimeTakenSec: req.TimeTakenSec,
		Answers:      req.Answers,
		CompletedAt:  time.Now(),
	}

	if err := s.repo.SaveResult(ctx, result); err != nil {
		return nil, err
	}

	ev := repository.TestEvent{
		EventID:    uuid.NewString(),
		UserID:     userID,
		TestID:     testID,
		Topic:      test.Topic,
		Score:      uint32(score),
		MaxScore:   uint32(maxScore),
		Percentage: pct,
		Correct:    uint32(correct),
		Incorrect:  uint32(incorrect),
		DurationMs: uint64(req.TimeTakenSec * 1000),
		CreatedAt:  time.Now(),
	}
	_ = s.analytics.InsertTestEvent(ctx, ev)
	if b, err := json.Marshal(ev); err == nil {
		_ = s.kafka.Send(ctx, []byte(userID), b)
	}

	return result, nil
}

// MyResults returns user results.
func (s *TestService) MyResults(ctx context.Context, userID string) ([]models.TestResult, error) {
	return s.repo.ListResultsByUser(ctx, userID)
}
