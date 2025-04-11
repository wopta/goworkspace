package draftbpnm

import (
	"errors"
	"testing"

	"github.com/wopta/goworkspace/models"
)

type validity struct {
	Result bool
	Step   int
}

func (v *validity) Type() string {
	return "validity"
}

type PolicyMock struct {
	Age  int `json:"age"`
	Name string
}

func (v *PolicyMock) Type() string {
	return "policy"
}

type mockLog struct {
	log []string
}

func (m *mockLog) Println(mes string) {
	m.log = append(m.log, mes)
}

func TestBpnmHappyPath(t *testing.T) {
	g, err := NewBpnmBuilder()
	log := mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddLocal("validationObject", new(validity))
	storage.AddGlobal("policyPr", &PolicyMock{Age: 3})
	p := new(models.Policy)
	p.Name = "pippo"
	g.SetPoolDate(storage)

	g.AddHandler("init", func(st StorageData) error {
		log.Println("init")
		st.AddLocal("validationObject", new(validity))
		d, _ := st.GetGlobal("policyPr")
		if d == nil {
			return errors.New("no polscy")
		}
		return nil
	})
	g.AddHandler("AEvent", func(st StorageData) error {
		log.Println("init A")
		st.AddLocal("validationObject", new(validity))
		st.AddLocal("error", &Error{Result: false})
		return nil
	})
	g.AddHandler("BEvent", func(st StorageData) error {
		log.Println("init B")
		st.AddLocal("error", &Error{Result: false})
		st.AddLocal("validationObject", new(validity))
		return nil
	})
	g.AddHandler("CEvent", func(st StorageData) error {
		log.Println("init C")
		return nil
	})
	flow, err := g.Build()
	if err != nil {
		t.Fatal(err)
	}
	err = flow.Run("emit")
	if err != nil {
		t.Fatal(err)
	}
	exps := []string{
		"init",
		"init B",
		"init A",
	}
	if len(exps) != len(log.log) {
		t.Fatalf("exp n message: %v,got: %v", len(exps), len(log.log))
	}
	for i, exp := range exps {
		if log.log[i] != exp {
			t.Fatalf("exp: %v,got: %v", exp[i], log.log[i])
		}
	}

}
func TestBpnmHappyPath2(t *testing.T) {
	g, err := NewBpnmBuilder()
	log := mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddLocal("validationObject", new(validity))
	storage.AddGlobal("policyPr", &PolicyMock{Age: 1})
	p := new(models.Policy)
	p.Name = "pippo"
	g.SetPoolDate(storage)

	g.AddHandler("init", func(st StorageData) error {
		log.Println("init")
		st.AddLocal("validationObject", &validity{Result: true, Step: 3})
		return nil
	})
	g.AddHandler("AEvent", func(st StorageData) error {
		log.Println("init A")
		st.AddLocal("error", &Error{Result: true})
		return nil
	})
	g.AddHandler("BEvent", func(st StorageData) error {
		log.Println("init B")
		st.AddLocal("error", &Error{Result: false})
		st.AddLocal("validationObject", new(validity))
		return nil
	})
	g.AddHandler("CEvent", func(st StorageData) error {
		log.Println("init C")
		return nil
	})
	flow, err := g.Build()
	if err != nil {
		t.Fatal(err)
	}
	err = flow.Run("emit")
	if err != nil {
		t.Fatal(err)
	}
	exps := []string{
		"init",
		"init A",
		"init B",
		"init A",
	}
	if len(exps) != len(log.log) {
		t.Fatalf("exp n message: %v,got: %v", len(exps), len(log.log))
	}
	for i, exp := range exps {
		if log.log[i] != exp {
			t.Fatalf("exp: %v,got: %v", exp[i], log.log[i])
		}
	}

}
func TestBpnmMissingOutput(t *testing.T) {
	g, err := NewBpnmBuilder()
	log := mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddGlobal("policyPr", &PolicyMock{Age: 1})
	p := new(models.Policy)
	p.Name = "pippo"
	g.SetPoolDate(storage)

	g.AddHandler("init", func(st StorageData) error {
		log.Println("init")
		return nil
	})
	g.AddHandler("AEvent", func(st StorageData) error {
		log.Println("init A")
		st.AddLocal("error", &Error{Result: true})
		return nil
	})
	g.AddHandler("BEvent", func(st StorageData) error {
		log.Println("init B")
		st.AddLocal("error", &Error{Result: false})
		st.AddLocal("validationObject", new(validity))
		return nil
	})
	g.AddHandler("CEvent", func(st StorageData) error {
		log.Println("init C")
		return nil
	})
	flow, err := g.Build()
	if err != nil {
		t.Fatal(err)
	}
	err = flow.Run("emit")
	if err == nil {
		t.Fatalf("should have error")
	}
	if err.Error() != "Activity: init, output: resource required is not found validationObject" {
		t.Fatalf("should have another error, got: %v", err.Error())
	}
	if len(log.log) != 1 {
		t.Fatalf("should have 1 log")
	}
}
func TestBpnmMissingInput(t *testing.T) {
	g, err := NewBpnmBuilder()
	log := mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddGlobal("policyPr", &PolicyMock{Age: 10})
	p := new(models.Policy)
	p.Name = "pippo"
	g.SetPoolDate(storage)

	g.AddHandler("init", func(st StorageData) error {
		log.Println("init")
		st.AddLocal("validationObject", new(validity))
		return nil
	})
	g.AddHandler("AEvent", func(st StorageData) error {
		log.Println("init A")
		st.AddLocal("error", &Error{Result: true})
		return nil
	})
	g.AddHandler("BEvent", func(st StorageData) error {
		log.Println("init B")
		return nil
	})
	g.AddHandler("CEvent", func(st StorageData) error {
		log.Println("init C")
		return nil
	})
	flow, err := g.Build()
	if err != nil {
		t.Fatal(err)
	}
	err = flow.Run("emit")
	if err == nil {
		t.Fatalf("should have error")
	}
	if err.Error() != "Activity: CEvent, input: resource required is not found error" {
		t.Fatalf("should have another error, got: %v", err.Error())
	}
	if len(log.log) != 1 {
		t.Fatalf("should have 1 log")
	}
}
func TestBpnmMissingHandler(t *testing.T) {
	g, err := NewBpnmBuilder()
	log := mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddGlobal("policyPr", &PolicyMock{Age: 10})
	p := new(models.Policy)
	p.Name = "pippo"
	g.SetPoolDate(storage)

	g.AddHandler("init", func(st StorageData) error {
		log.Println("init")
		st.AddLocal("validationObject", new(validity))
		return nil
	})
	g.AddHandler("BEvent", func(st StorageData) error {
		log.Println("init B")
		return nil
	})
	g.AddHandler("CEvent", func(st StorageData) error {
		log.Println("init C")
		return nil
	})
	_, err = g.Build()
	if err.Error() != "no handler registered for the activity: AEvent" {
		t.Fatalf("should have another error, got: %v", err.Error())
	}
	if len(log.log) != 0 {
		t.Fatalf("should have 0 log")
	}
}
