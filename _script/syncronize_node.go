package _script

import (
	"errors"
	"fmt"
	"strings"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	"gitlab.dev.wopta.it/goworkspace/product"
)

func SynchronizeEleads(dryRun bool) {
	node, err := network.GetNodeByUidErr("eleads")
	if err != nil {
		panic(err)
	}
	log.PrintStruct("\nOrigin node\n", node)
	pr := product.GetLatestActiveProduct("life", models.ECommerceChannel, nil, nil)
	if err = synchronizeFieldNode(node, *pr, "quotersurvey", "companyPrivacy"); err != nil {
		panic(err)
	}
	if err = synchronizeFieldNode(node, *pr, "guaranteeconfigurationstep->privacyConsent", "consens"); err != nil {
		panic(err)
	}
	log.PrintStruct("Node changed \n\n", node)
	if dryRun {
		return
	}
	saveNode(node)
}

func ReorderEleads(dryRun bool) {
	node, err := network.GetNodeByUidErr("eleads")
	if err != nil {
		panic(err)
	}
	log.PrintStruct("\nOrigin node\n", node)
	allNameStep := getAllStepsName(getProductFromNode(node, "life"))
	log.PrintStruct("---step\n\n", allNameStep)

	if err = reorderStepsNode(node, "life", []string{"guaranteeconfigurationstep", "quotercontractordata", "quotersurvey", "quoterstatements", "quoterbeneficiary", "quoteruploaddocuments", "quoterrecap", "quotersignpay", "quoterthankyou"}); err != nil {
		panic(err)
	}
	log.PrintStruct("Node changed \n\n", node)
	if dryRun {
		return
	}
	saveNode(node)
}

func SynchronizeBeprof(dryRun bool) {
	node, err := network.GetNodeByUidErr("beprof")
	if err != nil {
		panic(err)
	}
	pr := product.GetLatestActiveProduct("life", models.ECommerceChannel, nil, nil)
	if err = synchronizeFieldNode(node, *pr, "quotersurvey", "companyPrivacy"); err != nil {
		panic(err)
	}
	if err = synchronizeFieldNode(node, *pr, "guaranteeconfigurationstep->privacyConsent", "consens"); err != nil {
		panic(err)
	}
	log.PrintStruct("Node changed \n\n", node)
	if dryRun {
		return
	}
	saveNode(node)
}

func ReorderBeProf(dryRun bool) {
	node, err := network.GetNodeByUidErr("beprof")
	if err != nil {
		panic(err)
	}
	log.PrintStruct("\nOrigin node\n", node)
	allNameStep := getAllStepsName(getProductFromNode(node, "life"))
	log.PrintStruct("---step\n\n", allNameStep)

	if err = reorderStepsNode(node, "life", []string{"guaranteeconfigurationstep", "quotercontractordata", "quotersurvey", "quoterstatements", "quoterbeneficiary", "quoteruploaddocuments", "quoterrecap", "quotersignpay", "quoterthankyou"}); err != nil {
		panic(err)
	}
	log.PrintStruct("Node changed \n\n", node)
	if dryRun {
		return
	}
	saveNode(node)
}

func SynchronizeFacile(dryRun bool) {
	node, err := network.GetNodeByUidErr("facile")
	if err != nil {
		panic(err)
	}
	log.PrintStruct("\nOrigin node\n", node)
	pr := product.GetLatestActiveProduct("life", models.ECommerceChannel, nil, nil)
	if err = synchronizeFieldNode(node, *pr, "quotersurvey", "companyPrivacy"); err != nil {
		panic(err)
	}
	if err = synchronizeFieldNode(node, *pr, "guaranteeconfigurationstep->privacyConsent", "consens"); err != nil {
		panic(err)
	}
	log.PrintStruct("Node changed \n\n", node)
	if dryRun {
		return
	}
	saveNode(node)

}
func ReorderFacile(dryRun bool) {
	node, err := network.GetNodeByUidErr("facile")
	if err != nil {
		panic(err)
	}
	log.PrintStruct("\nOrigin node\n", node)

	if err = reorderStepsNode(node, "life", []string{"guaranteeconfigurationstep", "quotercontractordata", "quotersurvey", "quoterstatements", "quoterbeneficiary", "quoteruploaddocuments", "quoterrecap", "quotersignpay", "quoterthankyou"}); err != nil {
		panic(err)
	}
	log.PrintStruct("Node changed \n\n", node)
	if dryRun {
		return
	}
	saveNode(node)

}

