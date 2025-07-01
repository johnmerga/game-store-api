package validator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (ve ValidationErrors) Error() string {
	return fmt.Sprintf("validation failed with %d errors", len(ve.Errors))
}

func New() *Validator {
	validate := validator.New()

	// Register custom tag name function to use json tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{validate: validate}
}

func (v *Validator) ValidateStruct(s interface{}) error {
	if err := v.validate.Struct(s); err != nil {
		return v.formatValidationErrors(err)
	}
	return nil
}

func (v *Validator) ValidateAndParseJSON(r *http.Request, s interface{}) error {
	// Parse JSON body
	if err := json.NewDecoder(r.Body).Decode(s); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Validate struct
	return v.ValidateStruct(s)
}

func (v *Validator) formatValidationErrors(err error) ValidationErrors {
	var validationErrors []ValidationError

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrs {
			validationErrors = append(validationErrors, ValidationError{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Value:   fmt.Sprintf("%v", err.Value()),
				Message: v.getErrorMessage(err),
			})
		}
	}

	return ValidationErrors{Errors: validationErrors}
}

func (v *Validator) getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", err.Field())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", err.Field())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}
