package forms

import (
	"net/http"

	"tools2/app/internal/domain/common"
)

const (
	errRequired = "1文字以上で入力してください"
	errInvalid  = "入力内容を確認してください"
)

type commonInput struct {
	Name string `json:"name" validate:"trimmed_required"`
	Type Kind   `json:"type" validate:"required,oneof=server storage message"`
}

type serverInput struct {
	Host string `json:"server_host" validate:"trimmed_required"`
	Port string `json:"server_port" validate:"trimmed_required"`
}

type storageInput struct {
	Bucket string `json:"storage_bucket" validate:"trimmed_required"`
	Region string `json:"storage_region" validate:"trimmed_required"`
}

type messageInput struct {
	Topic string `json:"message_topic" validate:"trimmed_required"`
	QoS   string `json:"message_qos" validate:"trimmed_required"`
}

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
	collectFormErrors(errs, common.ValidateStruct(commonInput{
		Name: fs.Common.Name,
		Type: fs.Common.Type,
	}))

	switch fs.Common.Type {
	case KindServer:
		collectFormErrors(errs, common.ValidateStruct(serverInput{
			Host: fs.Server.Host,
			Port: fs.Server.Port,
		}))
	case KindStorage:
		collectFormErrors(errs, common.ValidateStruct(storageInput{
			Bucket: fs.Storage.Bucket,
			Region: fs.Storage.Region,
		}))
	case KindMessage:
		collectFormErrors(errs, common.ValidateStruct(messageInput{
			Topic: fs.Message.Topic,
			QoS:   fs.Message.QoS,
		}))
	}

	return errs
}

func collectFormErrors(errs FieldErrors, violations []common.FieldViolation) {
	for _, violation := range violations {
		if _, exists := errs[violation.Field]; exists {
			continue
		}
		switch violation.Tag {
		case "required", "trimmed_required":
			errs[violation.Field] = errRequired
		default:
			errs[violation.Field] = errInvalid
		}
	}
}
