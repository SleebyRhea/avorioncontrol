package ifaces

// IGalaxyCache describes the fields required to be considered a cached Galaxy
type IGalaxyCache interface {
	Name() string
	Path() string
}
