package form

type Errors struct {
	Name  string
	Phone string
}

func (e Errors) Any() bool {
	return e.Name != "" || e.Phone != ""
}

type State struct {
	Name    string
	Phone   string
	Errors  Errors
	Success string
}
