package forms

import "testing"

func TestValidateServerFields(t *testing.T) {
	errs := Validate(FormSet{
		Common: CommonForm{Name: "x", Type: KindServer},
		Server: ServerForm{},
	})

	if errs["server_host"] != errRequired {
		t.Fatalf("server_host error = %q", errs["server_host"])
	}
	if errs["server_port"] != errRequired {
		t.Fatalf("server_port error = %q", errs["server_port"])
	}
	if errs["storage_bucket"] != "" {
		t.Fatalf("unexpected storage_bucket error = %q", errs["storage_bucket"])
	}
}

func TestValidateRequiresType(t *testing.T) {
	errs := Validate(FormSet{
		Common: CommonForm{Name: "x"},
	})

	if errs["type"] != errRequired {
		t.Fatalf("type error = %q", errs["type"])
	}
}
