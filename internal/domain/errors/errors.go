package errors

import "errors"

// создаем глобальные переменные типа error
var (
	ErrTaskIsAlreadyExistsInDB = errors.New("such task is already in db")
	ErrCookie                  = errors.New("cookies doesn't exists in headers")
	ErrBadCredential           = errors.New("bad credentials")
	ErrPermissionDenied        = errors.New("bad credentials")
	ErrDecisionAlreadyMade     = errors.New("decision is already made")
)
