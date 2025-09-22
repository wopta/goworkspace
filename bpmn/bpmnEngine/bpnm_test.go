package bpmnEngine

import (
	"errors"
	"fmt"
	"testing"
)

type errorRandomForTest struct {
	Step        string
	Description string
	Result      bool
}

func (e *errorRandomForTest) GetType() string {
	return "error"
}

type randomObjectForTest struct {
	Result bool
	Step   int
}

func (v *randomObjectForTest) GetType() string {
	return "validity"
}

type policyMock struct {
	Age  int `json:"age"`
	Name string
}

func (v *policyMock) GetType() string {
	return "policy"
}

type mockLog struct {
	log []string
}

func (m *mockLog) println(mes string) {
	m.log = append(m.log, mes)
}

func (m *mockLog) printlnForTesting(t *testing.T) {
	t.Log("Actual log: ")
	for _, mes := range m.log {
		t.Log(" ", mes)
	}
}

func testLog(log *mockLog, exps []string, t *testing.T) {
	if len(exps) != len(log.log) {
		log.printlnForTesting(t)
		t.Fatalf("exp n message: %v,got: %v", len(exps), len(log.log))
	}
	for i, exp := range exps {
		if log.log[i] != exp {
			log.printlnForTesting(t)
			t.Fatalf("exp: %v,got: %v", exp, log.log[i])
		}
	}
}

func getInjectableFlow(log *mockLog) (*BpnmBuilder, error) {
	injectedFlow, err := NewBpnmBuilderRawPath("provaInjection.json")
	if err != nil {
		return nil, err
	}
	injectedFlow.SetStorage(NewStorageBpnm())
	return injectedFlow, IsError(
		injectedFlow.AddHandler("initPost", func(st *StorageBpnm) error {
			log.println("init post")
			return nil
		}),

		injectedFlow.AddHandler("pre-B", func(st *StorageBpnm) error {
			log.println("init pre-B")
			return nil
		}),
		injectedFlow.AddHandler("initPre", func(st *StorageBpnm) error {
			log.println("init pre")
			st.AddLocal("error", &errorRandomForTest{})
			return nil
		}),
		injectedFlow.AddHandler("save", func(st *StorageBpnm) error {
			log.println("end injected process")
			return nil
		}),
	)
}
func addDefaultHandlersForTest(g *BpnmBuilder, log *mockLog) error {
	return IsError(
		g.AddHandler("init", func(st *StorageBpnm) error {
			log.println("init")
			st.AddLocal("validationObject", new(randomObjectForTest))
			return nil
		}),
		g.AddHandler("AEvent", func(st *StorageBpnm) error {
			log.println("init A")
			st.AddLocal("validationObject", new(randomObjectForTest))
			st.AddLocal("error", &errorRandomForTest{Result: false})
			return nil
		}),
		g.AddHandler("BEvent", func(st *StorageBpnm) error {
			log.println("init B")
			st.AddLocal("validationObject", new(randomObjectForTest))
			st.AddLocal("error", &errorRandomForTest{Result: false})
			return nil
		}),
		g.AddHandler("CEvent", func(st *StorageBpnm) error {
			log.println("init C")
			return nil
		}),
		g.AddHandler("DEventWithRec", func(st *StorageBpnm) error {
			log.println("init D rec")
			return nil
		}),
		g.AddHandler("DRec1", func(st *StorageBpnm) error {
			log.println("recover D1")
			return nil
		}),
		g.AddHandler("DRec2", func(st *StorageBpnm) error {
			log.println("recover D2")
			return nil
		}),
	)
}

func TestBpnmHappyPath(t *testing.T) {
	g, err := NewBpnmBuilderRawPath("prova.json")
	log := &mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddLocal("validationObject", new(randomObjectForTest))
	storage.AddGlobal("policyPr", &policyMock{Age: 3})
	g.SetStorage(storage)
	addDefaultHandlersForTest(g, log)
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
	testLog(log, exps, t)
}

func TestBpnmHappyPath2(t *testing.T) {
	g, err := NewBpnmBuilderRawPath("prova.json")
	log := &mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddLocal("validationObject", new(randomObjectForTest))
	storage.AddGlobal("policyPr", &policyMock{Age: 3})
	g.SetStorage(storage)

	addDefaultHandlersForTest(g, log)
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
	testLog(log, exps, t)
}

