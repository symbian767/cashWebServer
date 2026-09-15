package services

import (
	"context"
	"document-service/models"
	"errors"
	"math/rand"
	"regexp"
	"unicode"
	"unicode/utf8"
)

var Session map[string]struct{}

const (
	keyLetter = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func CreateUser(ctx context.Context, user models.User) error {
	// Проверка логина и пароля
	if !validateLogin(user.Login) {
		return errors.New("invalid login")
	}
	if !validatePassword(user.Password) {
		return errors.New("invalid password")
	}
	// Создание нового пользователя
	return DB.CreateUser(ctx, user)
}

func genToken() string {
	bytesSlice := make([]byte, 16)
	for i := range bytesSlice {
		bytesSlice[i] = keyLetter[rand.Int63()%int64(len(keyLetter))]
	}
	return string(bytesSlice)
}

func AuthenticateUser(ctx context.Context, login, password string) (string, error) {
	// Аутентификация пользователя
	res, err := DB.AuthenticateUser(ctx, login, password)
	if err != nil {
		return "", err
	}
	if !res {
		return "", errors.New("неверно указан логин или пароль")
	}
	return genToken(), nil
}

func validateLogin(login string) bool {
	if utf8.RuneCountInString(login) < 8 {
		return false
	}
	regex := regexp.MustCompile(`^[a-zA-Z0-9]{8,}$`)
	return regex.MatchString(login)
}

func validatePassword(password string) bool {
	if utf8.RuneCountInString(password) < 8 {
		return false
	}

	var upper, lower, digit, special int
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			upper++
		case unicode.IsLower(r):
			lower++
		case unicode.IsDigit(r):
			digit++
		default:
			special++
		}
	}

	return upper >= 1 && lower >= 1 && digit >= 1 && special >= 1
}
