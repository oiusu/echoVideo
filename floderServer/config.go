package floderServer

import (
	"fmt"
	"net/http"
	"os"
)

/**
 * Created by chenc on 2018/9/21
 */

var (
	root_folder  *string // TODO: Find a way to be cleaner !
	uses_gzip    *bool
	template_dir *string
	Url_prefix   *string
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
	Url_prefix = new(string)
	uses_gzip = new(bool)

	*root_folder = "/Users/chenc/go/src/echoVideo/oss/"
	*template_dir = cwd + "/view/"
	*uses_gzip = true
	*Url_prefix = "/echoVideo"

	fmt.Printf("Sharing %s ...\n", *root_folder)
}

func ChangeRootPath(w http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" {
		req.ParseMultipartForm(32 << 20)
		rootPath := req.FormValue("rootPath")
		*root_folder = rootPath

		//reqPath := path.Clean(req.URL.Path)
		http.Redirect(w, req, *Url_prefix, http.StatusTemporaryRedirect)
	}
}
