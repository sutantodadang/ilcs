package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ValidationError struct {
	Namespace       string `json:"namespace"`
	Field           string `json:"field"`
	StructNamespace string `json:"structNamespace"`
	StructField     string `json:"structField"`
	Tag             string `json:"tag"`
	ActualTag       string `json:"actualTag"`
	Kind            string `json:"kind"`
	Type            string `json:"type"`
	Value           string `json:"value"`
	Param           string `json:"param"`
	Message         string `json:"message"`
}

func (v *ValidationError) Error() string {
	return v.Message
}

func NewValidationError(err error) []*ValidationError {

	var errs []*ValidationError

	for _, v := range err.(validator.ValidationErrors) {
		errs = append(errs, &ValidationError{
			Namespace:       v.Namespace(),
			Field:           v.Field(),
			StructNamespace: v.StructNamespace(),
			StructField:     v.StructField(),
			Tag:             v.Tag(),
			ActualTag:       v.ActualTag(),
			Kind:            v.Kind().String(),
			Type:            v.Type().String(),
			Value:           fmt.Sprintf("%v", v.Value()),
			Param:           v.Param(),
			Message:         v.Error(),
		})

	}

	return errs

}

func ValidateId(id string) error {

	_, err := uuid.Parse(id)

	if err != nil {
		return fmt.Errorf("invalid id")
	}

	return nil
}
