package harvester

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Validate InvokeScanRequest
func (r *InvokeScanRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ResourceName, validation.Required, validation.Length(0, 200)),
		validation.Field(&r.ResourceType, validation.Required, validation.Length(0, 50)),
	)
}
