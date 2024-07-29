package mga

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/wopta/goworkspace/lib"
)

func GetQuoterFileFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	const (
		basePath = "products/life/v2/"
		filename = "wopta-per-te-vita-pg.xltx"
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

	if authToken, err := lib.GetAuthTokenFromIdToken(r.Header.Get("Authorization")); err != nil || !authToken.IsNetworkNode {
		log.Printf("error validating token")
		err = errors.New("invalid token")
		return "", nil, err
	}

	rawDoc := lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), basePath+filename, "")

	outMap := map[string]string{
		"filename": filename,
		"rawDoc":   base64.StdEncoding.EncodeToString(rawDoc),
	}

	rawMap, err := json.Marshal(outMap)

	return string(rawMap), outMap, err
}
