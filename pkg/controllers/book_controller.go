package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/youssef-182/vue-skeleton-server/pkg/db"
	"github.com/youssef-182/vue-skeleton-server/pkg/models"
	"net/http"
)

func IndexBook(w http.ResponseWriter, r *http.Request) {
	var books []models.BookResp

	query := `SELECT 
    	books.id,
    	books.title,
    	books.alternative_title,
    	books.description,
    	books.author,
    	books.release_date,
    	user.id as "user.id",
    	user.first_name as "user.first_name",
    	user.last_name as "user.last_name",
    	user.username as "user.username",
    	user.email_address as "user.email_address",
    	user.birthdate as "user.birthdate",
    	user.gender as "user.gender"
    FROM books JOIN users as user on user.id = books.user_id`
	if err := db.DB.Select(&books, query); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": fmt.Sprintf("An error has occurred while trying to fetch the books. %v", err),
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string][]models.BookResp{
		"books": books,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ShowBook(w http.ResponseWriter, r *http.Request) {
	bookCtx := r.Context().Value("book").(models.BookResp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]models.BookResp{
		"book": bookCtx,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func StoreBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(float64)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	data := &models.BookReq{}

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

	query := `INSERT INTO books (user_id, title, alternative_title, description, author, release_date) VALUES (?, ?, ?, ?, ?, ?)`
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

	result, err := tx.Exec(query, userID, data.Title, data.AlternativeTitle, data.Description, data.Author, data.ReleaseDate)
	//result := db.DB.MustExec(query, user.FirstName, user.LastName, user.Username, user.EmailAddress, user.Password, user.Birthdate, user.Gender)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": fmt.Sprintf("An error has occurred while trying to store the book. %v", err),
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

	bookID, err := result.LastInsertId()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": fmt.Sprintf("An error has occurred while trying to grab the last book ID. %v", err),
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	//parsedReleaseDate, err := time.Parse("2006-01-02", data.ReleaseDate)
	//if err != nil {
	//	w.Header().Set("Content-Type", "application/json")
	//	w.WriteHeader(http.StatusBadRequest)
	//	if err := json.NewEncoder(w).Encode(map[string]interface{}{
	//		"message": fmt.Sprintf("An error has occurred while trying to parse the release date field. %v", err),
	//	}); err != nil {
	//		w.WriteHeader(http.StatusInternalServerError)
	//		return
	//	}
	//	return
	//}

	var bookResp models.BookResp
	if err := db.DB.Get(&bookResp, `SELECT 
    		books.id,
    		books.title,
    		books.alternative_title,
    		books.description,
    		books.author,
    		books.release_date,
    		user.id as "user.id",
    		user.first_name as "user.first_name",
    		user.last_name as "user.last_name",
    		user.username as "user.username",
    		user.email_address as "user.email_address",
    		user.birthdate as "user.birthdate",
    		user.gender as "user.gender"
    	FROM books JOIN users AS user ON books.user_id=user.id WHERE books.id=?`, bookID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": fmt.Sprintf("An error has occurred while trying to fetch the book. %v", err),
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]models.BookResp{
		"book": bookResp,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	bookCtx := r.Context().Value("book").(models.BookResp)

	data := &models.BookReq{}
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

	query := `UPDATE books SET 
		title=?,
		alternative_title=?,
		description=?,
		author=?,
		release_date=?
	WHERE id=?`

	//if t, err := time.Parse("2006-01-02", bookCtx.ReleaseDate); err != nil {
	//	w.Header().Set("Content-Type", "application/json")
	//	w.WriteHeader(http.StatusInternalServerError)
	//	if err := json.NewEncoder(w).Encode(map[string]interface{}{
	//		"message": fmt.Sprintf("An error has occurred while trying to parse the release date. %v", err),
	//	}); err != nil {
	//		w.WriteHeader(http.StatusInternalServerError)
	//		return
	//	}
	//	return
	//} else {
	//	fmt.Println(t)
	//	db.DB.MustExec(query, bookCtx.Title, bookCtx.AlternativeTitle, bookCtx.Description, bookCtx.Author, t, bookCtx.Id)
	//	w.Header().Set("Content-Type", "application/json")
	//	w.WriteHeader(http.StatusNoContent)
	//}
	//if t, err := time.Parse("2006-01-02T03:04:05Z", bookCtx.ReleaseDate); err != nil {
	//	fmt.Println(t)
	//	w.Header().Set("Content-Type", "application/json")
	//	w.WriteHeader(http.StatusInternalServerError)
	//	if err := json.NewEncoder(w).Encode(map[string]interface{}{
	//		"message": fmt.Sprintf("An error has occurred while trying to parse the release date. %v", err),
	//	}); err != nil {
	//		w.WriteHeader(http.StatusInternalServerError)
	//		return
	//	}
	//	return
	//} else {
	db.DB.MustExec(query, data.Title, data.AlternativeTitle, data.Description, data.Author, data.ReleaseDate, bookCtx.Id)
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(http.StatusText(http.StatusNoContent)))
	//}
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	bookCtx := r.Context().Value("book").(models.BookResp)

	query := `DELETE FROM books WHERE id=?`

	db.DB.MustExec(query, bookCtx.Id)
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(http.StatusText(http.StatusNoContent)))
}
