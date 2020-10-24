package favicon

import (
	"os"
	"net/http"
)

// have your image in /kodata-path/images/fav.png
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	koPath := os.Getenv("KO_DATA_PATH")
	http.ServeFile(w, r, koPath+"/images/fav.png")
}
