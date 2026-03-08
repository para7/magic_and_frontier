package forms

type Kind string

const (
	KindNone    Kind = ""
	KindServer  Kind = "server"
	KindStorage Kind = "storage"
	KindMessage Kind = "message"
)

type CommonForm struct {
	Name string
	Type Kind
}

type ServerForm struct {
	Host string
	Port string
}

type StorageForm struct {
	Bucket string
	Region string
}

type MessageForm struct {
	Topic string
	QoS   string
}

type FormSet struct {
	Common  CommonForm
	Server  ServerForm
	Storage StorageForm
	Message MessageForm
}

type FieldErrors map[string]string

func (e FieldErrors) Get(field string) string {
	return e[field]
}

func (e FieldErrors) Has(field string) bool {
	_, ok := e[field]
	return ok
}

func (e FieldErrors) Any() bool {
	return len(e) > 0
}
