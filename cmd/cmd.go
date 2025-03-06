package cmd

// CMD is an interface for all executable commands in the cmd package.
type CMD interface {
	Main() error
}
