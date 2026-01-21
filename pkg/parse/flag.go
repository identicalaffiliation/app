package parse

import "flag"

func FlagInit() string {
	var path *string = flag.String("path", "", "path to yaml config file")
	flag.Parse()

	return *path
}
