package db

import (
	"context"
	domains "minibank/internal/domain/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
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

	return err
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

	_, err := r.db.Exec(ctx, sqlQuery, id)
	if err != nil {
		return err
	}

	return nil
}
