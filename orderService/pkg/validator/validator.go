package validator

import (
	"fmt"
	"regexp"
	"user-svc/internal/service/param"
	"user-svc/pkg/errmsg"
	"user-svc/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
}

type Config struct{}

type Validator struct {
	repo Repository
}

func New(repo Repository) Validator {
	return Validator{repo: repo}
}

func (v Validator) ValidateRegisterRequest(req param.RegisterRequest) (map[string]string, error) {

	const op = "uservalidator.ValidateRegisterRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required, validation.Length(3, 50)),

		validation.Field(&req.Password, validation.Required,
			validation.Match(regexp.MustCompile("^[A-Za-z0-9]{4,}$"))),

		validation.Field(&req.PhoneNumber, validation.Required,
			validation.Match(regexp.MustCompile("^09[0-9]{9,}$")).Error(""),
			validation.By(v.checkPhoneUniqueness)),
	); err != nil {

		// fmt.Println("error validator is :", err)
		// fmt.Printf("type error is : %T\n", err)

		fieldErrors := make(map[string]string)
		errV, ok := err.(validation.Errors)
		if ok {
			for field, errs := range errV {
				if errs != nil {
					fieldErrors[field] = errs.Error()
				}
			}
		}

		return fieldErrors, richerror.New(op).WithMessage(errmsg.ErrorMsgInvalidInput).
			WithKind(richerror.KindInvalid).WithErr(err).
			WithMeta(map[string]interface{}{"req": req})
	}

	return nil, nil
}

// * Custom  Validation
func (v Validator) checkPhoneUniqueness(value interface{}) error {

	phoneNumber := value.(string)

	if isUnique, err := v.repo.IsPhoneNumberUnique(phoneNumber); err != nil || !isUnique {

		if err != nil {
			return err
		}

		if !isUnique {
			return fmt.Errorf(errmsg.ErrorMsgPhoneNumberIsNotUnique)
		}
	}
	return nil
}
