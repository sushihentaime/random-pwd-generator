package main

import (
	"bytes"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/edwardbktay/golang-random-password-generator/templates"

	"github.com/edwardbktay/golang-random-password-generator/password"

	"github.com/go-playground/form/v4"
	"github.com/julienschmidt/httprouter"
)

type application struct {
	logger      *slog.Logger
	formDecoder *form.Decoder
}

type templateData struct {
	Form     any
	Password string
	Error    any
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	formDecoder := form.NewDecoder()

	app := &application{
		logger:      logger,
		formDecoder: formDecoder,
	}

	srv := http.Server{
		Addr:         ":4000",
		Handler:      app.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  5 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	app.logger.Info("starting server", "addr", srv.Addr)

	err := srv.ListenAndServe()
	if err != nil {
		app.logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodPost, "/", app.generatePasswordHandler)
	return router
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := templateData{}
	data.Form = passwordGenerateForm{
		Uppercase: true,
		Lowercase: true,
		Numbers:   true,
	}

	app.render(w, r, http.StatusOK, data)
}

type passwordGenerateForm struct {
	Length    int  `form:"length"`
	Uppercase bool `form:"uppercase"`
	Lowercase bool `form:"lowercase"`
	Numbers   bool `form:"numbers"`
	Special   bool `form:"symbols"`
}

func (app *application) generatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var form passwordGenerateForm
	var data templateData

	err := app.decodeForm(r, &form)
	if err != nil {
		app.badRequestErrorResponse(w)
		return
	}

	if form.Length < 6 || form.Length > 32 {
		data.Error = "The length must be between 6 and 32 characters"
		data.Form = form
		app.render(w, r, http.StatusBadRequest, data)
		return
	}

	if !form.Uppercase && !form.Lowercase && !form.Numbers && !form.Special {
		data.Error = "At least one character set must be selected"
		data.Form = form
		app.render(w, r, http.StatusBadRequest, data)
		return
	}

	options := &password.PasswordGeneratorOptions{
		Length:    form.Length,
		UpperCase: form.Uppercase,
		LowerCase: form.Lowercase,
		Numbers:   form.Numbers,
		Special:   form.Special,
	}

	generator := password.New(options)

	password := generator.Generate()

	data.Form = form
	data.Password = password

	app.render(w, r, http.StatusOK, data)
}

func (app *application) decodeForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecodeError *form.InvalidDecoderError
		if errors.As(err, &invalidDecodeError) {
			panic(err)
		}
		return err
	}

	return nil
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, data any) {
	tmpl, err := template.ParseFS(templates.File, "html/*.html")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, data)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		url    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error("internal server error", "method", method, "url", url, "error", err, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) badRequestErrorResponse(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