func TestBpnmMissingOutput(t *testing.T) {
	g, err := NewBpnmBuilderRawPath("prova.json")
	log := mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddGlobal("policyPr", &policyMock{Age: 1})
	g.SetStorage(storage)
	addDefaultHandlersForTest(g, &log)
	g.setHandler("init", func(st *StorageBpnm) error {
		log.println("init")
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
	if err.Error() != "Process 'emit' with activity 'init' has an output error: Required local resource is not found 'validationObject'" {
		t.Fatalf("should have another error, got: %v", err.Error())
	}
	if len(log.log) != 1 {
		t.Fatalf("should have 1 log")
	}
}

func TestBpnmMissingHandler(t *testing.T) {
	g, err := NewBpnmBuilderRawPath("prova.json")
	log := mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddGlobal("policyPr", &policyMock{Age: 10})
	g.SetStorage(storage)
	addDefaultHandlersForTest(g, &log)
	g.setHandler("AEvent", nil)
	_, err = g.Build()
	if err.Error() != "No handler registered for the activity: 'AEvent'" {
		t.Fatalf("should have another error, got: %v", err.Error())
	}
	if len(log.log) != 0 {
		t.Fatalf("should have 0 log")
	}
}

func TestBpnmInjection(t *testing.T) {
	g, err := NewBpnmBuilderRawPath("prova.json")
	log := &mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddGlobal("policyPr", &policyMock{Age: 3})
	err = addDefaultHandlersForTest(g, log)
	if err != nil {
		t.Fatal(err)
	}
	g.SetStorage(storage)
	flowCatnat, err := getInjectableFlow(log)
	if err := g.Inject(flowCatnat); err != nil {
		t.Fatal(err)
	}

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
		"end injected process",
		"init A",
		"init post",
		"end injected process",
	}
	testLog(log, exps, t)
}

func TestBpnmWithMultipleInjection(t *testing.T) {
	g, err := NewBpnmBuilderRawPath("prova.json")
	log := mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddGlobal("policyPr", &policyMock{Age: 3})
	g.SetStorage(storage)
	flowCatnat, err := getInjectableFlow(&log)
	if err := g.Inject(flowCatnat); err != nil {
		t.Fatal(err)
	}
	err = g.Inject(flowCatnat)
	if err == nil {
		t.Fatalf("Should have an error")
	}

	if err.Error() != "Injection's been already done for: target process: 'emit', process: injected 'provaPost' with order 'Post'" {
		t.Fatalf("Injection's been already done for: target process: 'emit', process: injected 'provaPost' with order Post, got: %v'", err)
	}
}

func TestRunFromSpecificActivity(t *testing.T) {
	g, err := NewBpnmBuilderRawPath("prova.json")
	log := &mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddGlobal("policyPr", &policyMock{Age: 3})
	storage.AddLocal("validationObject", new(randomObjectForTest))
	g.SetStorage(storage)
	flowCatnat, err := getInjectableFlow(log)
	if err := g.Inject(flowCatnat); err != nil {
		t.Fatal(err)
	}

	addDefaultHandlersForTest(g, log)
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
		"end injected process",
	}
	testLog(log, exps, t)
}

