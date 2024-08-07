package mga

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/wopta/goworkspace/lib"
)

func GetQuoterFileFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	const (
		filePath = "products/life/v2/wopta-per-te-vita-v42.xltx"
		filename = "Wopta per te. Vita - V4.2.xltx"
	)
	var err error

	log.SetPrefix("[GetQuoterFileFx] ")
	log.Printf("Handler start -------------------------------------------------")
	defer func() {
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		log.Printf("Handler end -----------------------------------------------")
		log.SetPrefix("")
	}()

	rawDoc := lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filePath, "")

	outMap := map[string]string{
		"filename": filename,
		"rawDoc":   base64.StdEncoding.EncodeToString(rawDoc),
	}

	rawMap, err := json.Marshal(outMap)

	return string(rawMap), outMap, err
}
