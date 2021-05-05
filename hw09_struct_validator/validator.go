package hw09structvalidator

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const tagName = "validate"

var (
	ErrValueIsNotStruct          = errors.New("value is not a structure")
	ErrFieldTypeNotSupported     = errors.New("field type is not supported")
	ErrValidatorTypeNotSupported = errors.New("validator type is not supported")
	ErrValidatorWrongSettings    = errors.New("validator has wrong settings")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	b := bytes.NewBuffer(make([]byte, 0))

	for _, err := range v {
		fmt.Fprintf(b, "%s: %s, ", err.Field, err.Err.Error())
	}

	return strings.TrimSuffix(b.String(), ", ")
}

func strsToInts(strs []string) ([]int, error) {
	ints := make([]int, len(strs))

	for i, str := range strs {
		val, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}

		ints[i] = val
	}

	return ints, nil
}

func wrapValidateError(field, validator string, err error) error {
	return fmt.Errorf("unable to validate '%s' field with validator '%s': %w", field, validator, err)
}

func validateStrField(name, value string, rules []string) (ValidationErrors, error) {
	var errs ValidationErrors

	for _, rule := range rules {
		parts := strings.SplitN(rule, ":", 2)
		if len(parts) != 2 {
			return nil, wrapValidateError(name, rule, ErrValidatorWrongSettings)
		}

		validator, arg := parts[0], parts[1]

		switch validator {
		case "len":
			if ln, err := strconv.Atoi(arg); err != nil {
				return nil, wrapValidateError(name, validator, err)
			} else if err := validateStrLen(value, ln); err != nil {
				errs = append(errs, ValidationError{name, err})
			}
		case "in":
			values := strings.Split(arg, ",")

			if err := validateStrIn(value, values); err != nil {
				errs = append(errs, ValidationError{name, err})
			}
		case "regexp":
			if exp, err := regexp.Compile(arg); err != nil {
				return nil, wrapValidateError(name, validator, err)
			} else if err := validateStrRegExp(value, exp); err != nil {
				errs = append(errs, ValidationError{name, err})
			}
		default:
			return nil, wrapValidateError(name, validator, ErrValidatorTypeNotSupported)
		}
	}

	return errs, nil
}

func validateIntField(name string, value int64, rules []string) (ValidationErrors, error) {
	var errs ValidationErrors

	for _, rule := range rules {
		parts := strings.SplitN(rule, ":", 2)
		if len(parts) != 2 {
			return nil, wrapValidateError(name, rule, ErrValidatorWrongSettings)
		}

		validator, arg := parts[0], parts[1]

		switch validator {
		case "min":
			if min, err := strconv.Atoi(arg); err != nil {
				return nil, wrapValidateError(name, validator, err)
			} else if err := validateIntMin(value, min); err != nil {
				errs = append(errs, ValidationError{name, err})
			}
		case "max":
			if min, err := strconv.Atoi(arg); err != nil {
				return nil, wrapValidateError(name, validator, err)
			} else if err := validateIntMax(value, min); err != nil {
				errs = append(errs, ValidationError{name, err})
			}
		case "in":
			if values, err := strsToInts(strings.Split(arg, ",")); err != nil {
				return nil, wrapValidateError(name, validator, err)
			} else if err := validateIntIn(value, values); err != nil {
				errs = append(errs, ValidationError{name, err})
			}
		default:
			return nil, wrapValidateError(name, validator, ErrValidatorTypeNotSupported)
		}
	}

	return errs, nil
}

func validateStructField(name string, value reflect.Value, rules []string) (ValidationErrors, error) {
	var errs ValidationErrors

	for _, rule := range rules {
		switch rule {
		case "nested":
			var vErrs ValidationErrors

			if err := Validate(value.Interface()); errors.As(err, &vErrs) {
				errs = append(errs, ValidationError{name, vErrs})
			} else if err != nil {
				return nil, wrapValidateError(name, rule, err)
			}
		default:
			return nil, wrapValidateError(name, rule, ErrValidatorTypeNotSupported)
		}
	}

	return errs, nil
}

func validateField(name string, t reflect.Type, v reflect.Value, rules []string) (ValidationErrors, error) {
	switch t.Kind() {
	case reflect.Struct:
		return validateStructField(name, v, rules)

	case reflect.Slice:
		var errs ValidationErrors

		for i := 0; i < v.Len(); i++ {
			k := fmt.Sprintf("%s_%d", name, i)

			vErrs, err := validateField(k, t.Elem(), v.Index(i), rules)
			if err != nil {
				return nil, err
			}

			errs = append(errs, vErrs...)
		}

		return errs, nil

	case reflect.Int:
		return validateIntField(name, v.Int(), rules)

	case reflect.String:
		return validateStrField(name, v.String(), rules)

	default:
		return nil, fmt.Errorf("unable to validate '%s' field: %w", name, ErrFieldTypeNotSupported)
	}
}

func Validate(v interface{}) error {
	var errs ValidationErrors

	rt := reflect.TypeOf(v)

	if rt == nil || rt.Kind() != reflect.Struct {
		return ErrValueIsNotStruct
	}

	rv := reflect.ValueOf(v)

	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		tag := field.Tag.Get(tagName)

		if tag == "" {
			continue
		}

		rules := strings.Split(tag, "|")

		vErrs, err := validateField(field.Name, field.Type, rv.Field(i), rules)
		if err != nil {
			return err
		}

		errs = append(errs, vErrs...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}
