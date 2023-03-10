package repositories

import (
	"api/src/models"
	"database/sql"
	"fmt"
)

type users struct {
	db *sql.DB
}

// Creates a new repository of users using the users struct as base
func NewUserRepo(db *sql.DB) *users {
	return &users{db}
}

// Creates a user insert into database
func (repo users) Create(user models.User) (uint64, error) {
	statement, err := repo.db.Prepare("INSERT INTO users (username, email, passwd) VALUES(?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(user.Username, user.Email, user.Password)
	if err != nil {
		return 0, err
	}

	lastIdInserted, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastIdInserted), nil
}

// Deletes a user by Id
func (repo users) Delete(Id uint64) error {
	statement, err := repo.db.Prepare("DELETE FROM users WHERE userId = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(Id); err != nil {
		return err
	}

	return nil
}

// Updates the info of a user
func (repo users) Update(Id uint64, user models.User) error {
	statement, err := repo.db.Prepare("UPDATE users SET username = ?, email = ? WHERE userId = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(user.Username, user.Email, Id); err != nil {
		return err
	}

	return nil
}

// Returns all the users that contains the string content
func (repo users) Search(username string) ([]models.User, error) {
	username = fmt.Sprintf("%%%s%%", username) //%username%
	lines, err := repo.db.Query(
		"SELECT userId, username, email, createdAt FROM users WHERE username LIKE ?", username,
	)
	if err != nil {
		return nil, err
	}
	defer lines.Close()

	var users []models.User
	for lines.Next() {
		var user models.User
		if err = lines.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

// Fetches by Id from the database
func (repo users) SearchById(Id uint64) (models.User, error) {
	lines, err := repo.db.Query(
		"SELECT userId, username, email, createdAt FROM users WHERE userId = ?", Id,
	)
	if err != nil {
		// has to return empty user in case of error
		return models.User{}, err
	}
	defer lines.Close()

	var user models.User

	if lines.Next() {
		if err = lines.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

func (repo users) SearchByEmail(email string) (models.User, error) {
	line, err := repo.db.Query("SELECT userId, passwd FROM users WHERE email = ?", email)
	if err != nil {
		return models.User{}, err
	}
	defer line.Close()

	var user models.User
	if line.Next() {
		if err = line.Scan(&user.ID, &user.Password); err != nil {
			return models.User{}, err
		}
	}

	return models.User{}, err
}

func (repo users) SearchPassword(userId uint64) (string, error) {
	lines, err := repo.db.Query("SELECT passwd FROM users WHERE userId = ?", userId)
	if err != nil {
		return "", err
	}
	defer lines.Close()

	var user models.User

	if lines.Next() {
		if err = lines.Scan(&user.Password); err != nil {
			return "", err
		}
	}

	return user.Password, nil
}

func (repo users) UpdatePassword(userId uint64, password string) error {
	statement, err := repo.db.Prepare("UPDATE users SET passwd = ? WHERE userId = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(password, userId); err != nil {
		return err
	}

	return nil
}
