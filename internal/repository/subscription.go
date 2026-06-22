package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/PhosFactum/effective-mobile-test-go/internal/models"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(sub *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(query, sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate)
	return err
}

func (r *SubscriptionRepository) GetByID(id uuid.UUID) (*models.Subscription, error) {
	var sub models.Subscription
	var endDate sql.NullTime

	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID,
		&sub.StartDate, &endDate, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if endDate.Valid {
		sub.EndDate = &endDate.Time
	}

	return &sub, nil
}

func (r *SubscriptionRepository) List(limit, offset int) ([]models.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		var endDate sql.NullTime
		err := rows.Scan(
			&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID,
			&sub.StartDate, &endDate, &sub.CreatedAt, &sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if endDate.Valid {
			sub.EndDate = &endDate.Time
		}
		subs = append(subs, sub)
	}

	return subs, nil
}

func (r *SubscriptionRepository) Update(sub *models.Subscription) error {
	query := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, start_date = $3, end_date = $4, updated_at = NOW()
		WHERE id = $5
	`
	_, err := r.db.Exec(query, sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate, sub.ID)
	return err
}

func (r *SubscriptionRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM subscriptions WHERE id = $1", id)
	return err
}

func (r *SubscriptionRepository) TotalCost(startDate, endDate time.Time, userID *uuid.UUID, serviceName *string) (int, error) {
	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE start_date >= $1 AND start_date <= $2
	`
	args := []interface{}{startDate, endDate}
	argIndex := 3

	if userID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *userID)
		argIndex++
	}

	if serviceName != nil && *serviceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", argIndex)
		args = append(args, *serviceName)
	}

	var total int
	err := r.db.QueryRow(query, args...).Scan(&total)
	return total, err
}
