package main

import (
	"github.com/gorilla/mux"
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
	"encoding/json"
	"errors"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(r *mux.Router, db *sql.DB) {
	a.DB = db
	a.Router = r
	a.initializeRoutes()
}

func (a *App) initializeRoutes()  {
	a.Router.HandleFunc("/users", RecoverWrap(a.getUsers)).Methods("GET")
	a.Router.HandleFunc("/user", RecoverWrap(a.createUser)).Methods("POST")
	a.Router.HandleFunc("/user/{id:[0-9]+}", RecoverWrap(a.getUser)).Methods("GET")
	a.Router.HandleFunc("/user/{id:[0-9]+}", RecoverWrap(a.updateUser)).Methods("PUT")
	a.Router.HandleFunc("/user/{id:[0-9]+}", RecoverWrap(a.deleteUser)).Methods("DELETE")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	u := user{ID: id}
	if err := getUser(&u, a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func checkRequestBody(w http.ResponseWriter, r *http.Request){
	if r == nil {
		respondWithError(w, http.StatusInternalServerError, "The request body is nil")
		return
	}
}

func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	checkRequestBody(w, r)
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	users, err := getUsers(start, count, a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request){
	checkRequestBody(w, r)
	var u user
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Payload")
		return
	}
	defer r.Body.Close()
	if err := createUser(&u, a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, u)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request)  {
	checkRequestBody(w, r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	var u user
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil{
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	defer r.Body.Close()
	u.ID = id
	if err := updateUser(&u, a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request)  {
	checkRequestBody(w, r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Id")
		return
	}
	u := user{ID: id}
	if err := deleteUser(&u, a.DB); err != nil{
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func RecoverWrap(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				log.Println(err.Error())
				respondWithError(w, http.StatusInternalServerError, err.Error())
			}
		}()
		f(w, r)
	})
}
