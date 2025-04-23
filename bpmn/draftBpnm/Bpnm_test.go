package draftbpnm

import (
	"fmt"
	"testing"
)

type validity struct {
	Result bool
	Step   int
}

func (v *validity) GetType() string {
	return "validity"
}

type PolicyMock struct {
	Age  int `json:"age"`
	Name string
}

func (v *PolicyMock) GetType() string {
	return "policy"
}

type mockLog struct {
	log []string
}

func (m *mockLog) Println(mes string) {
	m.log = append(m.log, mes)
}

func TestBpnmHappyPath(t *testing.T) {
	g, err := NewBpnmBuilder("prova.json")
	log := mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddLocal("validationObject", new(validity))
	storage.AddGlobal("policyPr", &PolicyMock{Age: 3})
	g.SetStorage(storage)

	g.AddHandler("init", func(st StorageData) error {
		log.Println("init")
		st.AddLocal("validationObject", new(validity))
		st.AddLocal("garbage1", new(validity))
		st.AddLocal("garbage2", new(validity))
		st.AddLocal("garbage3", new(validity))
		st.AddLocal("error", &Error{Result: false})
		_, e := st.GetGlobal("policyPr")
		if e != nil {
			return e
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
		st.AddLocal("validationObject", new(validity))
		st.AddLocal("error", &Error{Result: false})
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
	err = flow.RunAt("emit", "init")
	if err != nil {
		t.Fatal(err)
	}
	exps := []string{
		"init",
		"init B",
		"init A",
		"init A",
	}
	if len(exps) != len(log.log) {
		t.Fatalf("exp n message: %v,got: %v", len(exps), len(log.log))
	}
	for i, exp := range exps {
		if log.log[i] != exp {
			t.Fatalf("exp: %v,got: %v", exp, log.log[i])
		}
	}

}

func TestBpnmHappyPath2(t *testing.T) {
	g, err := NewBpnmBuilder("prova.json")
	log := mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddLocal("validationObject", new(validity))
	storage.AddGlobal("policyPr", &PolicyMock{Age: 3})
	g.SetStorage(storage)

	g.AddHandler("init", func(st StorageData) error {
		log.Println("init")
		st.AddLocal("validationObject", &validity{Result: true, Step: 3})
		st.AddLocal("garbage1", new(validity))
		st.AddLocal("garbage2", new(validity))
		st.AddLocal("garbage3", new(validity))
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
	err = flow.RunAt("emit", "init")
	if err != nil {
		t.Fatal(err)
	}
	exps := []string{
		"init",
		"init B",
		"init A",
		"init A",
	}
	if len(exps) != len(log.log) {
		t.Fatalf("exp n message: %v,got: %v", len(exps), len(log.log))
	}
	for i, exp := range exps {
		if log.log[i] != exp {
			t.Fatalf("exp: %v,got: %v", exp, log.log[i])
		}
	}

}

func TestBpnmMissingOutput(t *testing.T) {
	g, err := NewBpnmBuilder("prova.json")
	log := mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddGlobal("policyPr", &PolicyMock{Age: 1})
	g.SetStorage(storage)

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
	err = flow.RunAt("emit", "init")
	if err == nil {
		t.Fatalf("should have error")
	}
	if err.Error() != "Process 'emit' with activity 'init' has error: Required local resource is not found 'validationObject'" {
		t.Fatalf("should have another error, got: %v", err.Error())
	}
	if len(log.log) != 1 {
		t.Fatalf("should have 1 log")
	}
}

func TestBpnmMissingInput(t *testing.T) {
	g, err := NewBpnmBuilder("prova.json")
	log := mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddGlobal("policyPr", &PolicyMock{Age: 10})
	g.SetStorage(storage)

	g.AddHandler("init", func(st StorageData) error {
		log.Println("init")
		st.AddLocal("validationObject", new(validity))
		st.AddLocal("garbage1", new(validity))
		st.AddLocal("garbage2", new(validity))
		st.AddLocal("garbage3", new(validity))
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
	err = flow.RunAt("emit", "init")
	if err == nil {
		t.Fatalf("should have error")
	}
	if err.Error() != "Process 'emit' with activity 'CEvent' has an input error: Required local resource is not found 'error'" {
		t.Fatalf("should have another error, got: %v", err.Error())
	}
	if len(log.log) != 1 {
		t.Fatalf("should have 1 log")
	}
}

func TestBpnmMissingHandler(t *testing.T) {
	g, err := NewBpnmBuilder("prova.json")
	log := mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddGlobal("policyPr", &PolicyMock{Age: 10})
	g.SetStorage(storage)

	g.AddHandler("init", func(st StorageData) error {
		log.Println("init")
		st.AddLocal("validationObject", new(validity))
		st.AddLocal("garbage1", new(validity))
		st.AddLocal("garbage2", new(validity))
		st.AddLocal("garbage3", new(validity))
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
	if err.Error() != "No handler registered for the activity: AEvent" {
		t.Fatalf("should have another error, got: %v", err.Error())
	}
	if len(log.log) != 0 {
		t.Fatalf("should have 0 log")
	}
}

func getFlowcatnat(log *mockLog) (*BpnmBuilder, error) {
	injectedFlow, err := NewBpnmBuilder("provaInjection.json")
	injectedFlow.SetStorage(NewStorageBpnm())
	injectedFlow.AddHandler("initPost", func(st StorageData) error {
		log.Println("init post")
		return nil
	})

	injectedFlow.AddHandler("pre-B", func(st StorageData) error {
		log.Println("init pre-B")
		return nil
	})
	injectedFlow.AddHandler("initPre", func(st StorageData) error {
		log.Println("init pre")
		return nil
	})
	return injectedFlow, err
}

func TestBpnmInjection(t *testing.T) {
	g, err := NewBpnmBuilder("prova.json")
	log := mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddGlobal("policyPr", &PolicyMock{Age: 3})
	g.SetStorage(storage)
	flowCatnat, err := getFlowcatnat(&log)
	if err := g.Inject(flowCatnat); err != nil {
		t.Fatal(err)
	}

	g.AddHandler("init", func(st StorageData) error {
		log.Println("init")
		st.AddLocal("validationObject", new(validity))
		st.AddLocal("garbage1", new(validity))
		st.AddLocal("garbage2", new(validity))
		st.AddLocal("garbage3", new(validity))
		_, e := st.GetGlobal("policyPr")
		if e != nil {
			return e
		}
		return nil
	})
	g.AddHandler("AEvent", func(st StorageData) error {
		log.Println("init A")
		st.AddLocal("validationObject", new(validity))
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
		"init pre",
		"init pre-B",
		"init B",
		"init A",
		"init post",
		"init A",
		"init post",
	}
	for i, exp := range exps {
		if log.log[i] != exp {
			t.Fatalf("exp: %v,got: %v", exp, log.log[i])
		}
	}
}

func TestBpnmWithMultipleInjection(t *testing.T) {
	g, err := NewBpnmBuilder("prova.json")
	log := mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddGlobal("policyPr", &PolicyMock{Age: 3})
	g.SetStorage(storage)
	flowCatnat, err := getFlowcatnat(&log)
	if err := g.Inject(flowCatnat); err != nil {
		t.Fatal(err)
	}
	err = g.Inject(flowCatnat)
	if err == nil {
		t.Fatalf("Should have an error")
	}

	if err.Error() != "Injection's been already done: target process: emit, process: injected provaPost" {
		t.Fatalf("Should have the error,exp: Injection's been already done: target process: emit, process: injected provaPost, got: %v", err.Error())
	}
}

func TestRunFromSpecificActivity(t *testing.T) {
	g, err := NewBpnmBuilder("prova.json")
	log := mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddGlobal("policyPr", &PolicyMock{Age: 3})
	storage.AddLocal("validationObject", new(validity))
	g.SetStorage(storage)
	flowCatnat, err := getFlowcatnat(&log)
	if err := g.Inject(flowCatnat); err != nil {
		t.Fatal(err)
	}

	g.AddHandler("init", func(st StorageData) error {
		log.Println("init")
		st.AddLocal("validationObject", new(validity))
		st.AddLocal("garbage1", new(validity))
		st.AddLocal("garbage2", new(validity))
		st.AddLocal("garbage3", new(validity))
		_, e := st.GetGlobal("policyPr")
		if e != nil {
			return e
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
	err = flow.RunAt("emit", "BEvent")
	if err != nil {
		t.Fatal(err)
	}

	exps := []string{
		"init pre",
		"init pre-B",
		"init B",
		"init A",
		"init post",
	}
	if len(exps) != len(log.log) {
		t.Fatalf("exp n message: %v,got: %v", len(exps), len(log.log))
	}
	for i, exp := range exps {
		if log.log[i] != exp {
			t.Fatalf("exp: %v,got: %v", exp, log.log[i])
		}
	}
}

func TestBpnmStoreClean(t *testing.T) {
	//this case test how the framework manage memory
	//at each cycles
	//it marks every output resource of each activities (T), after all activities(T) have finished, it clean the store leaving only the marked ones
	g, err := NewBpnmBuilder("prova.json")
	log := mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddLocal("validationObject", new(validity))
	storage.AddGlobal("policyPr", &PolicyMock{Age: 2})
	g.SetStorage(storage)

	g.AddHandler("init", func(st StorageData) error {
		log.Println("init")
		st.AddLocal("validationObject", new(validity))
		st.AddLocal("garbage1", new(validity))
		st.AddLocal("garbage2", new(validity))
		st.AddLocal("garbage3", new(validity))
		st.AddLocal("error", &Error{Result: false})
		_, e := st.GetGlobal("policyPr")
		if e != nil {
			return e
		}
		if len(st.getAllLocal()) != 5 { //output of init
			return fmt.Errorf("store hasn't been cleaned right, n resource %v", len(st.getAllLocal()))
		}
		return nil
	})
	g.AddHandler("AEvent", func(st StorageData) error {
		log.Println("init A")
		st.AddLocal("error", &Error{Result: false})
		d, e := GetData[*validity]("validationObject", st)
		if e != nil {
			return e
		}
		d.Step = 3
		p, e := GetData[*Error]("error", st)
		if e != nil {
			return e
		}
		p.Result = true
		return nil
	})
	g.AddHandler("BEvent", func(st StorageData) error {
		log.Println("init B")
		if len(st.getAllLocal()) != 2 { //output of AEvent
			return fmt.Errorf("Expected 2 resource from AEvent, got: %v", len(st.getAllLocal()))
		}
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
	err = flow.RunAt("emit", "init")
	if err != nil {
		t.Fatal(err)
	}
	exps := []string{
		"init",
		"init A",
		"init B",
	}
	if len(exps) != len(log.log) {
		t.Fatalf("exp n message: %v,got: %v", len(exps), len(log.log))
	}
	for i, exp := range exps {
		if log.log[i] != exp {
			t.Fatalf("exp: %v,got: %v", exp, log.log[i])
		}
	}

}

func TestMergeBuilder(t *testing.T) {
	log := &mockLog{}
	g, err := NewBpnmBuilder("prova.json")
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	g.AddHandler("init", func(st StorageData) error {
		log.Println("init")
		st.AddLocal("validationObject", new(validity))
		st.AddLocal("garbage1", new(validity))
		st.AddLocal("garbage2", new(validity))
		st.AddLocal("garbage3", new(validity))
		st.AddLocal("error", &Error{Result: false})
		_, e := st.GetGlobal("policyPr")
		if e != nil {
			return e
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
		st.AddLocal("validationObject", new(validity))
		st.AddLocal("error", &Error{Result: false})
		return nil
	})
	g.AddHandler("CEvent", func(st StorageData) error {
		log.Println("init C")
		return nil
	})
	storage.AddLocal("validationObject", new(validity))
	storage.AddGlobal("policyPr", &PolicyMock{Age: 2})
	g.SetStorage(storage)
	b2, err := getFlowcatnat(log)
	if err != nil {
		t.Fatal(err)
	}
	err = g.AddProcesses(b2)
	if err != nil {
		t.Fatal(err)
	}
	f, err := g.Build()
	if err != nil {
		t.Fatal(err)
	}
	err = f.Run("provaPre")
	if err != nil {
		t.Fatal(err)
	}
	exps := []string{
		"init pre",
		"init pre-B",
	}
	if len(exps) != len(log.log) {
		t.Fatalf("exp n message: %v,got: %v", len(exps), len(log.log))
	}
	for i, exp := range exps {
		if log.log[i] != exp {
			t.Fatalf("exp: %v,got: %v", exp, log.log[i])
		}
	}
}
