package mga

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/lib"
)

func getQuoterFileFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	const (
		filePath = "products/life/v2/wopta-per-te-vita-v42.xltx"
		filename = "Wopta per te. Vita - V4.2.xltx"
	)
	var err error

	log.AddPrefix("GetQuoterFileFx")
	log.Printf("Handler start -------------------------------------------------")
	defer func() {
		if err != nil {
			log.ErrorF("error: %s", err.Error())
		}
		log.Printf("Handler end -----------------------------------------------")
		log.PopPrefix()
	}()

	rawDoc := lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filePath, "")

	outMap := map[string]string{
		"filename": filename,
		"rawDoc":   base64.StdEncoding.EncodeToString(rawDoc),
	}

	rawMap, err := json.Marshal(outMap)

	return string(rawMap), outMap, err
}
