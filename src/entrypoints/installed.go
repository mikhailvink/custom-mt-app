package entrypoints

import "net/http"

func InstalledHandler(w http.ResponseWriter, _ *http.Request) {
	httpSuccess(w, []byte("{}"))
}
