package main

import (
	"ass1/pkg/models/mysql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Data struct {
	Title string

	NewsList []mysql.News
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w) // Use the notFound() helper
		return
	}
	news, err := app.NewsModel.GetNews()
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Создаем структуру для передачи данных в шаблон.
	data := Data{
		Title:    "Home",
		NewsList: news,
	}

	files := []string{
		"./ui/html/pages/homepage.html",
		"./ui/html/partials/news_section.html",
		"./ui/html/base.html",
		"./ui/html/partials/navbar.html",
		"./ui/html/partials/footer.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) // Use the serverError() helper.
		return
	}
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err) // Use the serverError() helper.
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper.
		return
	}
	w.Write([]byte("Create a new snippet..."))
}

func (app *application) foradmins(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/foradmins" {
		app.notFound(w) // Use the notFound() helper
		return
	}
	// Создаем структуру для передачи данных в шаблон.

	files := []string{
		"./ui/html/pages/foradmins.html",
		"./ui/html/base.html",
		"./ui/html/partials/navbar.html",
		"./ui/html/partials/news_section.html",
		"./ui/html/partials/footer.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) // Use the serverError() helper.
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err) // Use the serverError() helper.
	}
}

// Handler for creating news
func (app *application) createNews(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/foradmins/createnews" {
		app.notFound(w) // Use the notFound() helper
		return
	}

	if r.Method == http.MethodPost {
		// Handle form submission.
		title := r.FormValue("inputTitle")
		content := r.FormValue("inputText")
		tag := r.FormValue("tag") // Assuming you add the 'name' attribute to the select element.
		imageURL := ""            // Add logic to handle image URL.

		// Data validation
		if title == "" || content == "" || tag == "Open this select menu" {
			// Invalid data, return an error or show a validation message
			// You can customize this based on your application's needs
			http.Error(w, "Invalid data. Please fill in all fields.", http.StatusBadRequest)
			return
		}

		// Assuming you have a function IsValidTag to validate the tag
		if !IsValidTag(tag) {
			http.Error(w, "Invalid tag selected.", http.StatusBadRequest)
			return
		}

		if title != "" {
			_, err := app.NewsModel.CreateNews(title, content, tag, imageURL)
			if err != nil {
				app.serverError(w, err)
				return
			}
			http.Redirect(w, r, "/success", http.StatusSeeOther)
			return
		}
	}

	files := []string{
		"./ui/html/pages/create_news.html",
		"./ui/html/base.html",
		"./ui/html/partials/navbar.html",
		"./ui/html/partials/footer.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

// IsValidTag checks if the provided tag is a valid option
func IsValidTag(tag string) bool {
	validTags := map[string]bool{
		"For Students":   true,
		"For Staff":      true,
		"For Applicants": true,
	}
	return validTags[tag]
}

// Handler for deleting news
func (app *application) deleteNews(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/foradmins/deletenews" {
		app.notFound(w) // Use the notFound() helper
		return
	}
	if r.Method == http.MethodPost {
		// Handle form submission.
		deletetitle := r.FormValue("deleteTitle")

		// Use the CreateNews method to save the form data to the database.
		if deletetitle == "" {
			// Invalid data, return an error or show a validation message
			// You can customize this based on your application's needs
			http.Error(w, "Invalid data. Please fill in all fields.", http.StatusBadRequest)
			return
		}
		err := app.NewsModel.DeleteNews(deletetitle)
		if err != nil {
			app.serverError(w, err)
			return
		}

		// Redirect or display a success message.
		http.Redirect(w, r, "/success", http.StatusSeeOther)
		return
	}
	files := []string{
		"./ui/html/pages/delete_news.html",
		"./ui/html/base.html",
		"./ui/html/partials/navbar.html",
		"./ui/html/partials/footer.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

// Handler for deleting news
func (app *application) updateNews(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/foradmins/updatenews" {
		app.notFound(w) // Use the notFound() helper
		return
	}
	if r.Method == http.MethodPost {
		// Handle form submission.
		old := r.FormValue("inputTitle")
		next := r.FormValue("inputNextTitle")
		nextcontent := r.FormValue("inputText")
		tag := r.FormValue("tag")

		if old == "" || next == "" || nextcontent == "" || tag == "Open this select menu" {
			// If any required field is empty, return a validation error.
			app.clientError(w, http.StatusBadRequest)
			return
		}
		// Use the CreateNews method to save the form data to the database.
		err := app.NewsModel.UpdateNews(old, next, nextcontent, tag, "")
		if err != nil {
			app.serverError(w, err)
			return
		}

		// Redirect or display a success message.
		http.Redirect(w, r, "/success", http.StatusSeeOther)
		return
	}
	files := []string{
		"./ui/html/pages/update_news.html",
		"./ui/html/base.html",
		"./ui/html/partials/navbar.html",
		"./ui/html/partials/footer.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) forStudents(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/forstudents" {
		app.notFound(w) // Use the notFound() helper
		return
	}
	news, err := app.NewsModel.GetNewsByCategory("For Students")
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Создаем структуру для передачи данных в шаблон.
	data := Data{
		Title:    "Home",
		NewsList: news,
	}

	files := []string{
		"./ui/html/pages/forstudents.html",
		"./ui/html/partials/news_section.html",
		"./ui/html/base.html",
		"./ui/html/partials/navbar.html",
		"./ui/html/partials/footer.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) // Use the serverError() helper.
		return
	}
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err) // Use the serverError() helper.
	}
}
func (app *application) forStaff(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/forstaff" {
		app.notFound(w) // Use the notFound() helper
		return
	}
	news, err := app.NewsModel.GetNewsByCategory("For Staff")
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Создаем структуру для передачи данных в шаблон.
	data := Data{
		Title:    "Home",
		NewsList: news,
	}

	files := []string{
		"./ui/html/pages/forstaff.html",
		"./ui/html/partials/news_section.html",
		"./ui/html/base.html",
		"./ui/html/partials/navbar.html",
		"./ui/html/partials/footer.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) // Use the serverError() helper.
		return
	}
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err) // Use the serverError() helper.
	}
}
func (app *application) forApplicants(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/forapplicants" {
		app.notFound(w) // Use the notFound() helper
		return
	}
	news, err := app.NewsModel.GetNewsByCategory("For Applicants")
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Создаем структуру для передачи данных в шаблон.
	data := Data{
		Title:    "Home",
		NewsList: news,
	}

	files := []string{
		"./ui/html/pages/forapplicants.html",
		"./ui/html/partials/news_section.html",
		"./ui/html/base.html",
		"./ui/html/partials/navbar.html",
		"./ui/html/partials/footer.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) // Use the serverError() helper.
		return
	}
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err) // Use the serverError() helper.
	}
}
func (app *application) success(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/success" {
		app.notFound(w) // Use the notFound() helper
		return
	}

	files := []string{
		"./ui/html/pages/success.html",
		"./ui/html/partials/news_section.html",
		"./ui/html/base.html",
		"./ui/html/partials/navbar.html",
		"./ui/html/partials/footer.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) // Use the serverError() helper.
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err) // Use the serverError() helper.
	}
}
