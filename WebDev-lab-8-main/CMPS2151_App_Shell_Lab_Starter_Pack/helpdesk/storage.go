package main

import "sync"

type Submission struct {
    Kind        string
    Title       string
    Location    string
    Category    string
    Description string
}

var submissions []Submission
var submissionsMu sync.Mutex

// addSubmission stores one validated record in this process.
func addSubmission(record Submission) {
    submissionsMu.Lock()
    submissions = append(submissions, record)
    submissionsMu.Unlock()
}

// submissionsFor returns a snapshot for one form, in submission order.
func submissionsFor(kind string) []Submission {
    submissionsMu.Lock()
    var result []Submission
    for _, record := range submissions {
        if record.Kind == kind {
            result = append(result, record)
        }
    }
    submissionsMu.Unlock()
    return result
}
