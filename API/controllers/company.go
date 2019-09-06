package controllers

import (
	"net/http"

	"github.com/samueldaviddelacruz/go-job-board/API/email"

	"github.com/samueldaviddelacruz/go-job-board/API/models"
)

// NewCompany is used to create a new Company controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup
func NewCompany(us models.UserService, emailer *email.Client) *Company {
	return &Company{
		us:      us,
		emailer: emailer,
	}
}

// Company Represents a Company controller
type Company struct {
	us      models.UserService
	emailer *email.Client
}
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Create is used to process the signup form when a user
// submits it. This is used to create a new user account.
//
// POST /signup
func (u *Company) Create(w http.ResponseWriter, r *http.Request) {

	credentials := Credentials{
	}
	//company.
	parseJSON(w, r, &credentials)
	companyUser := models.User{
		RoleID:   1,
		Password: credentials.Password,
		Email:    credentials.Email,
	}

	if err := u.us.Create(&companyUser); err != nil {
		//vd.SetAlert(err)
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	/*
		err := u.signIn(w, &user)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		u.emailer.Welcome(user.Name, user.Email)
	*/
	//views.RedirectAlert(w, r, "/galleries", http.StatusFound, alert)
	//user2, _ := u.us.ByID(1)

	respondJSON(w, http.StatusCreated, "resource created successfully")
}

// Login is used to verify the provided email address and
// password and then log the user in if they are correct.
//
// POST /login
func (u *Company) Login(w http.ResponseWriter, r *http.Request) {
	credentials := Credentials{
	}

	parseJSON(w, r, &credentials)
	companyUser := models.User{
		Password: credentials.Password,
		Email:    credentials.Email,
	}
	_, err := u.us.Authenticate(companyUser.Email, companyUser.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			respondJSON(w, http.StatusUnauthorized, "Invalid email address")
		default:
			//vd.SetAlert(err)
			respondJSON(w, http.StatusInternalServerError, err)
		}
		//u.LoginView.Render(w, r, vd)
		return
	}

	err = u.signIn(w, &companyUser)
	if err != nil {
		//u.LoginView.Render(w, r, vd)
		respondJSON(w, http.StatusInternalServerError, err)
		return
	}
	user2, _ := u.us.ByID(1)
	respondJSON(w, http.StatusOK, user2)
}

// Logout is used to delete a users session cookie (remember_token)
// and then will update the user resource with a new remember
// token.
// POST /logout
func (u *Company) Logout(w http.ResponseWriter, r *http.Request) {
	/*
		cookie := http.Cookie{
			Name:     "remember_token",
			Value:    "",
			Expires:  time.Now(),
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		user := context.User(r.Context())
		//token, _ := rand.RememberToken()

		u.us.Update(user)
		http.Redirect(w, r, "/", http.StatusFound)
	*/
}

// ResetPwForm is used to process the forgot password form
// and the reset password form
type ResetPwForm struct {
	Email    string `schema:"email"`
	Token    string `schema:"token"`
	Password string `schema:"password"`
}

// POST /forgot
func (u *Company) InitiateReset(w http.ResponseWriter, r *http.Request) {
	//var vd views.Data
	var form ResetPwForm
	//vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		//vd.SetAlert(err)
		//u.ForgotPwView.Render(w, r, vd)
		return
	}
	token, err := u.us.InitiateReset(form.Email)
	if err != nil {
		//vd.SetAlert(err)
		//u.ForgotPwView.Render(w, r, vd)
		return
	}

	err = u.emailer.ResetPw(form.Email, token)
	if err != nil {
		//vd.SetAlert(err)
		//u.ForgotPwView.Render(w, r, vd)
		return
	}
	/*
		views.RedirectAlert(w, r, "/reset", http.StatusFound, views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Instructions for resetting your password have been emailed to you.",
		})
	*/
}

// CompleteReset processes the reset password form
//
//POST
func (u *Company) CompleteReset(w http.ResponseWriter, r *http.Request) {
	//var vd views.Data
	var form ResetPwForm
	//vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		//vd.SetAlert(err)
		//u.ResetPwView.Render(w, r, vd)
		return
	}
	user, err := u.us.CompleteReset(form.Token, form.Password)
	if err != nil {
		//vd.SetAlert(err)
		//u.ResetPwView.Render(w, r, vd)
		return
	}

	err = u.signIn(w, user)
	if err != nil {
		//vd.SetAlert(err)
		//u.LoginView.Render(w, r, vd)
		return
	}
	/*
		views.RedirectAlert(w, r, "/galleries", http.StatusFound, views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Your password has been reset and you have been logged in!",
		})
	*/
}

// signIn is used to sign the given user via cookies.
func (u *Company) signIn(w http.ResponseWriter, user *models.User) error {
	/*
		if user.Remember == "" {
			token, err := rand.RememberToken()
			if err != nil {
				return err
			}
			user.Remember = token

			err = u.us.Update(user)
			if err != nil {
				return err
			}
		}

		cookie := http.Cookie{
			Name:     "remember_token",
			Value:    user.Remember,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
	*/
	return nil
}
