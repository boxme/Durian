package main

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, countryCode int, mobileNumber string, password string) (string, error)
	Register(ctx context.Context, name string, countryCode int, mobileNumber string, password string) (string, error)
}

type authService struct {
	Db           *pgxpool.Pool
	jwtSecretKey []byte
}

type password struct {
	plaintextPassword *string
	hash              []byte
}

type JwtClaims struct {
	MobileNumber string
	jwt.RegisteredClaims
}

func CreateAuthService(db *pgxpool.Pool) AuthService {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	return &authService{
		Db:           db,
		jwtSecretKey: []byte(secretKey),
	}
}

func (as *authService) Register(
	ctx context.Context,
	name string,
	countryCode int,
	mobileNumber string,
	password string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return "", nil
	}
}

func (as *authService) Login(
	ctx context.Context,
	countryCode int,
	mobileNumber string,
	password string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		query := `SELECT id, hashed_password FROM users WHERE phone_number = @phoneNumber AND country_code = @countryCode;`
		args := pgx.NamedArgs{
			"phoneNumber": mobileNumber,
			"countryCode": countryCode}

		var userId uint64
		var hashedPassword string
		if err := as.Db.QueryRow(ctx, query, args).Scan(&userId, &hashedPassword); err != nil {
			return "", err
		}

		passwordMatched, err := isMatched(hashedPassword, password)
		if err != nil {
			return "", err
		}

		if !passwordMatched {
			return "", errors.New("wrong credentials")
		}

		expirationTime := time.Now().Add(5 * time.Minute)
		claim := jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		}
		jwtClaim := &JwtClaims{MobileNumber: mobileNumber, RegisteredClaims: claim}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim)
		tokenString, err := token.SignedString(as.jwtSecretKey)
		if err != nil {
			return "", err
		}

		return tokenString, nil
	}
}

func encrypt(plaintext string) (password, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	var encrypted password
	if err != nil {
		return encrypted, err
	}

	encrypted = password{
		plaintextPassword: &plaintext,
		hash:              hash,
	}

	return encrypted, nil
}

func isMatched(hashed string, plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plaintext))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil

		default:
			return false, err
		}
	}

	return true, nil
}
