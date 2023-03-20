package render

type InMsg interface {
	GetInfo() []string
}

type OutMsg interface {
	ToString() string
}
