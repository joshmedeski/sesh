package namer

// Gets the name from a directory
func dirName(n *RealNamer, path string) (string, error) {
	name := n.pathwrap.Base(path)
	return name, nil
}
