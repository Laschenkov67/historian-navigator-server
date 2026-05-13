package models

import "time"

// Role defines user role in the system.
type Role string

const (
	RoleTeacher Role = "teacher"
	RoleStudent Role = "student"
)

// User represents a system user.
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	FullName  string    `json:"full_name"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// Program — программа внеурочной деятельности по ФГОС.
type Program struct {
	ID           string    `json:"id"`
	AuthorID     string    `json:"author_id"`
	Title        string    `json:"title"`
	Category     string    `json:"category"` // kraevedenie, archaeology, debates, museum
	Description  string    `json:"description"`
	TargetBlock  string    `json:"target_block"`
	ContentBlock string    `json:"content_block"`
	ResultBlock  string    `json:"result_block"`
	Duration     int       `json:"duration_hours"`
	GradeLevel   string    `json:"grade_level"`
	IsTemplate   bool      `json:"is_template"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Test — тест по исторической теме.
type Test struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Topic       string     `json:"topic"`
	Description string     `json:"description"`
	Questions   []Question `json:"questions"`
	TimeLimit   int        `json:"time_limit_minutes"`
}

// Question — вопрос теста.
type Question struct {
	ID      string   `json:"id"`
	Text    string   `json:"text"`
	Options []string `json:"options"`
	Correct int      `json:"correct,omitempty"`
	Points  int      `json:"points"`
}

// TestResult — результат прохождения теста.
type TestResult struct {
	ID           string    `json:"id"`
	TestID       string    `json:"test_id"`
	UserID       string    `json:"user_id"`
	Score        int       `json:"score"`
	MaxScore     int       `json:"max_score"`
	Percentage   float64   `json:"percentage"`
	Correct      int       `json:"correct"`
	Incorrect    int       `json:"incorrect"`
	TimeTakenSec int       `json:"time_taken_sec"`
	Answers      []Answer  `json:"answers"`
	CompletedAt  time.Time `json:"completed_at"`
}

// Answer — ответ пользователя на вопрос.
type Answer struct {
	QuestionID string `json:"question_id"`
	Selected   int    `json:"selected"`
	IsCorrect  bool   `json:"is_correct"`
}
