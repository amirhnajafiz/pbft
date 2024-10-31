package cmd

// CMD is an interface for all executables in cmd module.
type CMD interface {
	Main() error
}
