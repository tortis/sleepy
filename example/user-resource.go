package main

import (
	"fmt"
	"net/http"

	"github.com/tortis/sleepy"
)

type User struct {
	Id        string `bson:"_id" sleepy:"readonly"`
	FirstName string `sleepy:"required"`
	LastName  string `sleepy:"required"`
	Email     string `sleepy:"required"`
	Password  string `json:",omitempty" sleepy:"required,writeonly"`
}

type UserResource struct {
	dbs int
}

func (u *UserResource) Generate() *sleepy.Resource {
	res := sleepy.NewResource("/users")
	res.Filter(hasAuthFilter)
	res.Route("/{uid}").
		Method("GET").
		To(u.getUser).
		OperationName("getUser").
		PathParam("uid", "ID of the user to search for.").
		Returns(User{})

	res.Route("").
		Method("POST").
		To(u.createUser).
		OperationName("createUser").
		Reads(User{}).
		Returns(User{})
	return res
}

func (u *UserResource) getUser(w http.ResponseWriter, r *http.Request, d sleepy.CallData) (interface{}, sleepy.Error) {
	return fmt.Sprintf("getUser! - %s", d["auth"].(string)), nil
}

func (u *UserResource) createUser(w http.ResponseWriter, r *http.Request, d sleepy.CallData) (interface{}, sleepy.Error) {
	return "Create User!", nil
}

func hasAuthFilter(w http.ResponseWriter, r *http.Request, d sleepy.CallData) sleepy.Error {
	if r.Header.Get("Authorization") == "" {
		http.Error(w, "Authorization header is missing.", http.StatusBadRequest)
		return ErrLogin
	}
	d["auth"] = r.Header.Get("Authorization")
	return nil
}
