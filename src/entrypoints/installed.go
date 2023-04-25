package entrypoints

import "net/http"

func (hc *HandlerCreator) InstalledHandler(w http.ResponseWriter, _ *http.Request) {
	hc.httpSuccess(w, []byte("{}"))
}
