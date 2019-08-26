package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/samueldaviddelacruz/go-job-board/controllers"
	"github.com/samueldaviddelacruz/go-job-board/email"
	"github.com/samueldaviddelacruz/go-job-board/middleware"
	"github.com/samueldaviddelacruz/go-job-board/models"
	"github.com/samueldaviddelacruz/go-job-board/rand"
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
		models.WithGallery(),
		models.WithImage(),
		models.WithOAuth(),
	)
	must(err)

	defer services.Close()
	//must(services.DestructiveReset())
	must(services.AutoMigrate())

	mgCfg := appCfg.Mailgun
	emailer := email.NewClient(
		email.WithSender("lenslocked-project-demo.net Support", "support@sandboxddba781be75b455ea3313563bb0b74b2.mailgun.org"),
		email.WithMailgun(mgCfg.Domain, mgCfg.APIKey, mgCfg.PublicAPIKEY),
	)

	r := mux.NewRouter()

	jobsC  := controllers.NewJobs()

	usersC := controllers.NewUsers(services.User, emailer)



	randBytes, err := rand.Bytes(32)
	must(err)
	csrfMw := csrf.Protect(randBytes, csrf.Secure(appCfg.IsProd()))

	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{
		User: userMw,
	}

	r.HandleFunc("/", jobsC.Search).Methods("GET")


	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/logout", requireUserMw.ApplyFn(usersC.Logout)).Methods("POST")

	r.Handle("/forgot", usersC.ForgotPwView).Methods("GET")
	r.HandleFunc("/forgot", usersC.InitiateReset).Methods("POST")

	r.HandleFunc("/reset", usersC.ResetPw).Methods("GET")
	r.HandleFunc("/reset", usersC.CompleteReset).Methods("POST")

	//assets
	assetsHandler := http.FileServer(http.Dir("./assets"))
	assetsHandler = http.StripPrefix("/assets/", assetsHandler)
	r.PathPrefix("/assets/").Handler(assetsHandler)
	// Image routes
	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

	fmt.Printf("Running on port :%d", appCfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", appCfg.Port), csrfMw(userMw.Apply(r)))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
