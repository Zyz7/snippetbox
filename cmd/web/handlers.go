package main

import (
	"errors"
	"example/internal/models"
	"example/internal/validator"
	"fmt"
	"net/http"
	"strconv"
)

// Represent the form data and validation errors for the form fields
type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userPassForm struct {
	Current             string `form:"current"`
	New                 string `form:"new"`
	Confirm             string `form:"confirm"`
	validator.Validator `form:"-"`
}

// Define a home handler function which writes a byte slice containing
// Change the signature of the home handler so it is defined as method against
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// Add to header to the response header map
	w.Header().Add("Server example", "Go")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Call the newTemplateData helper to get a templateData struct containing
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Use the new render helper
	app.render(w, r, http.StatusOK, "home.html", data)

	/*
		files := []string{
			"../../ui/html/base.html",
			"../../ui/html/partials/nav.html",
			"../../ui/html/pages/home.html",
		}

		ts, err := template.ParseFiles(files...)

		if err != nil {
			// Because the home handler is now a method against the application
			// struct it can access its fields including the structured logger
			//app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
			//http.Error(w, "Internal server error", http.StatusInternalServerError)
			app.serverError(w, r, err)
			return
		}

		data := templateData{
			Snippets: snippets,
		}

		// The Execute method on the template set to write the template content
		// as the response body
		err = ts.ExecuteTemplate(w, "base", data)

		if err != nil {
			//app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
			//http.Error(w, "Internal server error", http.StatusInternalServerError)
			app.serverError(w, r, err)
		}

		//for _, snippet := range snippets {
		//	fmt.Fprintf(w, "%+v\n", snippet)
		//}
	*/
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// If there is not matching key in the session data this will return the empty string
	//flash := app.sessionManager.PopString(r.Context(), "flash")

	// Call the newTemplateData helper to get a templateData struct containing
	data := app.newTemplateData(r)
	data.Snippet = snippet
	//data.Flash = flash

	app.render(w, r, http.StatusOK, "show.html", data)

	/*
		files := []string{
			"../../ui/html/base.html",
			"../../ui/html/partials/nav.html",
			"../../ui/html/pages/view.html",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		// Create an instance of a templateData struct holding the snippet data
		data := templateData{
			Snippet: snippet,
		}

		// Pass in the templateData struct when executing the template
		err = ts.ExecuteTemplate(w, "base", data)
		if err != nil {
			app.serverError(w, r, err)
		}

		// Write the snippet data as a plain text http response body
		//fmt.Fprintf(w, "%+v", snippet)
	*/
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Send a 201 status code
	//w.WriteHeader(201)
	//w.WriteHeader(http.StatusCreated)
	//w.Write([]byte("Save a new snippet"))

	// Declare a new empty instance of the snippetCreateForm struct
	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	/*
		title := "O snail"
		content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
		expires := 7

		id, err := app.snippets.Insert(title, content, expires)
	*/

	// Initialize a map to hold any validation errors for the form fields
	//fieldErrors := make(map[string]string)

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Use put to add a string value and the key "flash"
	app.sessionManager.Put(r.Context(), "flash", "Snippet succesfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {

	var form userSignupForm

	err := app.decodePostForm(r, &form)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characteres long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.html", data)

}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {

	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {

		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	// Use PopString to retrieve and remove a value from the session data in one step
	path := app.sessionManager.PopString(r.Context(), "redirectPathAfterLogin")
	if path != "" {
		http.Redirect(w, r, path, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {

	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You have been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Return true if the current request is from an authenticated user otherwise return false
func (app *application) isAuthenticated(r *http.Request) bool {

	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "about.html", data)
}

func (app *application) accountView(w http.ResponseWriter, r *http.Request) {

	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	user, err := app.users.Get(userID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	data := app.newTemplateData(r)
	data.User = user

	app.render(w, r, http.StatusOK, "account.html", data)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) accountPassUpdate(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	data.Form = userPassForm{}
	app.render(w, r, http.StatusOK, "password.html", data)
}

func (app *application) accountPassUpdatePost(w http.ResponseWriter, r *http.Request) {

	var form userPassForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Current), "current", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.New), "new", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.New, 8), "new", "This field must be at least 8 characteres long")
	form.CheckField(validator.NotBlank(form.Confirm), "confirm", "This field cannot be blank")
	form.CheckField(form.New == form.Confirm, "confirm", "Password do not match")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "password.html", data)
		return
	}

	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	err = app.users.PassUpdate(userID, form.Current, form.New)

	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {

			form.AddFieldError("current", "Current password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "password.html", data)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	// Use put to add a string value and the key "flash"
	app.sessionManager.Put(r.Context(), "flash", "Your password has been updated")

	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}

func (app *application) snippetUpdate(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet
	data.Form = snippetCreateForm{Title: snippet.Title, Content: snippet.Content, Expires: 365}
	app.render(w, r, http.StatusOK, "update.html", data)
}

func (app *application) snippetUpdatePost(w http.ResponseWriter, r *http.Request) {

	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Snippet = snippet
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "update.html", data)
		return
	}

	err = app.snippets.Update(form.Title, form.Content, form.Expires, id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Use put to add a string value and the key "flash"
	app.sessionManager.Put(r.Context(), "flash", "Snippet succesfully updated!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetDelete(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, r, http.StatusOK, "delete.html", data)
}

func (app *application) snippetDeletePost(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	err = app.snippets.Delete(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet succesfully deleted!")
	app.home(w, r)
}
