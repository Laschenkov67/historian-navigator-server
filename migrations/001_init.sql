CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY,
                                     email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('teacher','student')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE TABLE IF NOT EXISTS programs (
                                        id UUID PRIMARY KEY,
                                        author_id UUID REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    category VARCHAR(100),
    description TEXT,
    target_block TEXT,
    content_block TEXT,
    result_block TEXT,
    duration_hours INT DEFAULT 34,
    grade_level VARCHAR(50),
    is_template BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
    );

CREATE TABLE IF NOT EXISTS tests (
                                     id UUID PRIMARY KEY,
                                     title VARCHAR(500) NOT NULL,
    topic VARCHAR(200) NOT NULL,
    description TEXT,
    questions JSONB NOT NULL,
    time_limit INT DEFAULT 20
    );

CREATE TABLE IF NOT EXISTS test_results (
                                            id UUID PRIMARY KEY,
                                            test_id UUID REFERENCES tests(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    score INT,
    max_score INT,
    percentage NUMERIC(5,2),
    correct INT,
    incorrect INT,
    time_taken_sec INT,
    answers JSONB,
    completed_at TIMESTAMPTZ DEFAULT now()
    );
CREATE INDEX idx_results_user ON test_results(user_id);
CREATE INDEX idx_results_test ON test_results(test_id);