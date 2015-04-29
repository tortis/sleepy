package main

import "../sleeyp"

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

func (u *UserResource) Generate() *Resource {
	res := new(Resource)
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

func (u *UserResource) getUser(resp Response, req *Request) {

}

func (u *UserResource) createUser(resp Response, req *Requeset) {

}
