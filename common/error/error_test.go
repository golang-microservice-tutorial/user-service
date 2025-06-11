package error

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type userInput struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"min=18"`
	Code  string `validate:"customtag"` // Untuk test unknown tag
}

func TestErrorValidationResponse(t *testing.T) {
	validate := validator.New()

	err := validate.RegisterValidation("customtag", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == "ok"
	})
	if err != nil {
		t.Fatalf("gagal register custom validation tag: %v", err)
	}

	ErrValidator = map[string]string{
		"min":      "%s must be at least %s",
		"required": "%s is required",
		"email":    "%s is not valid email",
	}

	tests := []struct {
		name     string
		input    userInput
		expected []ValidationResponse
	}{
		{
			name:     "All valid",
			input:    userInput{Name: "John", Email: "john@example.com", Age: 20, Code: "ok"},
			expected: nil,
		},
		{
			name:  "Age below minimum",
			input: userInput{Name: "John", Email: "john@example.com", Age: 16, Code: "ok"},
			expected: []ValidationResponse{
				{Field: "Age", Message: "Age must be at least 18"},
			},
		},
		{
			name:  "required name",
			input: userInput{Name: "", Email: "jhon@example.com", Age: 20, Code: "ok"},
			expected: []ValidationResponse{
				{Field: "Name", Message: "field Name is required"},
			},
		},
		{
			name:  "Invalid email format",
			input: userInput{Name: "John", Email: "invalid-email", Age: 20, Code: "ok"},
			expected: []ValidationResponse{
				{Field: "Email", Message: "field Email is not valid email"},
			},
		},
		{
			name:  "Unknown tag (custom)",
			input: userInput{Name: "John", Email: "john@example.com", Age: 20, Code: "fail"},
			expected: []ValidationResponse{
				{Field: "Code", Message: "Something wrong on Code; customtag"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validate.Struct(tc.input)

			if err == nil {
				assert.Nil(t, err)
				return
			}
			// pastiin error emang dari validator
			var validatorErrs validator.ValidationErrors
			if errors.As(err, &validatorErrs) {
				result := ErrValidationResponse(err)
				assert.ElementsMatch(t, tc.expected, result)
			} else {
				t.Fatalf("error bukan dari validator, got: %T", err)
			}
		})
	}
}
