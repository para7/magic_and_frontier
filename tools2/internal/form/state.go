package form

type Errors struct {
	Name      string
	Phone     string
	Mode      string
	Latitude  string
	Longitude string
	Birthdate string
	Height    string
	Weight    string
}

func (e Errors) Any() bool {
	return e.Name != "" ||
		e.Phone != "" ||
		e.Mode != "" ||
		e.Latitude != "" ||
		e.Longitude != "" ||
		e.Birthdate != "" ||
		e.Height != "" ||
		e.Weight != ""
}

type State struct {
	Name      string
	Phone     string
	Mode      string
	Latitude  string
	Longitude string
	Birthdate string
	Height    string
	Weight    string
	Errors    Errors
	Success   string
}

const (
	ModeLatLng       = "latlng"
	ModeBirthdate    = "birthdate"
	ModeHeightWeight = "heightweight"
)

func IsValidMode(mode string) bool {
	return mode == ModeLatLng || mode == ModeBirthdate || mode == ModeHeightWeight
}

func NormalizeMode(mode string) string {
	if IsValidMode(mode) {
		return mode
	}
	return ModeLatLng
}
