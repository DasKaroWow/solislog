package solislog

// Extra stores contextual key-value fields attached to log records.
//
// Extra values can be rendered with template placeholders such as
// {extra[user_id]} or as a full JSON object with {extra}.
type Extra map[string]string
