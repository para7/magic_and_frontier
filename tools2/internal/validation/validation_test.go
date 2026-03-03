package validation

import "testing"

func TestValidate(t *testing.T) {
	tests := []struct {
		name         string
		inputName    string
		inputPhone   string
		wantNameErr  bool
		wantPhoneErr bool
	}{
		{name: "both empty", inputName: "", inputPhone: "", wantNameErr: true, wantPhoneErr: true},
		{name: "phone empty", inputName: "a", inputPhone: "", wantNameErr: false, wantPhoneErr: true},
		{name: "name empty", inputName: "", inputPhone: "123", wantNameErr: true, wantPhoneErr: false},
		{name: "phone non digit", inputName: "a", inputPhone: "abc", wantNameErr: false, wantPhoneErr: true},
		{name: "valid", inputName: "a", inputPhone: "09012345678", wantNameErr: false, wantPhoneErr: false},
		{name: "name spaces", inputName: "   ", inputPhone: "123", wantNameErr: true, wantPhoneErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.inputName, tt.inputPhone)
			if got := err.Name != ""; got != tt.wantNameErr {
				t.Fatalf("name error = %v, want %v", got, tt.wantNameErr)
			}
			if got := err.Phone != ""; got != tt.wantPhoneErr {
				t.Fatalf("phone error = %v, want %v", got, tt.wantPhoneErr)
			}
		})
	}
}
