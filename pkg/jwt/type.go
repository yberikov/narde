package jwt

type Type string

const (
	Refresh Type = "REFRESH"
	Access  Type = "ACCESS"
)
