package main

import (
	"alexedwards.net/snippetbox/pkg/models"
	"errors"
	"fmt"
	_ "github.com/gorilla/mux"
	"net/http"
	"strconv"
	_ "strings"
	_ "unicode/utf8"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	s, err := app.news.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "home.page.tmpl", &templateData{
		NewsArray: s,
	})
}

func (app *application) showNews(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	n, err := app.news.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		}
		app.serverError(w, err)
		return
	}
	comments, err := app.news.GetComment(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	userID := app.session.GetInt(r, "authenticatedUserID")
	app.render(w, r, "show.page.tmpl", &templateData{
		News:                n,
		CommentList:         comments,
		AuthenticatedUserID: userID,
	})
}
func (app *application) creationPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{})
}

func (app *application) createNews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title := r.FormValue("title")
	content := r.FormValue("content")
	category := r.FormValue("category")
	if title == "" || content == "" || category == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	if len(title) > 20 || len(content) < 10 || len(content) > 200 {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	validCategories := []string{"Students", "Staff", "Applicants", "Researches"}
	if !malika(validCategories, category) {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	id, err := app.news.Insert(title, content, category)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "News successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/news?id=%d", id), http.StatusSeeOther)
}
func malika(slice []string, s string) bool {
	for _, value := range slice {
		if value == s {
			return true
		}
	}
	return false
}
func (app *application) filterCategory(category string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := app.renderCategoryPage(w, r, category)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}
}
func (app *application) renderCategoryPage(w http.ResponseWriter, r *http.Request, category string) error {
	newsArray, err := app.news.GetByCategory(category)
	if err != nil {
		return err
	}
	app.render(w, r, "category.page.tmpl", &templateData{
		Category:  category,
		NewsArray: newsArray,
	})
	return nil
}

func (app *application) contacts(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, request, "contacts.page.tmpl", &templateData{})
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", nil)
}
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	fmt.Println(name + email + password)
	//role := models.RoleUser
	err = app.users.Insert(name, email, password)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Sign up is successful")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{})
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	ID, err := app.users.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			app.session.Put(r, "flash", "Invalid email or password.")
			app.render(w, r, "login.page.tmpl", &templateData{})
			return
		}
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "authenticatedUserID", ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) Admin(w http.ResponseWriter, r *http.Request) {
	userList, err := app.users.AdminGet()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "admin.page.tmpl", &templateData{UserArray: userList})
}
func (app *application) Role(w http.ResponseWriter, r *http.Request) {
	fmt.Println("before check")
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("after check")
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	userID := r.PostForm.Get("userID")
	newRole := r.PostForm.Get("newRole")
	fmt.Println(userID)
	fmt.Println(newRole)
	err = app.users.Role(userID, newRole)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, "/user/admin", http.StatusSeeOther)
}
func (app *application) AddComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	userId := app.session.GetInt(r, "authenticatedUserID")
	newsId, err := strconv.Atoi(r.FormValue("newsID"))
	if err != nil {
		app.serverError(w, err)
	}
	text := r.FormValue("text")
	err = app.comments.Insert(userId, newsId, text)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/news?id=%d", newsId), http.StatusSeeOther)
}
func (app *application) deleteComment(w http.ResponseWriter, r *http.Request) {
	userID := app.session.GetInt(r, "authenticatedUserID")
	user, err := app.users.Get(userID)

	commentID, err := strconv.Atoi(r.FormValue("commentID"))
	if err != nil || commentID < 1 {
		app.serverError(w, err)
		return
	}
	newsId, err := app.comments.GetNewsId(commentID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	authorId, err := app.comments.GetAuthorId(commentID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if user.Role != "admin" && userID != authorId {
		app.session.Put(r, "flash", "You can only delete your own comments!")
		http.Redirect(w, r, fmt.Sprintf("/news?id=%d", newsId), http.StatusSeeOther)
		return
	}
	err = app.comments.Delete(commentID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/news?id=%d", newsId), http.StatusSeeOther)
}
