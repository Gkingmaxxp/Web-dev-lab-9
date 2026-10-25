# CampusFix first form starter

This project implements the first CampusFix form completely. Run it before changing it.
The purple rail, purple buttons, gold university subtitle and pale background follow the supplied reference.
The lab retains its existing four-field contract. The reference image's additional Building field is not part of this exercise.

1. Open a terminal inside this folder.
2. Run `go run .` with Go 1.26 or later.
3. Open http://localhost:4000/reports/create.
4. Submit valid values. Read the list at http://localhost:4000/reports.
5. Open `VALIDATION-WALKTHROUGH.md` and perform its experiments.
6. Build your application's second form using its exact mapping in the lab.

Files to read in order:
- `main.go`: the three routes and static-file route.
- `templates/report-form.tmpl`: the controls and conditional error markup.
- `handlers.go`: raw values, validation, 422, save, 303.
- `static/token.css`: named colors.
- `static/main.css`: form layout, red borders and red error messages.
- `storage.go`: the supplied in-memory helper.
- `render.go`: the supplied template buffer helper.

The first form uses four posted names: title, location, category and description.
Do not copy a second main function or duplicate FormData/ListData types into your existing app.
Keep your previous reset and app shell when transferring the form. The standalone starter includes its own small reset.