func saveNode(node *models.NetworkNode) {
	err := lib.SetFirestoreErr(models.NetworkNodesCollection, node.Uid, node)
	if err != nil {
		panic(err)
	}
	if err = node.SaveBigQuery(); err != nil {
		panic(err)
	}
	log.Println("\n\nNode Saved")
}
func synchronizeFieldNode(node *models.NetworkNode, product models.Product, widgetPaths string, fieldToChange string) error {
	err := synchronizeStep(node, product, widgetPaths, fieldToChange)
	return err
}

// widgetPaths name->name
func synchronizeStep(node *models.NetworkNode, productToUse models.Product, widgetPaths string, fieldToChange string) error {
	getAttributesFromChild := func(step *models.Step, widgetPathChild string) *interface{} {
		if widgetPathChild == "" {
			return &step.Attributes
		}
		for _, child := range step.Children {
			if child.Widget == widgetPathChild {
				return &child.Attributes
			}
		}
		return nil

	}
	assign := func(attributesToOverride, attributetoUse *any, nameValua string) {
		if attributetoUse == nil {
			return
		}
		if *attributesToOverride == nil {
			*attributesToOverride = *attributetoUse
			return
		}
		(*attributesToOverride).(map[string]any)[nameValua] = (*attributetoUse).(map[string]any)[nameValua]
	}

	nameWidget, remainingPaths, _ := strings.Cut(widgetPaths, "->")
	productToChange := getProductFromNode(node, productToUse.Name)
	if productToChange == nil {
		return fmt.Errorf("ProductToChange: not found")
	}
	stepToChange := getStep(nameWidget, productToChange)
	if stepToChange == nil {
		return fmt.Errorf("StepToChange: not found")
	}
	stepToUse := getStep(nameWidget, &productToUse)
	if stepToUse == nil {
		return fmt.Errorf("StepToUse: not found")
	}

	assign(getAttributesFromChild(stepToChange, remainingPaths), getAttributesFromChild(stepToUse, remainingPaths), fieldToChange)

	//widgetPaths: can be "nameWidget" or "firstWidgetName->last"
	return nil
}

func reorderStepsNode(node *models.NetworkNode, productName string, widgetsName []string) error {
	product := getProductFromNode(node, productName)
	originalSteps := product.Steps
	if len(widgetsName) > len(originalSteps) {
		return errors.New("step required too high")
	}
	var finalSteps []models.Step = make([]models.Step, len(widgetsName))
	for i := range widgetsName {
		step := getStep(widgetsName[i], product)
		if step == nil {
			return fmt.Errorf("step not found: %v", widgetsName[i])
		}
		finalSteps[i] = *step
	}
	product.Steps = finalSteps
	return nil
}

func getProductFromNode(node *models.NetworkNode, nameProduct string) *models.Product {
	for i := range node.Products {
		if node.Products[i].Name == nameProduct {
			return &node.Products[i]

		}
	}
	return nil
}

func getStep(widgetName string, product *models.Product) *models.Step {
	for i := range product.Steps {
		if widgetName == product.Steps[i].Widget {
			return &product.Steps[i]
		}
	}
	return nil
}
func getAllStepsName(product *models.Product) (names []string) {
	for i := range product.Steps {
		names = append(names, product.Steps[i].Widget)
	}
	return names
}
