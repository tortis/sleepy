package sleepy

type Call struct {
	path    string
	method  string
	handler Handler
	filters []Filter
	model   callDataModel
}

type callDataModel struct {
	bodyIn    interface{}
	bodyOut   interface{}
	pathVars  []inputVar
	queryVars []inputVar
}

type inputVar struct {
	typ      string
	name     string
	desc     string
	required bool
}

// Call builder method
func (c *Call) Method(method string) *Call {
	c.method = method
	return c
}

// Call builder method
func (c *Call) To(fn Handler) *Call {
	c.handler = fn
	return c
}

// Call builder method
func (c *Call) PathParam(name, desc string) *Call {
	c.model.pathVars = append(c.model.pathVars, inputVar{
		typ:      "string",
		name:     name,
		desc:     desc,
		required: true})
	return c
}

// Call builder method
func (c *Call) Reads(m interface{}) *Call {
	c.model.bodyIn = m
	return c
}

// Call builder method
func (c *Call) Returns(m interface{}) *Call {
	c.model.bodyOut = m
	return c
}

// -- Example data model

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
		PathParam("uid", "ID of the user to search for.").
		Returns(User{})

	res.Route("").
		Method("POST").
		To(u.createUser).
		Reads(User{}).
		Returns(User{})
	return res
}

func (u *UserResource) getUser(resp Response, req *Request) {

}

func (u *UserResource) createUser(resp Response, req *Requeset) {

}
