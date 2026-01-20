package _script

import (
	"reflect"
	"testing"

	"github.com/wopta/goworkspace/models"
)

func TestLoadUsersFromPath(t *testing.T) {
	jsonFilePath := "test/email-list.json"
	result, err := loadUsersFromPath(jsonFilePath)
	if err != nil {
		t.Fatalf("expected error: %v", err)
	}
	expected := []woptaUser{{"Diogo", "Carvalho", "diogo@wopta.it"},
		{"Yousef", "Hammar", "yousef@wopta.it"}}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("expected %v got %v", expected, result)
	}
}

func TestAddPrefixToEmailAddress(t *testing.T) {
	address := "vinz.clortho@gozer.evil"
	result := addPrefixToEmailAddress(address, "keymaster")
	expected := "vinz.clortho+keymaster@gozer.evil"
	if result != expected {
		t.Fatalf("expected %v got %v", expected, result)
	}
}

func TestInitDevEnvForThisUser(t *testing.T) {
	wu := woptaUser{"Vinz", "Clortho", "vinz.clortho@gozer.evil"}
	u := models.User{Name: "Vinz"}
	ue := initDevEnvForThisUser(wu)
	p := ue.Policies[0]
	got := p.Agent.Name
	exp := u.Name
	if got != exp {
		t.Fatalf("expected %v got %v", exp, got)
	}
}

func TestInitDevEnv(t *testing.T) {
	err := InitDevEnv()
	if err != nil {
		t.Fatalf("expected error: %v", err)
	}
}
