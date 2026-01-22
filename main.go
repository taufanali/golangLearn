package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var category []Category

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(category)

		} else if r.Method == "POST" {
			var categoryBaru Category
			err := json.NewDecoder(r.Body).Decode(&categoryBaru)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			categoryBaru.ID = len(category) + 1
			category = append(category, categoryBaru)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(categoryBaru)
		} 
	})
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/api/categories/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		if r.Method == "GET" {
			for _, cat := range category {
				if cat.ID == id {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(cat)
					return
				}
			}
			http.Error(w, "Category not found", http.StatusNotFound)

		} else if r.Method == "DELETE" {
			for i, cat := range category {
				if cat.ID == id {
					category = append(category[:i], category[i+1:]...)
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}
			http.Error(w, "Category not found", http.StatusNotFound)
		} else if r.Method == "PUT" {
			var updatedCategory Category
			err := json.NewDecoder(r.Body).Decode(&updatedCategory)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			for i, cat := range category {
				if cat.ID == id {
					updatedCategory.ID = id
					category[i] = updatedCategory
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(updatedCategory)
					return
				}
			}
			http.Error(w, "Category not found", http.StatusNotFound)
		}
	})

	http.ListenAndServe(":8080", nil)
}

