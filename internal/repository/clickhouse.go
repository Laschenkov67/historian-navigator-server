package repository

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// NewClickHouseConn returns CH connection.
func NewClickHouseConn(dsn string) (driver.Conn, error) {
	opts, err := clickhouse.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	opts.DialTimeout = 5 * time.Second
	conn, err := clickhouse.Open(opts)
	if err != nil {
		return nil, err
	}
	return conn, conn.Ping(context.Background())
}

// AnalyticsRepo stores analytics in ClickHouse.
type AnalyticsRepo struct {
	conn driver.Conn
}

// NewAnalyticsRepo constructor.
func NewAnalyticsRepo(c driver.Conn) *AnalyticsRepo { return &AnalyticsRepo{conn: c} }

// TestEvent logged for each completed test.
type TestEvent struct {
	EventID    string
	UserID     string
	TestID     string
	Topic      string
	Score      uint32
	MaxScore   uint32
	Percentage float64
	Correct    uint32
	Incorrect  uint32
	DurationMs uint64
	CreatedAt  time.Time
}

// InsertTestEvent stores test event.
func (r *AnalyticsRepo) InsertTestEvent(ctx context.Context, e TestEvent) error {
	return r.conn.Exec(ctx,
		`INSERT INTO test_events(event_id,user_id,test_id,topic,score,max_score,percentage,correct,incorrect,duration_ms,created_at)
		 VALUES(?,?,?,?,?,?,?,?,?,?,?)`,
		e.EventID, e.UserID, e.TestID, e.Topic, e.Score, e.MaxScore, e.Percentage,
		e.Correct, e.Incorrect, e.DurationMs, e.CreatedAt)
}

// TopicStats aggregated stats by topic.
type TopicStats struct {
	Topic        string  `json:"topic"`
	Attempts     uint64  `json:"attempts"`
	AvgPercent   float64 `json:"avg_percent"`
	AvgCorrect   float64 `json:"avg_correct"`
	AvgIncorrect float64 `json:"avg_incorrect"`
}

// GetUserStats returns user statistics.
func (r *AnalyticsRepo) GetUserStats(ctx context.Context, userID string) ([]TopicStats, error) {
	rows, err := r.conn.Query(ctx,
		`SELECT topic, count() AS attempts, avg(percentage), avg(correct), avg(incorrect)
		 FROM test_events WHERE user_id=? GROUP BY topic`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []TopicStats
	for rows.Next() {
		var s TopicStats
		if err := rows.Scan(&s.Topic, &s.Attempts, &s.AvgPercent, &s.AvgCorrect, &s.AvgIncorrect); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

// GetTestAnalytics aggregates data by test.
func (r *AnalyticsRepo) GetTestAnalytics(ctx context.Context, testID string) (map[string]any, error) {
	var attempts uint64
	var avgPct, avgTime float64
	err := r.conn.QueryRow(ctx,
		`SELECT count(), avg(percentage), avg(duration_ms) FROM test_events WHERE test_id=?`, testID).
		Scan(&attempts, &avgPct, &avgTime)
	if err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(ctx,
		`SELECT toDate(created_at) d, avg(percentage) FROM test_events
		 WHERE test_id=? GROUP BY d ORDER BY d`, testID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type Point struct {
		Date time.Time `json:"date"`
		Avg  float64   `json:"avg"`
	}
	var series []Point
	for rows.Next() {
		var p Point
		if err := rows.Scan(&p.Date, &p.Avg); err != nil {
			return nil, err
		}
		series = append(series, p)
	}

	return map[string]any{
		"attempts":    attempts,
		"avg_percent": avgPct,
		"avg_time_ms": avgTime,
		"time_series": series,
	}, nil
}
