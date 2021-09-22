package dto

type (
	User struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
		Role  int    `json:"role,omitempty"`
	}
)
