package controllers

import (
	"github.com/samueldaviddelacruz/go-job-board/views"
	"net/http"
)

func NewJobs() *Jobs {
	return &Jobs{
		SearchView:    views.NewView("bootstrap", "jobs/search"),
	}
}
type JobPost struct {
	Title string
	Location string
	Category string
	Description string
	ApplicationEmail string
	Company string
}



// GET /
func (j *Jobs) Search(w http.ResponseWriter, r *http.Request) {
	mockJobs := []JobPost{
		{
			Title:"React Developer",
			Location:"Remote",
			Category:"Web Development",
			Company:"Netflix",
		},
		{
			Title:"Android Developer",
			Location:"Germany",
			Category:"Mobile",
			Company:"Apple",
		},
		{
			Title:"Postgres DBA",
			Location:"Ontario",
			Category:"DBA/Devops",
			Company:"Microsoft",
		},
		{
			Title:"Senior Automation QA",
			Location:"NY",
			Category:"QA",
			Company:"Google",
		},
	}
	var vd views.Data
	vd.Yield = mockJobs
	j.SearchView.Render(w, r, vd)
}

type Jobs struct {
	SearchView    *views.View
}
