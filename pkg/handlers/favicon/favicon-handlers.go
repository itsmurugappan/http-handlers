package favicon

import (
	"os"
	"net/http"
)

// have your image in /opt/images/fav.png
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	koPath := os.Getenv("KO_DATA_PATH")
	http.ServeFile(w, r, koPath+"/images/fav.png")
}
