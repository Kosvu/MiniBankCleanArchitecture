package db

import (
	"context"
	"errors"
	"log/slog"
	domains "minibank/internal/domain/users"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewUserRepository(db *pgxpool.Pool, log *slog.Logger) *UserRepository {
	return &UserRepository{
		db:  db,
		log: log,
	}
}

func (r *UserRepository) Create(ctx context.Context, user domains.User) error {
	sqlQuery := `
		INSERT INTO users(id, full_name, balance)
		VALUES ($1, $2, $3);
	`
	_, err := r.db.Exec(
		ctx,
		sqlQuery,
		user.ID,
		user.FullName,
		user.Balance,
	)

	if err != nil {
		r.log.Error("failed to create user", "err", err, "full_name", user.FullName)
		return err
	}

	return nil
}

func (r *UserRepository) GetUser(ctx context.Context, id uuid.UUID) (domains.User, error) {
	sqlQuery := `
		SELECT id, full_name, balance
		FROM users
		WHERE id=$1;
	`

	var user domains.User
	err := r.db.QueryRow(ctx, sqlQuery, id).Scan(&user.ID, &user.FullName, &user.Balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Warn("user not found", "user_id", id)
			return domains.User{}, domains.ErrUserNotFound
		}

		r.log.Error("failed to get user by id", "err", err, "user_id", id)
		return domains.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context, limit, offset int) ([]domains.User, error) {
	sqlQuery := `
		SELECT id, full_name, balance
		FROM users
		ORDER BY full_name
		LIMIT $1 OFFSET $2
		`

	rows, err := r.db.Query(ctx, sqlQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]domains.User, 0)

	for rows.Next() {
		var user domains.User
		if err := rows.Scan(&user.ID, &user.FullName, &user.Balance); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, user domains.User) error {
	sqlQuery := `
		UPDATE users
		SET full_name = $1, balance = $2
		WHERE id = $3
  	`
	tag, err := r.db.Exec(ctx, sqlQuery, user.FullName, user.Balance, user.ID)
	if err != nil {
		r.log.Error("failed to update user", "err", err, "full_name", user.FullName)
		return err
	}

	//Количество строк, которое изменилось
	if tag.RowsAffected() == 0 {
		return domains.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sqlQuery := `
	DELETE FROM users
	WHERE id = $1
 	`

	tag, err := r.db.Exec(ctx, sqlQuery, id)
	if err != nil {
		r.log.Error("failed to delete user", "err", err, "user_id", id)
		return err
	}

	if tag.RowsAffected() == 0 {
		return domains.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) CreateTransactionRecord(ctx context.Context, tx domains.Transaction) error {
	sqlQuery := `
	INSERT INTO transactions(id, user_id, transaction_type, amount, created_at)
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, sqlQuery, tx.ID, tx.UserID, tx.TransactionType, tx.Amount, tx.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUserTransactions(ctx context.Context, userID uuid.UUID) ([]domains.Transaction, error) {
	sqlQuery := `
	SELECT id, user_id, transaction_type, amount, created_at
	FROM transactions
	WHERE user_id = $1
	ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, sqlQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transaction []domains.Transaction

	for rows.Next() {
		var tx domains.Transaction

		err := rows.Scan(&tx.ID, &tx.UserID, &tx.TransactionType, &tx.Amount, &tx.CreatedAt)
		if err != nil {
			return nil, err
		}

		transaction = append(transaction, tx)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transaction, nil

}
