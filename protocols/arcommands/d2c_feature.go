package arcommands

// D2CFeature ...
// TODO: Document this
type D2CFeature interface {
	ID() uint8
	Name() string
	D2CClasses() map[uint8]D2CClass
}
