package main

import (
    "net/http"
    "strings"
    "unicode/utf8"
)

type FormData struct {
	PageTitle        string
	CurrentPage      string
	FormAction       string
	ListURL          string
	TitleLabel       string
	LocationLabel    string
	CategoryLabel    string
	DescriptionLabel string
	SubmitLabel      string
	Categories       []string
	CategoryValid    bool
	Title            string
	Location         string
	Category         string
	Description      string
	TitleError       string
	LocationError    string
	CategoryError    string
	DescriptionError string
}

type ListData struct {
	PageTitle      string
	CurrentPage    string
	FormURL        string
	FormLabel      string
	OtherListURL   string
	OtherListLabel string
	EmptyMessage   string
	Items          []Submission
}

type formConfig struct {
	kind             string
	pageTitle        string
	currentPage      string
	formTemplate     string
	listTemplate     string
	formAction       string
	listURL          string
	titleLabel       string
	locationLabel    string
	categoryLabel    string
	descriptionLabel string
	submitLabel      string
	listTitle        string
	formLabel        string
	otherListURL     string
	otherListLabel   string
	emptyMessage     string
	categories       []string
}

var ticketConfig = formConfig{
	kind:             "ticket",
	pageTitle:        "Support ticket",
	currentPage:      "support-ticket",
	formTemplate:     "ticket-form.tmpl",
	listTemplate:     "tickets.tmpl",
	formAction:       "/tickets/create",
	listURL:          "/tickets",
	titleLabel:       "Ticket title",
	locationLabel:    "Room or location",
	categoryLabel:    "Support category",
	descriptionLabel: "Description",
	submitLabel:      "Create ticket",
	listTitle:        "Support tickets",
	formLabel:        "New ticket",
	otherListURL:     "/software-requests",
	otherListLabel:   "Software requests",
	emptyMessage:     "No support tickets have been submitted yet.",
	categories:       []string{"Network", "Email", "Hardware", "Other"},
}

var softwareRequestConfig = formConfig{
	kind:             "software",
	pageTitle:        "Software request",
	currentPage:      "software-request",
	formTemplate:     "software-form.tmpl",
	listTemplate:     "software-requests.tmpl",
	formAction:       "/software-requests/create",
	listURL:          "/software-requests",
	titleLabel:       "Software name",
	locationLabel:    "Room or location",
	categoryLabel:    "Request type",
	descriptionLabel: "Description",
	submitLabel:      "Submit software request",
	listTitle:        "Software requests",
	formLabel:        "New software request",
	otherListURL:     "/tickets",
	otherListLabel:   "Support tickets",
	emptyMessage:     "No software requests have been submitted yet.",
	categories:       []string{"Installation", "Upgrade", "Licence", "Other"},
}

func ticketCreate(w http.ResponseWriter, r *http.Request) {
	renderForm(w, http.StatusOK, ticketConfig, FormData{})
}

func ticketCreatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	form := FormData{
		Title:       r.PostForm.Get("title"),
		Location:    r.PostForm.Get("location"),
		Category:    r.PostForm.Get("category"),
		Description: r.PostForm.Get("description"),
	}
	title := strings.TrimSpace(form.Title)
	location := strings.TrimSpace(form.Location)
	description := strings.TrimSpace(form.Description)

	if title == "" {
		form.TitleError = "Enter an issue title."
	} else if utf8.RuneCountInString(title) > 80 {
		form.TitleError = "Use at most 80 code points."
	}
	if location == "" {
		form.LocationError = "Enter a room or location."
	} else if utf8.RuneCountInString(location) > 40 {
		form.LocationError = "Use at most 40 code points."
	}
	if !categoryIsValid(form.Category, ticketConfig.categories) {
		form.CategoryError = "Choose a listed support category."
	}
	if description == "" {
		form.DescriptionError = "Describe the issue."
	} else if utf8.RuneCountInString(description) > 500 {
		form.DescriptionError = "Use at most 500 code points."
	}

	if form.TitleError != "" || form.LocationError != "" ||
		form.CategoryError != "" || form.DescriptionError != "" {
		renderForm(w, http.StatusUnprocessableEntity, ticketConfig, form)
		return
	}

	addSubmission(Submission{
		Kind:        "ticket",
		Title:       title,
		Location:    location,
		Category:    form.Category,
		Description: description,
	})
	http.Redirect(w, r, "/tickets", http.StatusSeeOther)
	return
}

func listTickets(w http.ResponseWriter, r *http.Request) {
	renderList(w, ticketConfig)
}

func softwareCreate(w http.ResponseWriter, r *http.Request) {
	renderForm(w, http.StatusOK, softwareRequestConfig, FormData{})
}

func softwareCreatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	form := FormData{
		Title:       r.PostForm.Get("title"),
		Location:    r.PostForm.Get("location"),
		Category:    r.PostForm.Get("category"),
		Description: r.PostForm.Get("description"),
	}
	title := strings.TrimSpace(form.Title)
	location := strings.TrimSpace(form.Location)
	description := strings.TrimSpace(form.Description)

	if title == "" {
		form.TitleError = "Enter a software name."
	} else if utf8.RuneCountInString(title) > 80 {
		form.TitleError = "Use at most 80 code points."
	}
	if location == "" {
		form.LocationError = "Enter a room or location."
	} else if utf8.RuneCountInString(location) > 40 {
		form.LocationError = "Use at most 40 code points."
	}
	if !categoryIsValid(form.Category, softwareRequestConfig.categories) {
		form.CategoryError = "Choose a listed request type."
	}
	if description == "" {
		form.DescriptionError = "Enter a description."
	} else if utf8.RuneCountInString(description) > 500 {
		form.DescriptionError = "Use at most 500 code points."
	}

	if form.TitleError != "" || form.LocationError != "" ||
		form.CategoryError != "" || form.DescriptionError != "" {
		renderForm(w, http.StatusUnprocessableEntity, softwareRequestConfig, form)
		return
	}

	addSubmission(Submission{
		Kind:        "software",
		Title:       title,
		Location:    location,
		Category:    form.Category,
		Description: description,
	})
	http.Redirect(w, r, "/software-requests", http.StatusSeeOther)
	return
}

func listSoftwareRequests(w http.ResponseWriter, r *http.Request) {
	renderList(w, softwareRequestConfig)
}

func renderForm(w http.ResponseWriter, status int, config formConfig, form FormData) {
	form.PageTitle = config.pageTitle
	form.CurrentPage = config.currentPage
	form.FormAction = config.formAction
	form.ListURL = config.listURL
	form.TitleLabel = config.titleLabel
	form.LocationLabel = config.locationLabel
	form.CategoryLabel = config.categoryLabel
	form.DescriptionLabel = config.descriptionLabel
	form.SubmitLabel = config.submitLabel
	form.Categories = config.categories
	form.CategoryValid = categoryIsValid(form.Category, config.categories)
	renderTemplate(w, status, config.formTemplate, form)
}

func renderList(w http.ResponseWriter, config formConfig) {
	data := ListData{
		PageTitle:      config.listTitle,
		CurrentPage:    strings.TrimPrefix(config.listURL, "/"),
		FormURL:        config.formAction,
		FormLabel:      config.formLabel,
		OtherListURL:   config.otherListURL,
		OtherListLabel: config.otherListLabel,
		EmptyMessage:   config.emptyMessage,
		Items:          submissionsFor(config.kind),
	}
	renderTemplate(w, http.StatusOK, config.listTemplate, data)
}

func categoryIsValid(category string, options []string) bool {
	for _, option := range options {
		if category == option {
			return true
		}
	}
	return false
}
