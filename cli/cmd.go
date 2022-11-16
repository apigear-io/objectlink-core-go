package cli

type Command struct {
	Usage string
	Names []string
	Exec  func(args []string) error
	Help  string
}
