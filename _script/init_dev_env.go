package _script

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type woptaUser struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
}

type userEnv struct {
	WoptaUser   woptaUser
	DevEnvUsers []models.User
	NodeGraph   []models.NetworkNode
	Policies    []models.Policy
}

var (
	//roles  = []string{lib.UserRoleAdmin, lib.UserRoleManager, lib.UserRoleAgency, lib.UserRoleAgent + "online"}
	devEnv = make([]userEnv, 0)
)

func addPrefixToEmailAddress(address, prefix string) string {
	splittedAddress := strings.Split(address, "@")
	splittedAddress[0] += "+" + prefix
	return strings.Join(splittedAddress, "@")
}

func loadUsersFromPath(filePath string) ([]woptaUser, error) {
	var woptaUsers []woptaUser

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("can't open file path %s: %w", filePath, err)
	}
	defer jsonFile.Close()

	file, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("can't read file %v: %w", file, err)
	}

	err = json.Unmarshal(file, &woptaUsers)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal file %v into %v: %w", file, woptaUsers, err)
	}

	return woptaUsers, nil
}

func saveFileToPath(fileName, filePath string) error {
	byteBlob, err := json.Marshal(devEnv)
	if err != nil {
		return fmt.Errorf("error with json.Marshal: %w", err)
	}

	cwd, _ := os.Getwd()
	f, err := os.Create(filepath.Join(cwd, filePath, fileName))
	if err != nil {
		return fmt.Errorf("error with os.Create: %w", err)
	}
	defer f.Close()

	_, err = f.Write(byteBlob)
	if err != nil {
		return fmt.Errorf("error with f.Write: %w", err)
	}

	return nil
}

const alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func uniqueID() string {
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("dev env script: crypto/rand.Read error: %v", err))
	}
	for i, byt := range b {
		b[i] = alphanum[int(byt)%len(alphanum)]
	}
	return string(b)
}

func createUserWithRole(user woptaUser, role string) models.User {
	newUser := models.User{
		Uid:     uniqueID(),
		Name:    user.Name,
		Surname: user.Surname,
		Mail:    addPrefixToEmailAddress(user.Email, role),
		Role:    role,
	}
	return newUser
}

func createNetworkNode(userMail string, role string, parentUid string) models.NetworkNode {
	node := models.NetworkNode{
		Uid:          uniqueID(),
		Type:         role,
		Role:         role,
		Mail:         userMail,
		ParentUid:    parentUid,
		CreationDate: time.Now(),
		UpdatedDate:  time.Now(),
	}
	return node
}

func emptyPolicy(date time.Time) models.Policy {
	var p models.Policy
	p.Uid = uniqueID()
	p.EmitDate = date

	return p
}

func policyToSign(date time.Time) models.Policy {
	p := emptyPolicy(date)
	p.Status = models.PolicyStatusToSign
	p.StatusHistory = []string{models.PolicyStatusInitLead, models.PolicyStatusProposal,
		models.PolicyStatusContact, models.PolicyStatusToSign}
	p.CompanyEmit = true
	p.IsSign = false
	p.IsPay = false
	p.IsDeleted = false
	p.CompanyEmitted = false
	p.IsReserved = false
	p.Annuity = 0

	return p
}

func policyToPay(date time.Time) models.Policy {
	p := emptyPolicy(date)
	p.Status = models.PolicyStatusToPay
	p.StatusHistory = []string{models.PolicyStatusInitLead, models.PolicyStatusProposal,
		models.PolicyStatusContact, models.PolicyStatusToSign,
		models.PolicyStatusSign, models.PolicyStatusToPay}
	p.CompanyEmit = true
	p.IsSign = true
	p.IsPay = false
	p.IsDeleted = false
	p.CompanyEmitted = false
	p.IsReserved = false
	p.Annuity = 0

	return p
}

func policyPaid(date time.Time) models.Policy {
	p := emptyPolicy(date)
	p.Status = models.PolicyStatusPay
	p.StatusHistory = []string{models.PolicyStatusInitLead, models.PolicyStatusProposal,
		models.PolicyStatusContact, models.PolicyStatusToSign,
		models.PolicyStatusSign, models.PolicyStatusToPay,
		models.PolicyStatusPay}
	p.CompanyEmit = true
	p.IsSign = true
	p.IsPay = true
	p.IsDeleted = false
	p.CompanyEmitted = false
	p.IsReserved = false
	p.Annuity = 0

	return p
}

