package utils

import (
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
)

func CanBeAccessedBy(role, producerUid, nodeUid string) bool {
	return role == models.UserRoleAdmin || producerUid == nodeUid || network.IsParentOf(nodeUid, producerUid)
}
