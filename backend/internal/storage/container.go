package storage

import (
	"backend/internal/models"
	"context"
	"fmt"
)

func (s *PostgresqlDB) CreateStatus(status *models.ContainerStatus) error {
	conn, err := s.db.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	defer conn.Release()

	var existingID int
	checkSQL := `SELECT id FROM statuses WHERE ip = $1`
	err = conn.QueryRow(context.Background(), checkSQL, status.IP).Scan(&existingID)

	if err == nil {
		updateSQL := `UPDATE statuses SET ping_time = $1, last_check = $2 WHERE id = $3`
		_, err = conn.Exec(context.Background(), updateSQL, status.PingTime, status.LastCheck, existingID)
		if err != nil {
			return fmt.Errorf("failed to update status: %w", err)
		}
	} else {
		insertSQL := `INSERT INTO statuses (ip, ping_time, last_check) VALUES ($1, $2, $3) RETURNING id`
		err = conn.QueryRow(context.Background(), insertSQL, status.IP, status.PingTime, status.LastCheck).Scan(&status.ID)
		if err != nil {
			return fmt.Errorf("failed to create status: %w", err)
		}
	}

	return nil
}

func (s *PostgresqlDB) ReadAllStatuses() ([]*models.ContainerStatus, error) {
	conn, err := s.db.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	var statuses []*models.ContainerStatus
	rows, err := conn.Query(context.Background(), `SELECT id, ip, ping_time, last_check FROM statuses`)
	if err != nil {
		return nil, fmt.Errorf("failed to query all statuses: %w", err)
	}

	for rows.Next() {
		expr := &models.ContainerStatus{}
		err := rows.Scan(&expr.ID, &expr.IP, &expr.PingTime, &expr.LastCheck)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row into status: %w", err)
		}
		statuses = append(statuses, expr)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration over rows: %w", err)
	}

	return statuses, nil
}
