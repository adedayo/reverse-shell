package reverse

//ShellOut output of a shell command run
type ShellOut struct {
	User, Dir, Hostname, StdOut, StdErr string
}
