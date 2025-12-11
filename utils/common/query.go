package common

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/utils"
)

func AddLog(data *utils.Common, q *sqlx.Tx) error {
	_, err := q.Exec("INSERT INTO tbl_log (from_user_id, action, component) VALUES ($1, $2, $3)", data.FromUserID, data.Action, data.Component)
	return err
}

func ExecuteTransaction(ctx context.Context, db *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if pErr := recover(); pErr != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				fmt.Printf("failed to rollback transaction after panic: %v", rollbackErr)
			}
			panic(pErr)
		}

		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("transaction failed: %w, rollback failed: %v", err, rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", commitErr)
			}
		}
	}()

	err = fn(tx)
	return err
}
