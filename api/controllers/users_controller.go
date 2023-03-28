package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"strconv"

	"github.com/ahmed-deftoner/blogs-backend/api/auth"
	"github.com/ahmed-deftoner/blogs-backend/api/models"
	"github.com/ahmed-deftoner/blogs-backend/api/response"
	"github.com/ahmed-deftoner/blogs-backend/api/utils/formaterror"
	"github.com/gorilla/mux"
)

var (
	authUserName      = "AKIA2HQWJ6EXY5XVJQZ6"
	authPassword      = "BBcZEJ4MNjCjurtvP46F2Mr3Q+gvSkdwu5MxXO8RmWNt"
	smtpServerAddr    = "email-smtp.ap-south-1.amazonaws.com"
	smtpServerPort    = "587"
	destinationEmails = []string{"ahmedghtwhts786@gmail.com"}
	senderEmail       = "ahmedghtwhts786@gmail.com"
)

func (server *Server) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	err := auth.TokenValid(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	fmt.Println(uid)
	user := models.User{}
	updatedUser, err := user.ConfirmUser(server.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		response.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	fmt.Println(updatedUser)
	response.JSON(w, http.StatusOK, updatedUser)
}

func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user.Prepare()
	err = user.Validate("")
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userCreated, err := user.SaveUser(server.DB)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())

		response.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	tok, err := auth.CreateToken(userCreated.Id)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	url := "http://localhost:8080/confirm/?token=" + tok
	go func() {
		fmt.Println("sending emails example")

		msg := []byte("Subject: test email\r\n" +
			"\r\n" +
			"Here comes the email content\r\n<a href=\"" + url + "\">Click here to confirm your email</a>")

		auth := smtp.PlainAuth("", authUserName, authPassword, smtpServerAddr)

		err := smtp.SendMail(smtpServerAddr+":"+smtpServerPort,
			auth, senderEmail, destinationEmails, msg)

		if err != nil {
			fmt.Printf("Error to sending email: %s", err)
			return
		}

		fmt.Println("email sent success")
	}()
	fmt.Println(url)
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.Id))
	response.JSON(w, http.StatusCreated, userCreated)
}

func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {

	user := models.User{}

	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, users)
}

func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserById(server.DB, uint32(uid))
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	response.JSON(w, http.StatusOK, userGotten)
}

func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(uid) {
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedUser, err := user.UpdateUser(server.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		response.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	response.JSON(w, http.StatusOK, updatedUser)
}

func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	user := models.User{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != 0 && tokenID != uint32(uid) {
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = user.DeleteUser(server.DB, uint32(uid))
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	response.JSON(w, http.StatusNoContent, "")
}
