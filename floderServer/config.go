package floderServer

import (
	"fmt"
	"os"
)

/**
 * Created by chenc on 2018/9/21
 */

var (
	root_folder  *string // TODO: Find a way to be cleaner !
	uses_gzip    *bool
	template_dir *string
)

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error while getting current directory.")
		return
	}
	fmt.Println(cwd)

	root_folder = new(string)
	template_dir = new(string)
	uses_gzip = new(bool)

	*root_folder = "/Users/chenc/go/src/echoVideo/oss/"
	*template_dir = cwd + "/view/"
	*uses_gzip = true

	fmt.Printf("Sharing %s ...\n", *root_folder)
}
