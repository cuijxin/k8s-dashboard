package clusterrolebinding

import "github.com/cuijxin/k8s-dashboard/src/backend/resource/dataselect"

// The code below allows to perform complex data section on []ClusterRoleBinding
type ClusterRoleBindingCell ClusterRoleBinding

func (self ClusterRoleBindingCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
	switch name {
	case dataselect.NameProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Name)
	case dataselect.CreationTimestampProperty:
		return dataselect.StdComparableTime(self.ObjectMeta.CreationTimestamp.Time)
	case dataselect.NamespaceProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Namespace)
	default:
		// if name is not supported then just return a constant dummy value, sort will haveno effect.
		return nil
	}
}

func toCells(std []ClusterRoleBinding) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = ClusterRoleBindingCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []ClusterRoleBinding {
	std := make([]ClusterRoleBinding, len(cells))
	for i := range std {
		std[i] = ClusterRoleBinding(cells[i].(ClusterRoleBindingCell))
	}
	return std
}
