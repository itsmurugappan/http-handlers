package favicon

import (
	"net/http"
	"os"
)

// have your image in /kodata-path/images/fav.png
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	koPath := os.Getenv("KO_DATA_PATH")
	http.ServeFile(w, r, koPath+"/images/fav.png")
}

// HandleFavicon handles favicon requests
func HandleFavicon() {
	http.HandleFunc("/favicon.ico", FaviconHandler)
}
