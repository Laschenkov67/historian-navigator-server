package repository

import (
	"context"
	"encoding/json"

	"github.com/historian/backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TestRepo DB access for tests.
type TestRepo struct {
	pool *pgxpool.Pool
}

// NewTestRepo constructor.
func NewTestRepo(p *pgxpool.Pool) *TestRepo { return &TestRepo{pool: p} }

// List returns all tests.
func (r *TestRepo) List(ctx context.Context) ([]models.Test, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,title,topic,description,questions,time_limit FROM tests ORDER BY title`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Test
	for rows.Next() {
		var t models.Test
		var qb []byte
		if err := rows.Scan(&t.ID, &t.Title, &t.Topic, &t.Description, &qb, &t.TimeLimit); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(qb, &t.Questions)
		out = append(out, t)
	}
	return out, nil
}

// Get fetches test by id.
func (r *TestRepo) Get(ctx context.Context, id string) (*models.Test, error) {
	var t models.Test
	var qb []byte
	err := r.pool.QueryRow(ctx,
		`SELECT id,title,topic,description,questions,time_limit FROM tests WHERE id=$1`, id).
		Scan(&t.ID, &t.Title, &t.Topic, &t.Description, &qb, &t.TimeLimit)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(qb, &t.Questions)
	return &t, nil
}

// SaveResult stores result.
func (r *TestRepo) SaveResult(ctx context.Context, res *models.TestResult) error {
	ab, _ := json.Marshal(res.Answers)
	_, err := r.pool.Exec(ctx,
		`INSERT INTO test_results(id,test_id,user_id,score,max_score,percentage,correct,incorrect,time_taken_sec,answers,completed_at)
		 VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		res.ID, res.TestID, res.UserID, res.Score, res.MaxScore, res.Percentage,
		res.Correct, res.Incorrect, res.TimeTakenSec, ab, res.CompletedAt)
	return err
}

// ListResultsByUser returns user's results.
func (r *TestRepo) ListResultsByUser(ctx context.Context, userID string) ([]models.TestResult, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,test_id,user_id,score,max_score,percentage,correct,incorrect,time_taken_sec,answers,completed_at
		 FROM test_results WHERE user_id=$1 ORDER BY completed_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.TestResult
	for rows.Next() {
		var r models.TestResult
		var ab []byte
		if err := rows.Scan(&r.ID, &r.TestID, &r.UserID, &r.Score, &r.MaxScore, &r.Percentage,
			&r.Correct, &r.Incorrect, &r.TimeTakenSec, &ab, &r.CompletedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(ab, &r.Answers)
		out = append(out, r)
	}
	return out, nil
}
