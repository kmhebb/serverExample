package errors

type InvalidDetails struct {
	Field string
	Value interface{}
}

func (d InvalidDetails) Map() map[string]interface{} {
	return map[string]interface{}{
		"field": d.Field,
		"value": d.Value,
	}
}

func Invalid(field string, value interface{}) *Error {
	return E(EINVALID).WithExtra(InvalidDetails{
		Field: field,
		Value: value,
	}).WithMessage(
		"You provided an invalid %s: %v.",
		field, value,
	).WithErrorf(
		"invalid %s: %v",
		field, value,
	)
}
