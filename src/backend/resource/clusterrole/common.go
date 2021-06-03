package clusterrole

import (
	"github.com/cuijxin/k8s-dashboard/src/backend/resource/dataselect"
)

// The code below allows to perform complex data section on []ClusterRole

type RoleCell ClusterRole

func (self RoleCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
	switch name {
	case dataselect.NameProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Name)
	case dataselect.CreationTimestampProperty:
		return dataselect.StdComparableTime(self.ObjectMeta.CreationTimestamp.Time)
	case dataselect.NamespaceProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Namespace)
	default:
		// if name is not supported then just return a constant dummy value, sort will have no effect.
		return nil
	}
}

func toCells(std []ClusterRole) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = RoleCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []ClusterRole {
	std := make([]ClusterRole, len(cells))
	for i := range std {
		std[i] = ClusterRole(cells[i].(RoleCell))
	}
	return std
}
