package sleepy

type Resource struct {
	path  string
	name  string
	calls []*Call
}

func (r *Resource) Route(path string) *Call {
	c := &Call{path: path}
	r.calls = append(r.calls, c)
	return c
}
