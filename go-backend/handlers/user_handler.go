package handlers

import (
	"database/sql"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sedmo/integra-coding-assessment/go-backend/db"
	"github.com/sedmo/integra-coding-assessment/go-backend/models"
)

// @Summary Get users
// @Description Get all users
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User
// @Router /users [get]
func GetUsers(c echo.Context) error {
	query := db.Psql.Select("*").From("users")
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	rows, err := db.DB.Query(sqlQuery, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.UserStatus, &user.Department); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		users = append(users, user)
	}
	return c.JSON(http.StatusOK, users)
}

// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.User true "New User"
// @Success 201 {object} models.User
// @Router /users [post]
func CreateUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := user.Validate(); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.JSON(http.StatusBadRequest, validationErrors.Error())
	}

	// Check if username already exists
	existingUser, err := getUserByUsername(user.UserName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if existingUser != nil {
		return c.JSON(http.StatusConflict, "username already exists")
	}

	query := db.Psql.Insert("users").
		Columns("user_name", "first_name", "last_name", "email", "user_status", "department").
		Values(user.UserName, user.FirstName, user.LastName, user.Email, user.UserStatus, user.Department).
		Suffix("RETURNING user_id")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	err = db.DB.QueryRow(sqlQuery, args...).Scan(&user.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, user)
}

// @Summary Update user
// @Description Update an existing user
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param user body models.User true "Updated User"
// @Success 200 {object} models.User
// @Router /users/{id} [put]
func UpdateUser(c echo.Context) error {
	id := c.Param("id")
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := user.Validate(); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.JSON(http.StatusBadRequest, validationErrors.Error())
	}

	// Check if username already exists for a different user
	existingUser, err := getUserByUsername(user.UserName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if existingUser != nil && existingUser.UserID != user.UserID {
		return c.JSON(http.StatusConflict, "username already exists")
	}

	query := db.Psql.Update("users").
		Set("user_name", user.UserName).
		Set("first_name", user.FirstName).
		Set("last_name", user.LastName).
		Set("email", user.Email).
		Set("user_status", user.UserStatus).
		Set("department", user.Department).
		Where(squirrel.Eq{"user_id": id})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	_, err = db.DB.Exec(sqlQuery, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

// @Summary Delete user
// @Description Delete a user
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {string} string "User deleted"
// @Router /users/{id} [delete]
func DeleteUser(c echo.Context) error {
	id := c.Param("id")

	query := db.Psql.Delete("users").Where(squirrel.Eq{"user_id": id})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	_, err = db.DB.Exec(sqlQuery, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "User deleted")
}

func getUserByUsername(username string) (*models.User, error) {
	query := db.Psql.Select("*").From("users").Where(squirrel.Eq{"user_name": username})
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var user models.User
	err = db.DB.QueryRow(sqlQuery, args...).Scan(&user.UserID, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.UserStatus, &user.Department)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
