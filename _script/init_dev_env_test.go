package _script

import (
	"reflect"
	"testing"
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

func TestInitDevEnv(t *testing.T) {
	err := InitDevEnv()
	if err != nil {
		t.Fatalf("expected error: %v", err)
	}
}
