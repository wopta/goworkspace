package addedndum

import (
	"time"

	"github.com/wopta/goworkspace/document/internal/dto"
	"github.com/wopta/goworkspace/document/internal/engine"
	"github.com/wopta/goworkspace/models"
)

type LifeAddendumGenerator struct {
	*baseGenerator
	dto *dto.BeneficiariesDTO
}

func NewLifeAddendumGenerator(engine *engine.Fpdf, policy *models.Policy, node *models.NetworkNode,
	product models.Product) *LifeAddendumGenerator {
	LifeAddendumDTO := dto.NewBeneficiariesDto()
	LifeAddendumDTO.FromPolicy(*policy, product)
	return &LifeAddendumGenerator{
		baseGenerator: &baseGenerator{
			engine:      engine,
			now:         time.Now(),
			signatureID: 0,
			networkNode: node,
			policy:      policy,
		},
		dto: LifeAddendumDTO,
	}
}

func (lag *LifeAddendumGenerator) Contract() ([]byte, error) {
	//lag.mainHeader()
	//
	//lag.engine.NewPage()
	//
	//lag.mainFooter()
	//
	//lag.engine.NewLine(10)
	//
	//lag.introSection()
	//
	//lag.engine.NewLine(10)
	//
	//lag.whoWeAreSection()
	//
	//lag.engine.NewLine(10)
	//
	//lag.insuredDetailsSection()
	//
	//lag.engine.NewPage()
	//
	//lag.guaranteesDetailsSection()
	//
	//lag.engine.NewPage()
	//
	//lag.deductibleSection()
	//
	//lag.engine.NewLine(10)
	//
	//lag.dynamicDeductibleSection()
	//
	//lag.engine.NewLine(5)
	//
	//lag.detailsSection()
	//
	//lag.engine.NewLine(5)
	//
	//lag.specialConditionsSection()
	//
	//lag.engine.NewLine(5)
	//
	//lag.bondSection()
	//
	//lag.engine.NewPage()
	//
	//lag.resumeSection()
	//
	//lag.engine.NewLine(5)
	//
	//lag.howYouCanPaySection()
	//
	//lag.engine.NewLine(5)
	//
	//lag.emitResumeSection()
	//
	//lag.engine.NewLine(5)
	//
	//lag.statementsFirstPart()
	//
	//lag.engine.NewPage()
	//
	//lag.claimsStatement()
	//
	//lag.statementsSecondPart()
	//
	//lag.engine.NewLine(3)
	//
	//lag.qbePrivacySection()
	//
	//lag.engine.NewPage()
	//
	//lag.qbePersonalDataSection()
	//
	//lag.engine.NewLine(5)
	//
	//lag.commercialConsentSection()
	//
	//lag.annexSections()
	//
	//lag.woptaHeader()
	//
	//lag.woptaFooter()
	//
	//lag.woptaPrivacySection()

	return lag.engine.RawDoc()
}
