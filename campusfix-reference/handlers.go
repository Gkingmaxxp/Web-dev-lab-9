package main

import (
    "net/http"
    "strings"
    "unicode/utf8"
)

type FormData struct {
    Title, Location, Category, Description string
    TitleError, LocationError, CategoryError, DescriptionError string
}

type ListData struct {
    Items []Submission
}

func reportCreate(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, 200, "report-form.tmpl", FormData{})
}

func reportCreatePost(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }
    // Keep original strings so a 422 can show what the person entered.
    form := FormData{
        Title: r.PostForm.Get("title"),
        Location: r.PostForm.Get("location"),
        Category: r.PostForm.Get("category"),
        Description: r.PostForm.Get("description"),
    }
    // Check trimmed text, but preserve the original strings in form.
    title := strings.TrimSpace(form.Title)
    location := strings.TrimSpace(form.Location)
    category := form.Category
    description := strings.TrimSpace(form.Description)

    // A message here becomes HTML, then a red border and red text.
    if title == "" {
        form.TitleError = "Enter a title."
    } else if utf8.RuneCountInString(title) > 80 {
        form.TitleError = "Use at most 80 code points."
    }
    if location == "" {
        form.LocationError = "Enter a location."
    } else if utf8.RuneCountInString(location) > 40 {
        form.LocationError = "Use at most 40 code points."
    }
    if category != "Electrical" && category != "Plumbing" &&
        category != "Furniture" && category != "Other" {
        form.CategoryError = "Choose a listed category."
    }
    if description == "" {
        form.DescriptionError = "Enter a description."
    } else if utf8.RuneCountInString(description) > 500 {
        form.DescriptionError = "Use at most 500 code points."
    }
    // Finish every check before rendering. No record is saved on this path.
    if form.TitleError != "" || form.LocationError != "" ||
        form.CategoryError != "" || form.DescriptionError != "" {
        renderTemplate(w, 422, "report-form.tmpl", form)
        return
    }
    // Only valid input reaches this save. Kind comes from our own handler.
    addSubmission(Submission{
        Kind: "report",
        Title: title,
        Location: location,
        Category: category,
        Description: description,
    })
    // Tell the browser to send GET /reports. The browser makes that request.
    http.Redirect(w, r, "/reports", http.StatusSeeOther)
    return
}

func listReports(w http.ResponseWriter, r *http.Request) {
    data := ListData{Items: submissionsFor("report")}
    renderTemplate(w, 200, "reports.tmpl", data)
}
