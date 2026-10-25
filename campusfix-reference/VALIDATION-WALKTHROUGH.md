# Trace one validation error

## First observe the running form

Fill Title, Category and Description with valid values. Put THREE SPACES in Room or location.
Submit. Spaces let this request pass the browser's required check, so the server can demonstrate its own rule.
Expected: a 422 response, red location border, red location message, retained valid values and no saved record.
An entirely empty required input is normally blocked by the browser before the server receives a request.

## 1 The handler creates a message

File: handlers.go. Function: reportCreatePost.

```go
location := strings.TrimSpace(form.Location)
if location == "" {
    form.LocationError = "Enter a location."
}
```

form.Location retains the original three spaces. The local location value is empty after trimming.
The form.LocationError string tells the template why this field failed.
All checks run. The final error branch calls:

```go
renderTemplate(w, 422, "report-form.tmpl", form)
return
```

That return prevents the handler reaching addSubmission. The 422 status describes the response.
A status code does not apply CSS or turn a border red.

## 2 The template creates HTML

File: templates/report-form.tmpl. Find the complete field wrapper for location.

```html
<div class="field">
  <label for="location">Room or location</label>
  <input id="location" name="location" value="{{.Location}}"
         {{if .LocationError}}
           aria-describedby="location-error"
           aria-invalid="true"
         {{end}}
         required>
  {{if .LocationError}}
    <p class="error" id="location-error">{{.LocationError}}</p>
  {{end}}
</div>
```

The template if checks whether this string is nonempty. This is a TEMPLATE rule; Go strings are not booleans.
For our invalid submission, Go emits this HTML:

```html
<div class="field">
  <label for="location">Room or location</label>
  <input id="location" name="location" value="   "
         aria-describedby="location-error" aria-invalid="true" required>
  <p class="error" id="location-error">Enter a location.</p>
</div>
```

The browser receives the resulting HTML. It never executes Go template instructions.
aria-describedby associates the control with the message. Matching id values are essential.

## 3 CSS matches the generated HTML

File: static/token.css.

```css
--danger: #b42318;
--danger-bg: #fff8f7;
```

File: static/main.css. Find the SERVER ERROR BORDER and SERVER ERROR MESSAGE comments.

```css
input[aria-invalid="true"],
select[aria-invalid="true"],
textarea[aria-invalid="true"] {
  border: 2px solid var(--danger);
  background: var(--danger-bg);
}
.error {
  color: var(--danger);
  font-size: .9rem;
  font-weight: 600;
  margin-top: 6px;
}
```

The attribute selector matches the failed control. It gives that control a red border.
The .error selector matches the paragraph's class. It makes the MESSAGE red.
aria-invalid has accessibility meaning, but it does not automatically give a field this red styling.
Our CSS supplies the visible appearance. The focus rules keep an invalid field's outline red while repairing it.
We do not use :invalid here because it represents the browser's validation state, which differs from the server's result.

## 4 Correct and submit again

Enter E201 and submit. A new request creates a new FormData value with empty error strings.
All rules pass, so the handler saves once, returns 303 with Location: /reports, and the browser requests GET /reports.
Opening a fresh blank form also produces no error attributes or error paragraphs.
No JavaScript removes the old message while you type; the next server response replaces the page.

## Keep the layout stable

Each .field contains its own label, control and conditional error paragraph.
.form-grid places FIELD WRAPPERS in its columns. It does not place error paragraphs in separate columns.
Keep error paragraphs inside their field wrapper when building the second form.

## Three small experiments

Make one change at a time in DevTools. Predict first and restore it before the next experiment.

1. Remove aria-invalid from the failed input. The red border disappears, but the red error paragraph remains.
2. Restore the attribute. Remove class="error" from the paragraph. Its red text disappears, but the field border stays red.
3. Inspect a fresh GET form. Neither aria-invalid nor the conditional error paragraph should be present.

Record what HTML matched each CSS selector. Explain the result before changing code.
If stuck, identify the first missing link: submitted value, error string, generated attribute/paragraph, or matching CSS rule.
Then ask a precise question about that link.
