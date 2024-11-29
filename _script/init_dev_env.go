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

func generatePoliciesForThisUser(user woptaUser) []models.Policy {
	/*
		REMINDER for myself
		- CustomerUser must be created on the fly with fake data and then attached to policies

		TODO
			Per ogni agency/agent creati al punto 1 procedere con creazione:
			•	Polizza pagata
			•	Polizza eliminata
			•	Polizza da pagare (emessa attualmente)
			•	Polizza da pagare (emessa X mesi precedenti)  caso insoluto
			•	Polizza da firmare
			•	Proposta da emettere
			•	Proposta da emettere (con startDate 30 giorni nel passato)
			•	Proposta in riservato
			•	Proposta approvata
			•	Proposta rifiutata
			•	Quietanza (polizza in draft renew) da pagare con mandato
			•	Quietanza da pagare senza mandato
			•	Quietanza eliminata
			•	Quietanza pagata
	*/
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
	//// Customer
	//cstUsr := createUserWithRole(user, lib.UserRoleCustomer)
	//ue.DevEnvUsers = append(ue.DevEnvUsers, cstUsr)

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
