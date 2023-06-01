package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/youssef-182/vue-skeleton-server/pkg/db"
	"github.com/youssef-182/vue-skeleton-server/pkg/models"
	"github.com/youssef-182/vue-skeleton-server/pkg/security"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

//type MyTime time.Time

//func (t *MyTime) UnmarshalJSON(data []byte) error {
//	if string(data) == "null" || string(data) == `""` {
//		return nil
//	}
//	return json.Unmarshal(data, (*time.Time)(t))
//}

func Register(w http.ResponseWriter, r *http.Request) {
	data := &models.UserRegisterRequest{}

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": fmt.Sprintf("An error has occurred while trying to parse the request. %v", err),
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	query := `INSERT INTO users (first_name, last_name, username, email_address, password, birthdate, gender) VALUES (?, ?, ?, ?,?, ?, ?)`
	tx, err := db.DB.Begin()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": fmt.Sprintf("An error has occurred while trying to begin a transaction. %v", err),
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": fmt.Sprintf("An error has occurred while trying to hash the user's password. %v", err),
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	result, err := tx.Exec(query, data.FirstName, data.LastName, data.Username, data.EmailAddress, string(hashedPassword), data.Birthdate, data.Gender)
	//result := db.DB.MustExec(query, user.FirstName, user.LastName, user.Username, user.EmailAddress, user.Password, user.Birthdate, user.Gender)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": fmt.Sprintf("An error has occurred while trying to store the user. %v", err),
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err = tx.Commit()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": fmt.Sprintf("An error has occurred while trying to commit the transaction. %v", err),
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	userID, err := result.LastInsertId()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": fmt.Sprintf("An error has occurred while trying to grab the last user ID. %v", err),
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	parsedBirthdate, err := time.Parse("2006-01-02", data.Birthdate)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": fmt.Sprintf("An error has occurred while trying to parse the birthdate field. %v", err),
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	user := models.User{
		Id:           userID,
		FirstName:    data.FirstName,
		LastName:     data.LastName,
		Username:     data.Username,
		EmailAddress: data.EmailAddress,
		Password:     string(hashedPassword),
		Birthdate:    parsedBirthdate,
		Gender:       data.Gender,
	}

	key := os.Getenv("JWT_SECRET_KEY")
	JwtStruct := security.Jwt{
		SecretKey: []byte(key),
		Claims: jwt.MapClaims{
			"id":       user.Id,
			"username": user.Username,
		},
		SigningMethod: jwt.SigningMethodHS256,
	}

	token, err := JwtStruct.GenerateToken()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": fmt.Sprintf("An error has occurred while trying to generate the token. %v", err),
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user":  user,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	data := &models.UserLoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": fmt.Sprintf("An error has occurred while trying to parse the request. %v", err),
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	var user models.User
	query := `SELECT * FROM users WHERE username=? LIMIT 1`
	err := db.DB.Get(&user, query, data.Username)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": fmt.Sprintf("Couldn't find user with the username: %s", data.Username),
			"errors": map[string][]string{
				"username": []string{"Username unavailable."},
				"password": []string{},
			},
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Wrong password",
			"errors": map[string][]string{
				"username": []string{},
				"password": []string{"Wrong password, please try again."},
			},
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	key := os.Getenv("JWT_SECRET_KEY")
	JwtStruct := security.Jwt{
		SecretKey: []byte(key),
		Claims: jwt.MapClaims{
			"id":       user.Id,
			"username": user.Username,
		},
		SigningMethod: jwt.SigningMethodHS256,
	}

	token, err := JwtStruct.GenerateToken()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": fmt.Sprintf("An error has occurred while trying to generate the token. %v", err),
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user":  user,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
