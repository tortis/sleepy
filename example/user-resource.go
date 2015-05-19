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

func (u *UserResource) getUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "getUser!")
}

func (u *UserResource) createUser(w http.ResponseWriter, r *http.Request) {
	//newUser := req.Data.(User)

	fmt.Fprintf(w, "createUser!")
}
