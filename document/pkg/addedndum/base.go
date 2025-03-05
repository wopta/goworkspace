package addedndum

import (
	"fmt"
	"strings"
	"time"

	"github.com/wopta/goworkspace/document/internal/engine"
	"github.com/wopta/goworkspace/models"
)

const (
	tabDimension = 15
)

type baseGenerator struct {
	engine      *engine.Fpdf
	now         time.Time
	signatureID uint32
	networkNode *models.NetworkNode
	policy      *models.Policy
}

func (bg *baseGenerator) Save(rawDoc []byte) (string, error) {
	filename := strings.ReplaceAll(fmt.Sprintf("%s/%s/"+models.ProposalDocumentFormat, "temp", bg.policy.Uid,
		bg.policy.NameDesc, bg.policy.ProposalNumber), " ", "_")
	return bg.engine.Save(rawDoc, filename)
}
