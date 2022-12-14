package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/LachlanStephan/ls_server/internal/models"
	"github.com/LachlanStephan/ls_server/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type blogCreateForm struct {
	Title   string
	Content string
	User_id int
	validator.Validator
}

func (app *application) blogView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	blog, err := app.blogs.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Blog = blog

	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) blogCreate(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// check that the user is logged in here
	// if not
	// redirect them to login page
	// this can be done once auth is implemented

	data := app.newTemplateData(r)
	data.Form = blogCreateForm{
		User_id: 1, // replace hardcoded id with id from session
	}

	app.render(w, http.StatusOK, "blog-create.tmpl.html", data)
}

func (app *application) blogCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := blogCreateForm{}
	// this needs to be retrieved from session
	// not from form
	form.User_id = validator.CastUserId(r.PostForm.Get("user_id"))
	form.Title = r.PostForm.Get("title")
	form.Content = r.PostForm.Get("content")

	validPublishers, err := app.users.GetAdminUsers()
	if err != nil {
		app.serverError(w, err)
		return
	}

	form.CheckField(validator.ValidUserId(form.User_id, validPublishers), "user_id", "You do not have permission to publish")
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 50), "title", "Title cannot have more than 50 chars")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Content, 50), "content", "Title cannot have more than 50 chars")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "blog-create.tmpl.html", data)
		return
	}

	id, err := app.blogs.Insert(form.User_id, form.Title, form.Content)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.SessionManager.Put(r.Context(), "flash", "Blog created successfully")

	http.Redirect(w, r, fmt.Sprintf("/blog/view/%d", id), http.StatusSeeOther)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	blogs, err := app.blogs.Recent()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.BlogLinks = blogs

	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) blog(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/blog" {
		app.notFound(w)
		return
	}

	blogs, err := app.blogs.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, http.StatusOK, "blog.tmpl.html", &templateData{
		BlogLinks: blogs,
	})
}

// func (app *application) admin(w http.ResponseWriter, r *http.Request) {
// if r.URL.Path != "/admin" {
// app.notFound(w)
// return
// }

// app.render(w, http.StatusOK, "admin.tmpl.html", nil)
// /*
// this will be an admin view
// just get the right templates/html and show it
// */
// }

// func (app *application) adminLogin(password string) bool {
// 	return false
/*
	hash password and see if matches the db for admin user
*/
// }
