package solislog

import (
	"maps"
)

type Extra map[string]string

func cloneExtra(src Extra) Extra {
	if src == nil {
		return Extra{}
	}
	return maps.Clone(src)
}

func mergeExtra(base Extra, override Extra) Extra {
	if base == nil {
		base = Extra{}
	}
	if override == nil {
		override = Extra{}
	}

	merged := maps.Clone(base)
	maps.Insert(merged, maps.All(override))

	return merged
}
