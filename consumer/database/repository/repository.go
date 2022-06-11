package repository

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MentoringRequest struct {
	UUID             string    `gorm:"column:uuid"`
	TrackTitle       string    `json:"track_title"`
	ExerciseIconURL  string    `gorm:"column:exercise_icon_url"`
	ExerciseTitle    string    `gorm:"column:exercise_title"`
	StudentHandle    string    `gorm:"column:student_handle"`
	StudentAvatarUrl string    `gorm:"column:student_avatar_url"`
	Action           string    `gorm:"column:action"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
	AddedAt          time.Time `gorm:"column:added_at"`
}

type Repository struct {
	db *gorm.DB
}

func (r *Repository) SaveMentoringRequest(request MentoringRequest) error {
	return r.db.
		Table("mentoring_requests").
		Create(&request).Error
}

func New() (*Repository, error) {
	var err error
	var repo Repository

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("EXERCISM_EVENTS_DB_HOST"),
		os.Getenv("EXERCISM_EVENTS_DB_USER"),
		os.Getenv("EXERCISM_EVENTS_DB_PASSWORD"),
		os.Getenv("EXERCISM_EVENTS_DB_NAME"),
		os.Getenv("EXERCISM_EVENTS_DB_PORT"))

	repo.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("getting database session: %w", err)
	}

	return &repo, nil
}
