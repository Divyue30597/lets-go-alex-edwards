package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Divyue30597/snippetbox-lets-go/internal/models"
	"github.com/Divyue30597/snippetbox-lets-go/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
	Title   string
	Content string
	Expires string
	validator.Validator
	// FieldErrors map[string]string
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

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {
	// 	// http.NotFound(w, r)
	// 	app.notFound(w)
	// 	return
	// }

	// panic("oops! something went wrong")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// for _, snippet := range snippets {
	// 	fmt.Fprintf(w, "%+v\n", snippet)
	// }

	// files := []string{
	// 	"./ui/html/pages/home.tmpl",
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// }
	// // templateSet, err := template.ParseFiles("./ui/html/pages/home.tmpl")

	// templateSet, err := template.ParseFiles(files...)
	// if err != nil {
	// 	// app.errorLog.Print(err.Error())
	// 	app.serverError(w, err)
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }

	// data := &templateData{
	// 	Snippets: snippets,
	// }

	// // err = templateSet.Execute(w, nil)
	// err = templateSet.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	// app.errorLog.Print(err.Error())
	// 	app.serverError(w, err)
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// }

	// replacing the above logic with render method on app
	// app.render(w, http.StatusOK, "home.tmpl", &templateData{
	// 	Snippets: snippets,
	// })

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)

	// w.Write([]byte("Hello from snippetbox."))
}

func (app *application) snippetViewHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	id, err := strconv.Atoi(params.ByName("id"))
	// r.URL.Query().Get("id")
	if err != nil || id < 1 {
		// http.NotFound(w, r)
		app.notFound(w)
		return
	}

	// Using the snippets to get a record
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/pages/view.tmpl",
	// }

	// templateSet, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// data := &templateData{
	// 	Snippet: snippet,
	// }

	// err = templateSet.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	app.serverError(w, err)
	// }

	// Changing the above operations to 2 different functions one for reading and another for rendering from the cache
	// app.render(w, http.StatusOK, "view.tmpl", &templateData{
	// 	Snippet: snippet,
	// })

	data := app.newTemplateData(r)
	data.Snippet = snippet

	// flash := app.sessionManager.PopString(r.Context(), "flash")
	// data.Flash = flash

	app.render(w, http.StatusOK, "view.tmpl", data)

	// w.Write([]byte("Display a specific snippet..."))
	// fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	// fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Can be removed since we are using httprouter to choose the request methods.
	// if r.Method != http.MethodPost {
	// 	w.Header().Set("Allow", http.MethodPost)
	// 	// w.Header()["Date"] = nil
	// 	// w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}
	// 	// http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Creating a form instance
	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: r.PostForm.Get("expires"),
		// FieldErrors: map[string]string{},
	}

	// title := r.PostForm.Get("title")
	// content := r.PostForm.Get("content")
	// expires := r.PostForm.Get("expires")

	// fieldErrors := make(map[string]string)

	// if strings.TrimSpace(form.Title) == "" {
	// 	form.FieldErrors["title"] = "This field cannot be blank"
	// } else if utf8.RuneCountInString(form.Title) > 100 {
	// 	form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	// }

	// if strings.TrimSpace(form.Content) == "" {
	// 	form.FieldErrors["content"] = ""
	// }

	// if form.Expires != "1 day" && form.Expires != "7 days" && form.Expires != "365 days" {
	// 	form.FieldErrors["expires"] = "This field must be equal to 1 day, 7 days, 365 days."
	// }

	// if len(form.FieldErrors) > 0 {
	// 	data := app.newTemplateData(r)
	// 	data.Form = form
	// 	app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
	// 	fmt.Fprint(w, form.FieldErrors)
	// 	return
	// }

	// Above logic is replace here with the validator package.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedString(form.Expires, "1 day", "7 days", "365 days"), "expires", "This field must be equal to 1 day or 7 days or 365 days")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		// fmt.Fprint(w, form.FieldErrors)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	// w.Write([]byte("Create a snippet..."))
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetCreateHandler(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("Display the form for creating a new snippet..."))
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: "365 days",
	}

	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) userSignUpHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl", data)
}

// Using formDecoder
func (app *application) userSignUpPostHandler(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	// err = app.formDecoder.Decode(&form, r.PostForm)
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using our helper functions.
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	// If there are any errors, redisplay the signup form along with a 422
	// status code.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your Sign up was successful. Please log in.")

	// fmt.Fprintln(w, "Create a new user...")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLoginHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "Display an HTML form for logging in the user...")
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl", data)
}

func (app *application) userLoginPostHandler(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	// fmt.Println(err)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Use the RenewToken() method on the current session to change the session
	// ID. It's good practice to generate a new session ID when the
	// authentication state or privilege levels changes for the user (e.g. login
	// and logout operations).
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticateUserID", id)

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)

	// fmt.Fprintln(w, "Authenticate and login a user...")
}

func (app *application) userLogoutPostHandler(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticateUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	// fmt.Fprintln(w, "Logout user...")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
