package ifaces

// IGalaxyCache describes the fields required to be considered a cached Galaxy
type IGalaxyCache interface {
	Players() IPlayerCache
	Name() string
	Path() string
}
