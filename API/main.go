package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/samueldaviddelacruz/go-job-board/API/controllers"
	"github.com/samueldaviddelacruz/go-job-board/API/email"

	"github.com/samueldaviddelacruz/go-job-board/API/middleware"
	"github.com/samueldaviddelacruz/go-job-board/API/models"
)

func main() {

	boolPtr := flag.Bool("prod", false,
		"Provide this flag in production. This ensures that a config.json file is provided before the application starts")
	flag.Parse()
	appCfg := LoadConfig(*boolPtr)
	postgresConfig := appCfg.Database

	services, err := models.NewServices(
		models.WithGorm(
			postgresConfig.Dialect(),
			postgresConfig.ConnectionInfo()),
		models.WithLogMode(!appCfg.IsProd()),
		models.WithUser(appCfg.Pepper, appCfg.HMACKey),
		models.WithJobPost(),
		models.WithOAuth(),
	)
	must(err)

	defer services.Close()
	must(services.DestructiveReset())
	//must(services.AutoMigrate())

	mgCfg := appCfg.Mailgun
	emailer := email.NewClient(
		email.WithSender("lenslocked-project-demo.net Support", "support@sandboxddba781be75b455ea3313563bb0b74b2.mailgun.org"),
		email.WithMailgun(mgCfg.Domain, mgCfg.APIKey),
	)

	r := mux.NewRouter()

	jobsC := controllers.NewJobs(services.JobPost)

	companyC := controllers.NewCompany(services.User, emailer)

	//randBytes, err := rand.Bytes(32)
	must(err)
	//csrfMw := csrf.Protect(randBytes, csrf.Secure(appCfg.IsProd()))

	userMw := middleware.Company{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{
		Company: userMw,
	}

	r.HandleFunc("/jobs", jobsC.List).Methods("GET")
	r.HandleFunc("/jobs", jobsC.Create).Methods("POST")

	r.HandleFunc("/company/signup", companyC.Create).Methods("POST")

	r.HandleFunc("/company/login", companyC.Login).Methods("POST")
	r.HandleFunc("/company/logout", requireUserMw.ApplyFn(companyC.Logout)).Methods("POST")

	r.HandleFunc("/company/forgot", companyC.InitiateReset).Methods("POST")

	r.HandleFunc("/company/reset", companyC.CompleteReset).Methods("POST")

	fmt.Printf("Running on port :%d", appCfg.Port)
	must(http.ListenAndServe(fmt.Sprintf(":%d", appCfg.Port), userMw.Apply(r)))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
