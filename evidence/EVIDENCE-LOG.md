## Experiment 1 — Blank form GET
**Prediction:** Both forms return 200, all fields blank, no error paragraphs, no records added.

**Method:** Opened each form and list in the browser; viewed page source on the ticket form.
**Observed:**
- GET /tickets/create → 200, blank form, no errors
- GET /software-requests/create → 200, blank form, no errors
- No `class="error"` elements in page source
- /tickets: "Saved records: 0" + "No support tickets have been submitted yet."
- /software-requests: "Saved records: 0" + "No software requests have been submitted yet."
**Explanation:** **Explanation:** `renderForm(... FormData{})` passes empty error strings, so every `{{if .XError}}` block is skipped. GET handlers never call `addSubmission`.

## Experiment 2 — Invalid values
**Prediction:** POST with spaces title + empty location returns 422, both errors shown, valid category/description kept, no record saved.

**Method:** curl POST with title=+++ and location= to both forms; also submitted in browser.

**Observed:**
- POST /tickets/create → 422
- Body has aria-invalid="true" on #title and #location
- Errors: "Enter an issue title." and "Enter a room or location."
- category=Network retained (selected); description "Test description" retained
- Follow-up GET /tickets → still "Saved records: 0"
- Same on /software-requests/create with category=Installation

**Explanation:** `TrimSpace("   ")` → "" triggers TitleError; `location == ""` triggers LocationError. Any error → `renderForm(..., 422, ...)` before `addSubmission`.

## Experiment 3 — Forged category
**Prediction:** POST category=NotAllowed returns 422, category error shown, no forged option in HTML, no record saved.

**Method:** curl POST with category=NotAllowed to both forms.

**Observed:**
- POST /tickets/create → 422
- Contains `<p class="error" id="category-error">Choose a listed support category.</p>`
- No `<option value="NotAllowed">` anywhere in the response
- Placeholder `<option value="" selected>` is active
- No record added (list count unchanged)
- Same on /software-requests/create

**Explanation:** `categoryIsValid("NotAllowed", ...)` returns false → CategoryError set, CategoryValid=false → placeholder gets `selected`. The `{{range .Categories}}` loop only emits allowlisted options.

## Experiment 4 — PRG
**Prediction:** POST returns 303 with Location: /tickets, browser follows with GET 200.

**Method:** Browser Network panel on a valid submission; curl without -L; follow-up GET with curl.

**Observed:**
- Browser: POST /tickets/create → 303, response header `location: /tickets`; browser auto-follows with GET /tickets → 200
- URL bar now http://localhost:4000/tickets
- curl -i (no -L): `HTTP/1.1 303 See Other` + `Location: /tickets`
- Follow-up `curl -i /tickets` → 200, "Projector broken" visible
- Same on /software-requests/create → 303 → Location: /software-requests

**Explanation:** `http.Redirect(w, r, "/tickets", http.StatusSeeOther)` emits the 303. Browser auto-issues GET to that Location. GET handler only reads the slice — no mutation.

## Experiment 5 — Refresh and separation
**Prediction:** Counts stay at 2 and 1 after three refreshes; each list shows only its own records.

**Method:** Restarted server, submitted 2 tickets (form 1) + 1 software request (form 2), refreshed each list 3× with Network panel open.

**Observed:**
- /tickets: "Saved records: 2" on all 3 refreshes; shows Wi-Fi down + Printer jam only
- /software-requests: "Saved records: 1" on all 3 refreshes; shows VS Code only
- No cross-contamination between kinds
- Network on refresh: Request Method = GET, Status = 200
- Counts do not increase on refresh

**Explanation:** `renderList(w, config)` calls `submissionsFor(config.kind)` which filters by Kind. GET handlers only read — never call `addSubmission`.

## Experiment 6 — Restart
**Prediction:** After restart both lists show 0 and their empty messages.

**Method:** Ctrl+C, then `go run .` again, then opened both lists.

**Observed:**
- Before restart: /tickets=2, /software-requests=1
- After restart: /tickets=0 with empty message; /software-requests=0 with empty message

**Explanation:** `var submissions []Submission` in storage.go lives in process memory only. On Ctrl+C the process exits and the slice is discarded. Nothing is written to disk.

## Boundary tests
**Prediction:** 80/40/500 return 303; 81/41/501 return 422.

**Method:** curl POST to both forms with exact boundary lengths; string lengths verified in PowerShell.

**Observed:**

Form 1 (tickets):
| Field | Length | Status | Expected |
|-------|--------|--------|----------|
| Title | 80 | 303 | 303 ✅ |
| Title | 81 | 422 | 422 ✅ |
| Location | 40 | 303 | 303 ✅ |
| Location | 41 | 422 | 422 ✅ |
| Description | 500 | 303 | 303 ✅ |
| Description | 501 | 422 | 422 ✅ |

Form 2 (software requests): identical results.

**Explanation:** `utf8.RuneCountInString` counts code points. Handler uses strict `>` comparisons, so 80 > 80 is false (pass) while 81 > 80 is true (fail). Same logic for 40/41 and 500/501.

## Browser space-only demo
**Prediction:** Typing three spaces in Title submits past the browser and server returns 422.

**Method:** Typed three spaces in Title, valid values elsewhere, submitted in browser.

**Observed:**
- Browser submitted (no HTML5 validation error)
- Server returned 422 with "Enter an issue title."
- Title field shows red border + error paragraph

**Explanation:** `required` only rejects the truly empty string; spaces count as content. `strings.TrimSpace("   ")` → "" → handler rejects. This is why the server must re-validate after trimming.

