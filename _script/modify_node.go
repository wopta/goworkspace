package _script

import (
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
)

func ModifyNode() error {
	// return modifyNodeFacile("ptns_remunerazione_zero")
	return nil
}

func modifyNodeFacile(warrant string) error {
	node, err := network.GetNodeByUidErr("facile")
	if err != nil {
		return err
	}

	node.Warrant = warrant

	err = lib.SetFirestoreErr(models.NetworkNodesCollection, node.Uid, node)
	if err != nil {
		return err
	}
	return node.SaveBigQuery("")
}
