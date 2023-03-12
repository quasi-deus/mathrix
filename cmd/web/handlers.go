package main
import (
	"errors"
	"net/http"
	"strconv"
	"fmt"
	"time"

	"mathrix.ceg.com/internal/models"
	"mathrix.ceg.com/internal/validator"

	"github.com/julienschmidt/httprouter"
)
type userSignupForm struct{
	UserID int `form:"userid"`
	Name string `form:"name"`
	URN string `form:"urn"`
	Phone int64 `form:"phone"`
	College string `form:"college"`
	Dept string `form:"dept"`
	Year int `form:"year"`
	Degree string `form:"degree"`
	Email string `form:"email"`
	Password string `form:"password"`
	Authority bool `form:"authority"`
	validator.Validator `form:"-"`
}

type userLoginForm struct{
	URN string `form:"urn"`
	Password string `form:"password"`
	validator.Validator `form:"-"`
}

type eventCreateForm struct{
	EventID int `form:"eventid"`
	EventName string `form:"eventname"`
	Content string `form:"content"`
	Venue string `form:"venue"`
	Technicality bool `form:"technicality"`
	EventDate time.Time `form:"eventdate"`
	validator.Validator `form:"-"`
}
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	
	events, err:=app.events.ViewAll()
	if err!=nil{
		app.serverError(w, err)
		return
	}

	data:=app.newTemplateData(r)
	data.Events=events

	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) event(w http.ResponseWriter, r *http.Request) {
	events, err:=app.events.ViewAll()
	if err!=nil{
		app.serverError(w, err)
		return
	}

	data:=app.newTemplateData(r)
	data.Events=events

	app.render(w, http.StatusOK, "event.tmpl", data)
}

func (app *application) user(w http.ResponseWriter, r *http.Request) {
	users, err:=app.users.ViewAll()
	if err!=nil{
		app.serverError(w, err)
		return
	}

	data:=app.newTemplateData(r)
	data.Users=users

	app.render(w, http.StatusOK, "user.tmpl", data)
}

func (app *application) eventView(w http.ResponseWriter, r *http.Request) {
	
	params:=httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	event, err := app.events.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data:=app.newTemplateData(r)
	data.Event=event
	
	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) eventCreate(w http.ResponseWriter, r *http.Request) {
	data:=app.newTemplateData(r)
	data.Form = eventCreateForm{
		Technicality: true,
	}
	app.render(w, http.StatusOK, "create.tmpl", data)
}
func (app *application) eventCreatePost (w http.ResponseWriter, r *http.Request) {
	var form eventCreateForm
	err:=app.decodePostForm(r, &form)
	if err!=nil{
		app.clientError(w, http.StatusBadRequest)
		return
	}
	
	form.CheckField(validator.NotBlank(form.EventName), "eventname", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.EventName, 100), "eventname", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Venue), "venue", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Technicality, true, false), "technicality", "This field can only have technical or non-technical")

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}
	id, err := app.events.Insert(form.EventName, form.Content, form.Venue, form.Technicality, form.EventDate)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Event has been successfully created")
	http.Redirect(w, r, fmt.Sprintf("/event/view/%d", id), http.StatusSeeOther)
}

func (app *application) eventUpdate(w http.ResponseWriter, r *http.Request) {
	params:=httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	event, err := app.events.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data:=app.newTemplateData(r)
	data.Event=event
	data.Form = eventCreateForm{
		EventID: event.EventID,
		EventName: event.EventName,
		Content: event.Content,
		Venue: event.Venue,
		Technicality: event.Technicality,
		EventDate: event.EventDate,
	}
	app.render(w, http.StatusOK, "update.tmpl", data)
}
func (app *application) eventUpdatePost (w http.ResponseWriter, r *http.Request) {
	var form eventCreateForm
	err:=app.decodePostForm(r, &form)
	if err!=nil{
		app.clientError(w, http.StatusBadRequest)
		return
	}
	
	form.CheckField(validator.NotBlank(form.EventName), "eventname", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.EventName, 100), "eventname", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Venue), "venue", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Technicality, true, false), "technicality", "This field can only have technical or non-technical")

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "update.tmpl", data)
		return
	}
	id, err := app.events.Update(form.EventID,form.EventName, form.Content, form.Venue, form.Technicality, form.EventDate)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Event has been successfully updated")
	http.Redirect(w, r, fmt.Sprintf("/event/view/%d", id), http.StatusSeeOther)
}

