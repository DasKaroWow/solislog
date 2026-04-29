package solislog

import (
	"maps"
)

// Extra stores contextual key-value fields attached to log records.
//
// Extra values can be rendered with template placeholders such as
// {extra[user_id]} or as a full JSON object with {extra}.
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
