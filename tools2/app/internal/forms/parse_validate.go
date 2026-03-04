package forms

import (
	"net/http"
	"strings"
)

const (
	errRequired = "1文字以上で入力してください"
)

func FromQuery(r *http.Request) FormSet {
	q := r.URL.Query()
	return FormSet{
		Common: CommonForm{
			Name: q.Get("name"),
			Type: Kind(q.Get("type")),
		},
		Server: ServerForm{
			Host: q.Get("server_host"),
			Port: q.Get("server_port"),
		},
		Storage: StorageForm{
			Bucket: q.Get("storage_bucket"),
			Region: q.Get("storage_region"),
		},
		Message: MessageForm{
			Topic: q.Get("message_topic"),
			QoS:   q.Get("message_qos"),
		},
	}
}

func FromPostForm(r *http.Request) FormSet {
	return FormSet{
		Common: CommonForm{
			Name: r.PostFormValue("name"),
			Type: Kind(r.PostFormValue("type")),
		},
		Server: ServerForm{
			Host: r.PostFormValue("server_host"),
			Port: r.PostFormValue("server_port"),
		},
		Storage: StorageForm{
			Bucket: r.PostFormValue("storage_bucket"),
			Region: r.PostFormValue("storage_region"),
		},
		Message: MessageForm{
			Topic: r.PostFormValue("message_topic"),
			QoS:   r.PostFormValue("message_qos"),
		},
	}
}

func Validate(fs FormSet) FieldErrors {
	errs := FieldErrors{}

	requireField(errs, "name", fs.Common.Name, errRequired)
	requireField(errs, "type", string(fs.Common.Type), errRequired)

	switch fs.Common.Type {
	case KindServer:
		requireField(errs, "server_host", fs.Server.Host, errRequired)
		requireField(errs, "server_port", fs.Server.Port, errRequired)
	case KindStorage:
		requireField(errs, "storage_bucket", fs.Storage.Bucket, errRequired)
		requireField(errs, "storage_region", fs.Storage.Region, errRequired)
	case KindMessage:
		requireField(errs, "message_topic", fs.Message.Topic, errRequired)
		requireField(errs, "message_qos", fs.Message.QoS, errRequired)
	}

	return errs
}

func requireField(errs FieldErrors, key, value, msg string) {
	if strings.TrimSpace(value) == "" {
		errs[key] = msg
	}
}
