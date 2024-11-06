package cmd

// CMD is an interface for all executables in cmd package.
type CMD interface {
	Main() error
}
