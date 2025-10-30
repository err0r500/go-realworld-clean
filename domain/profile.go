package domain

type Profile struct {
	Username  string
	Bio       string
	Image     string
	Following bool // Whether the requesting user follows this profile
}
