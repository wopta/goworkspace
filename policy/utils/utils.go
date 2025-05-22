package utils

import (
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
)

func CanBeAccessedBy(role, producerUid, nodeUid string) bool {
	return role == models.UserRoleAdmin || producerUid == nodeUid || network.IsParentOf(nodeUid, producerUid)
}