func (app *application) eventAddPost(w http.ResponseWriter, r *http.Request) {
	params:=httprouter.ParamsFromContext(r.Context())
	userid:=app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	eventid, err:=strconv.Atoi(params.ByName("id"))
	if err != nil || userid < 1 || eventid < 1{
		app.notFound(w)
		return
	}
	err=app.users.RegisterEvent(userid,eventid)
	if err!=nil{
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "You have registered for this event.")
	http.Redirect(w, r, "/event",http.StatusSeeOther)
}

func (app *application) userDelete(w http.ResponseWriter, r *http.Request) {
	params:=httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	err=app.users.Delete(id)
	if err!=nil{
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "You have deleted the user.")
	http.Redirect(w, r, "/user",http.StatusSeeOther)
}
func (app *application) userUpdate(w http.ResponseWriter, r *http.Request){
	data:=app.newTemplateData(r)
	params:=httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	user, err := app.users.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data.User=user
	data.Form = userSignupForm{
		UserID: user.UserID,
		Name: user.Name,
		URN: user.URN,
		Phone: user.Phone,
		College: user.College,
		Dept: user.Dept,
		Year: user.Year,
		Degree: user.Degree,
		Email: user.Email,
		Authority: user.Authority,
	}
	app.render(w, http.StatusOK, "signup.tmpl", data)
}
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	if(!app.isAuthorized(r) && app.isAuthenticated(r)){
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	data:=app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl", data)
}
func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.URN), "urn", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.College), "college", "This field cannot be blank")
	form.CheckField(validator.IntegerRange(form.Phone, 6000000000, 9999999999), "phone", "This field must be a valid phone number")
	form.CheckField(validator.PermittedValue(form.Year, 1,2,3,4,5), "year", "This field can only have values 1-5")
	form.CheckField(validator.NotBlank(form.Dept), "dept", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Degree), "degree", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}
	if(form.UserID!=0){
		err=app.users.Update(form.UserID, form.Name, form.URN, form.Phone, form.College, form.Dept, form.Year, form.Degree, form.Email, form.Password, form.Authority)
	}else{
		err=app.users.Insert(form.Name, form.URN, form.Phone, form.College, form.Dept, form.Year, form.Degree, form.Email, form.Password, form.Authority)
	}
	if err!=nil{
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	if(app.isAuthorized(r)){
		app.sessionManager.Put(r.Context(), "flash", "Database has been updated successfully.")
		http.Redirect(w, r, "/user",http.StatusSeeOther)
	}else{
		app.sessionManager.Put(r.Context(), "flash", "You have signed up successfully. Please log in.")
		http.Redirect(w, r, "/user/login",http.StatusSeeOther)
	}
}
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	if(app.isAuthenticated(r)){
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	data:=app.newTemplateData(r)
	data.Form=userLoginForm{}
	app.render(w, http.StatusOK,"login.tmpl", data)

}
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {

	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.URN), "urn", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.URN, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("URN or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	authority, err:=app.users.Authorized(id)
	if err!=nil{
		app.serverError(w, err)
		return
	}
	if authority {
		app.sessionManager.Put(r.Context(), "authorizedUserID", id)
	}
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/event", http.StatusSeeOther)
}
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err:=app.sessionManager.RenewToken(r.Context())
	if err!=nil{
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authorizedUserID")
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(),"flash","You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