func policyUnresolved(date time.Time) models.Policy {
	p := emptyPolicy(date)
	p.Status = models.PolicyStatusUnsolved
	p.StatusHistory = []string{models.PolicyStatusInitLead, models.PolicyStatusProposal,
		models.PolicyStatusContact, models.PolicyStatusToSign,
		models.PolicyStatusSign, models.PolicyStatusToPay,
		models.PolicyStatusPay, models.PolicyStatusDraftRenew,
		models.PolicyStatusToPay, models.PolicyStatusUnsolved}
	p.CompanyEmit = true
	p.IsSign = true
	p.IsPay = false
	p.IsDeleted = false
	p.CompanyEmitted = true
	p.IsReserved = true
	p.Annuity = 1

	return p
}

func policyRenewed(date time.Time) models.Policy {
	p := emptyPolicy(date)
	p.Status = models.PolicyStatusRenewed
	p.StatusHistory = []string{models.PolicyStatusInitLead, models.PolicyStatusProposal,
		models.PolicyStatusContact, models.PolicyStatusToSign,
		models.PolicyStatusSign, models.PolicyStatusToPay,
		models.PolicyStatusPay, models.PolicyStatusDraftRenew,
		models.PolicyStatusToPay, models.PolicyStatusPay,
		models.PolicyStatusRenewed}
	p.CompanyEmit = true
	p.IsSign = true
	p.IsPay = true
	p.IsDeleted = false
	p.CompanyEmitted = true
	p.IsReserved = true
	p.Annuity = 1

	return p
}

func generatePolicies() []models.Policy {
	var policies []models.Policy

	//PolicyStatusToSign
	policies = append(policies, policyToSign(time.Now()))
	policies = append(policies, policyToSign(time.Now().AddDate(0, -3, 0)))

	//PolicyStatusToPay:
	policies = append(policies, policyToPay(time.Now()))

	//PolicyStatusPay:
	policies = append(policies, policyPaid(time.Now()))

	//PolicyStatusUnsolved:
	policies = append(policies, policyUnresolved(time.Now()))

	//PolicyStatusRenewed:
	policies = append(policies, policyRenewed(time.Now()))

	return policies
}

func initDevEnvForThisUser(user woptaUser) userEnv {
	ue := userEnv{}
	ue.WoptaUser = user
	// Admin
	admin := createUserWithRole(user, lib.UserRoleAdmin)
	ue.DevEnvUsers = append(ue.DevEnvUsers, admin)
	// Area manager
	arManUsr := createUserWithRole(user, lib.UserRoleManager)
	ue.DevEnvUsers = append(ue.DevEnvUsers, arManUsr)
	aMnn := createNetworkNode(arManUsr.Mail, lib.UserRoleManager, "")
	ue.NodeGraph = append(ue.NodeGraph, aMnn)
	// Agency
	agyUsr := createUserWithRole(user, lib.UserRoleAgency)
	ue.DevEnvUsers = append(ue.DevEnvUsers, agyUsr)
	aYnn := createNetworkNode(agyUsr.Mail, lib.UserRoleAgency, aMnn.Uid)
	ue.NodeGraph = append(ue.NodeGraph, aYnn)
	// Agent (online)
	agtOUsr := createUserWithRole(user, lib.UserRoleAgent+"_online")
	ue.DevEnvUsers = append(ue.DevEnvUsers, agtOUsr)
	aTOnn := createNetworkNode(agtOUsr.Mail, lib.UserRoleAgent, aYnn.Uid)
	ue.NodeGraph = append(ue.NodeGraph, aTOnn)
	// Agent (rimessa)
	agtRUsr := createUserWithRole(user, lib.UserRoleAgent+"_rimessa")
	ue.DevEnvUsers = append(ue.DevEnvUsers, agtRUsr)
	aTRnn := createNetworkNode(agtRUsr.Mail, lib.UserRoleAgent, aYnn.Uid)
	ue.NodeGraph = append(ue.NodeGraph, aTRnn)

	ue.Policies = generatePolicies()
	cstUsr := createUserWithRole(user, lib.UserRoleCustomer)
	contr := models.Contractor{Uid: uniqueID(), Name: cstUsr.Name, Surname: cstUsr.Surname, Mail: cstUsr.Mail}
	for i := 0; i < len(ue.Policies); i++ {
		ue.Policies[i].Agent = &agtOUsr
		ue.Policies[i].Contractor = contr
	}
	return ue
}

func InitDevEnv() error {
	users, err := loadUsersFromPath("test/email-list.json")
	if err != nil {
		return fmt.Errorf("input list file: %w", err)
	}
	for _, v := range users {
		ue := initDevEnvForThisUser(v)
		devEnv = append(devEnv, ue)
	}

	err = saveFileToPath("dev-env-struct.json", "/test")
	if err != nil {
		return fmt.Errorf("json output file: %w", err)
	}

	return nil
}