func TestBpnmStoreClean(t *testing.T) {
	//this case test how the framework manage memory
	//at each cycles
	//it marks every output resource of each activities (T), after all activities(T) have finished, it clean the store leaving only the marked ones
	g, err := NewBpnmBuilderRawPath("prova.json")
	log := &mockLog{}
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddLocal("validationObject", new(randomObjectForTest))
	storage.AddGlobal("policyPr", &policyMock{Age: 2})
	g.SetStorage(storage)
	addDefaultHandlersForTest(g, log)

	g.setHandler("init", func(st *StorageBpnm) error {
		log.println("init")
		st.AddLocal("validationObject", new(randomObjectForTest))
		st.AddLocal("error", &errorRandomForTest{Result: false})
		st.AddLocal("error1", &errorRandomForTest{Result: false})
		st.AddLocal("error2", &errorRandomForTest{Result: false})
		st.AddLocal("error3", &errorRandomForTest{Result: false})
		_, e := st.GetGlobal("policyPr")
		if e != nil {
			return e
		}
		if len(st.getAllLocals()) != 5 { //output of init
			return fmt.Errorf("store hasnt the right number of resources %v", len(st.getAllLocals()))
		}
		return nil
	})
	g.setHandler("AEvent", func(st *StorageBpnm) error {
		log.println("init A")
		st.AddLocal("error", &errorRandomForTest{Result: false})
		d, e := GetData[*randomObjectForTest]("validationObject", st)
		if e != nil {
			return e
		}
		d.Step = 3
		p, e := GetData[*errorRandomForTest]("error", st)
		if e != nil {
			return e
		}
		p.Result = true
		return nil
	})
	g.setHandler("BEvent", func(st *StorageBpnm) error {
		log.println("init B")
		if len(st.getAllLocals()) != 2 { //output of AEvent
			return fmt.Errorf("Expected 2 resource from AEvent, got: %v", len(st.getAllLocals()))
		}
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
	testLog(log, exps, t)
}

func TestMergeBuilder(t *testing.T) {
	log := &mockLog{}
	g, err := NewBpnmBuilderRawPath("prova.json")
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	addDefaultHandlersForTest(g, log)
	g.AddHandler("end", func(sd *StorageBpnm) error {
		log.println("end")
		return nil
	})
	storage.AddLocal("validationObject", new(randomObjectForTest))
	storage.AddGlobal("policyPr", &policyMock{Age: 2})
	g.SetStorage(storage)
	b2, err := getInjectableFlow(log)
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
	testLog(log, exps, t)
}

func TestErrorWithoutRecover(t *testing.T) {
	log := &mockLog{}
	g, err := NewBpnmBuilderRawPath("prova.json")
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	addDefaultHandlersForTest(g, log)
	g.setHandler("AEvent", func(sd *StorageBpnm) error {
		return errors.New("error")
	})
	storage.AddLocal("validationObject", new(randomObjectForTest))
	storage.AddGlobal("policyPr", &policyMock{Age: 2})
	g.SetStorage(storage)
	f, err := g.Build()
	f.process["emit"].recover = nil
	if err != nil {
		t.Fatal(err)
	}
	err = f.Run("emit")
	if err == nil {
		t.Fatalf("should have error")
	}
	if err.Error() != "error" {
		t.Fatalf("should have another error, got: %v", err.Error())
	}
}

func TestRecoverWithFunction(t *testing.T) {
	log := &mockLog{}
	g, err := NewBpnmBuilderRawPath("prova.json")
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	storage.AddLocal("validationObject", new(randomObjectForTest))
	storage.AddGlobal("policyPr", &policyMock{Age: 2})
	addDefaultHandlersForTest(g, log)
	g.setHandler("DEventWithRec", func(st *StorageBpnm) error {
		log.println("init D")
		return errors.New("error of DEvent")
	})
	g.SetStorage(storage)
	f, err := g.Build()
	if err != nil {
		t.Fatal(err)
	}
	err = f.RunAt("emit", "DEventWithRec")
	if err == nil {
		t.Fatal("Should have an error")
	}
	exps := []string{
		"init D",
		"recover D1",
		"recover D2",
	}
	testLog(log, exps, t)
}

func TestRecoverFromPanic(t *testing.T) {
	log := &mockLog{}
	g, err := NewBpnmBuilderRawPath("prova.json")
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	addDefaultHandlersForTest(g, log)
	g.setHandler("DEventWithRec", func(st *StorageBpnm) error {
		log.println("init D rec")
		panic("fjsdklfjd")
	})
	storage.AddLocal("validationObject", new(randomObjectForTest))
	storage.AddGlobal("policyPr", &policyMock{Age: 2})
	g.SetStorage(storage)
	f, err := g.Build()
	if err != nil {
		t.Fatal(err)
	}
	err = f.RunAt("emit", "DEventWithRec")
	if err == nil {
		t.Fatal("Should have an error")
	}
	exps := []string{
		"init D rec",
		"recover D1",
		"recover D2",
	}
	testLog(log, exps, t)
}
func TestEndActivity(t *testing.T) {
	log := &mockLog{}
	g, err := NewBpnmBuilderRawPath("prova.json")
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	addDefaultHandlersForTest(g, log)
	storage.AddLocal("validationObject", new(randomObjectForTest))
	storage.AddGlobal("policyPr", &policyMock{Age: 2})
	g.SetStorage(storage)
	g.AddHandler("end_emit", func(sd *StorageBpnm) error {
		log.println("end")
		return nil
	})
	f, err := g.Build()
	if err != nil {
		t.Fatal(err)
	}
	err = f.RunAt("emit", "init")
	if err != nil {
		t.Fatal(err)
	}
	exps := []string{
		"init",
		"init A",
		"end",
	}
	testLog(log, exps, t)
}

func TestDontCallEndAfterInit(t *testing.T) {
	//i've set "callEndIfStop": false,
	log := &mockLog{}
	g, err := NewBpnmBuilderRawPath("prova.json")
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	addDefaultHandlersForTest(g, log)
	storage.AddLocal("validationObject", &randomObjectForTest{})
	storage.AddGlobal("policyPr", &policyMock{Age: 20})
	g.SetStorage(storage)
	g.AddHandler("end_emit", func(sd *StorageBpnm) error {
		log.println("end")
		return nil
	})
	f, err := g.Build()
	if err != nil {
		t.Fatal(err)
	}
	err = f.RunAt("emit", "init")
	if err != nil {
		t.Fatal(err)
	}
	exps := []string{
		"init",
	}
	testLog(log, exps, t)
}

func TestHandlerLessTrue(t *testing.T) {
	//i've set "handlerless": true,
	log := &mockLog{}
	g, err := NewBpnmBuilderRawPath("prova.json")
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorageBpnm()
	addDefaultHandlersForTest(g, log)
	storage.AddLocal("validationObject", &randomObjectForTest{Step: 3})
	storage.AddGlobal("policyPr", &policyMock{Age: 2})
	g.setHandler("AEvent", func(sd *StorageBpnm) error {
		sd.AddLocal("error", &errorRandomForTest{Result: true})
		log.println("init A")
		return nil
	})
	g.setHandler("BEvent", nil) //remove the BEvent's handler
	g.SetStorage(storage)
	f, err := g.Build()
	if err != nil {
		t.Fatal(err)
	}
	err = f.RunAt("emit", "init")
	if err != nil {
		t.Fatal(err)
	}
	exps := []string{
		"init",
		"init A",
	}
	testLog(log, exps, t)
}
