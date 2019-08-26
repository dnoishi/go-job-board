package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	llctx "github.com/samueldaviddelacruz/lenslocked.com/context"
	"github.com/samueldaviddelacruz/lenslocked.com/dbx"

	"github.com/gorilla/csrf"
	"github.com/samueldaviddelacruz/lenslocked.com/models"
	"golang.org/x/oauth2"
)

func NewAuths(os models.OAuthService, configs map[string]*oauth2.Config) *Oauths {
	return &Oauths{
		os: os,

		configs: configs,
	}
}

// Users Represents a Users controller
type Oauths struct {
	os      models.OAuthService
	configs map[string]*oauth2.Config
}

func (o *Oauths) Connect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	oauthConfig, ok := o.configs[service]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", http.StatusBadRequest)
		return
	}
	state := csrf.Token(r)
	cookie := http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	url := oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (o *Oauths) Callback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	oauthConfig, ok := o.configs[service]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", http.StatusBadRequest)
		return
	}

	r.ParseForm()
	state := r.FormValue("state")
	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if cookie == nil || cookie.Value != state {
		http.Error(w, "Invalid state provided", http.StatusBadRequest)
		return
	}

	cookie.Value = ""
	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)
	code := r.FormValue("code")

	token, err := oauthConfig.Exchange(context.TODO(), code)

	user := llctx.User(r.Context())
	existing, err := o.os.Find(user.ID, service)
	if err == models.ErrNotFound {

	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		o.os.Delete(existing.ID)
	}

	userOAuth := models.OAuth{
		UserID:  user.ID,
		Token:   *token,
		Service: service,
	}
	err = o.os.Create(&userOAuth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%+v", token)
}

func (o *Oauths) DropboxTest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]

	r.ParseForm()
	path := r.FormValue("path")

	user := llctx.User(r.Context())
	userOath, err := o.os.Find(user.ID, service)
	if err != nil {
		panic(err)
	}
	token := userOath.Token
	folders, files, err := dbx.List(token.AccessToken, path)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "Folders=", folders)
	fmt.Fprintln(w, "Files=", files)

}
