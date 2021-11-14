package utils

type UpsertQuery struct {
	Query      string
	Parameters []interface{}
	Table      string
}
