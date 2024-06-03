package main

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*User) error
	UpdateUser(int, *User) error
	DeleteUser(int) (*User, error)
	GetUserByID(int) (*User, error)
	GetAllUsers() ([]*User, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStoreage() (*PostgresStore, error) {
	godotenv.Load()
	connStr := os.Getenv("PG_CONN_STR")

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateUserTable()
}

func (s *PostgresStore) CreateUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS UserList (
		ID serial primary key,
		UUID char(36),
		username varchar(50),
		current_game varchar(50),
		current_level serial,
		created_at timestamp
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateUser(user *User) error {
	query := `INSERT INTO UserList (
		UUID,
		username,
		current_game,
		current_level,
		created_at
	)
	values (
		$1,
		$2,
		$3,
		$4,
		$5
	)`

	_, err := s.db.Exec(query, user.UUID, user.Username, user.CurrentGame, user.CurrentLevel, user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteUser(id int) (*User, error) {
	user := new(User)

	if row, err := s.db.Query("select * from userlist where ID = $1", id); err != nil {
		return nil, err
	} else {
		for row.Next() {
			if err := row.Scan(
				&user.ID,
				&user.UUID,
				&user.Username,
				&user.CurrentGame,
				&user.CurrentLevel,
				&user.CreatedAt); err != nil {
				return nil, err
			}
		}
	}

	_, err := s.db.Query("delete from userlist where ID = $1", id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresStore) UpdateUser(id int, user *User) error {
	query := `update userlist 
		SET
		UUID = $1,
		username = $2,
		current_game = $3,
		current_level = $4,
		created_at = $5
		where ID = $6`
	_, err := s.db.Exec(query, user.UUID, user.Username, user.CurrentGame, user.CurrentLevel, user.CreatedAt, id)
	if err != nil {
		return err
	}

	return nil
}
func (s *PostgresStore) GetUserByID(id int) (*User, error) {
	rows, err := s.db.Query("select * from userlist where ID = $1", id)
	if err != nil {
		return nil, err
	}

	user := new(User)
	for rows.Next() {
		if err := rows.Scan(
			&user.ID,
			&user.UUID,
			&user.Username,
			&user.CurrentGame,
			&user.CurrentLevel,
			&user.CreatedAt); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (s *PostgresStore) GetAllUsers() ([]*User, error) {
	rows, err := s.db.Query("select * from userlist")
	if err != nil {
		return nil, err
	}

	users := []*User{}
	for rows.Next() {
		user := new(User)
		if err := rows.Scan(
			&user.ID,
			&user.UUID,
			&user.Username,
			&user.CurrentGame,
			&user.CurrentLevel,
			&user.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
