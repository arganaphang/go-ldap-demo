package db

import (
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
}

type UserDB struct {
	db *sql.DB
}

const create string = `
CREATE TABLE IF NOT EXISTS users (
  id VARCHAR NOT NULL PRIMARY KEY,
  username VARCHAR NOT NULL UNIQUE,
  name VARCHAR NOT NULL,
  email VARCHAR NOT NULL UNIQUE
);`

func NewDatabase() (*UserDB, error) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(create); err != nil {
		return nil, err
	}
	return &UserDB{
		db: db,
	}, nil
}

func Add(u *UserDB, username, name, email string) error {
	_, err := u.db.Exec("INSERT INTO users VALUES(?,?,?,?);", uuid.New().String(), username, name, email)
	if err != nil {
		return err
	}
	return nil
}

func Get(u *UserDB, username string) (*User, error) {
	row := u.db.QueryRow("SELECT id, username, name, email FROM users WHERE UPPER(username)=UPPER(?)", username)
	user := User{}
	if err := row.Scan(&user.ID, &user.Username, &user.Name, &user.Email); err != nil {
		logrus.Info(err.Error())
		return nil, err
	}
	return &user, nil
}

func GetAll(u *UserDB) ([]User, error) {
	rows, err := u.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	data := []User{}
	for rows.Next() {
		i := User{}
		err = rows.Scan(&i.ID, &i.Username, &i.Name, &i.Email)
		if err != nil {
			return nil, err
		}
		data = append(data, i)
	}
	return data, nil
}
