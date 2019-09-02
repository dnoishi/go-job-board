package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/samueldaviddelacruz/go-job-board/API/models"
)

type Jobs struct {
}

func NewJobs() *Jobs {
	return &Jobs{}
}

// GET /jobs
func (j *Jobs) List(w http.ResponseWriter, r *http.Request) {
	mockJobs := []models.JobPost{
		{
			Title:    "React Developer",
			Location: "Remote",
			Category: "Web Development",
			//Company:  "Netflix",
		},
		{
			Title:    "Android Developer",
			Location: "Germany",
			Category: "Mobile",
			//Company:  "Apple",
		},
		{
			Title:    "Postgres DBA",
			Location: "Ontario",
			Category: "DBA/Devops",
			//Company:  "Microsoft",
		},
		{
			Title:    "Senior Automation QA",
			Location: "NY",
			Category: "QA",
			//Company:  "Google",
		},
	}

	respondJSON(w, http.StatusOK, mockJobs)
}

//POST /Create
func (j *Jobs) Create(w http.ResponseWriter, r *http.Request) {

	jobPost := models.JobPost{}
	err := json.NewDecoder(r.Body).Decode(&jobPost)
	if err != nil {
		respondJSON(w, 404, "Could not read new Book")
		return
	}
	//isbn, err := CreateBook(a, book)
	if err != nil {
		respondJSON(w, 404, "Could not create new Book")
		return
	}
	//respondJSON(w, 200, fmt.Sprintf("Created new Book with ISBN:%s", isbn))
}
