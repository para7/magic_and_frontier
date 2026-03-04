package validation

import (
	"testing"

	"tools2/internal/form"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name          string
		state         form.State
		wantNameErr   bool
		wantPhoneErr  bool
		wantModeErr   bool
		wantLatErr    bool
		wantLngErr    bool
		wantBirthErr  bool
		wantHeightErr bool
		wantWeightErr bool
	}{
		{
			name: "latlng valid",
			state: form.State{
				Name: "a", Phone: "09012345678", Mode: form.ModeLatLng, Latitude: "35.0", Longitude: "139.0",
			},
		},
		{
			name: "birthdate valid",
			state: form.State{
				Name: "a", Phone: "09012345678", Mode: form.ModeBirthdate, Birthdate: "2000-01-01",
			},
		},
		{
			name: "heightweight valid",
			state: form.State{
				Name: "a", Phone: "09012345678", Mode: form.ModeHeightWeight, Height: "170", Weight: "60.5",
			},
		},
		{
			name: "invalid mode",
			state: form.State{
				Name: "a", Phone: "09012345678", Mode: "unknown",
			},
			wantModeErr: true,
		},
		{
			name: "latlng non number",
			state: form.State{
				Name: "a", Phone: "09012345678", Mode: form.ModeLatLng, Latitude: "x", Longitude: "",
			},
			wantLatErr: true, wantLngErr: true,
		},
		{
			name: "birthdate required",
			state: form.State{
				Name: "a", Phone: "09012345678", Mode: form.ModeBirthdate, Birthdate: "",
			},
			wantBirthErr: true,
		},
		{
			name: "heightweight non number",
			state: form.State{
				Name: "a", Phone: "09012345678", Mode: form.ModeHeightWeight, Height: "a", Weight: "",
			},
			wantHeightErr: true, wantWeightErr: true,
		},
		{
			name: "common validation still works",
			state: form.State{
				Name: " ", Phone: "abc", Mode: form.ModeBirthdate, Birthdate: "2000-01-01",
			},
			wantNameErr: true, wantPhoneErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.state)
			if got := err.Name != ""; got != tt.wantNameErr {
				t.Fatalf("name error = %v, want %v", got, tt.wantNameErr)
			}
			if got := err.Phone != ""; got != tt.wantPhoneErr {
				t.Fatalf("phone error = %v, want %v", got, tt.wantPhoneErr)
			}
			if got := err.Mode != ""; got != tt.wantModeErr {
				t.Fatalf("mode error = %v, want %v", got, tt.wantModeErr)
			}
			if got := err.Latitude != ""; got != tt.wantLatErr {
				t.Fatalf("latitude error = %v, want %v", got, tt.wantLatErr)
			}
			if got := err.Longitude != ""; got != tt.wantLngErr {
				t.Fatalf("longitude error = %v, want %v", got, tt.wantLngErr)
			}
			if got := err.Birthdate != ""; got != tt.wantBirthErr {
				t.Fatalf("birthdate error = %v, want %v", got, tt.wantBirthErr)
			}
			if got := err.Height != ""; got != tt.wantHeightErr {
				t.Fatalf("height error = %v, want %v", got, tt.wantHeightErr)
			}
			if got := err.Weight != ""; got != tt.wantWeightErr {
				t.Fatalf("weight error = %v, want %v", got, tt.wantWeightErr)
			}
		})
	}
}
