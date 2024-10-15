package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"tech_check/internal/constants"
	"tech_check/internal/dto"
	"tech_check/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
)

type Auth struct {
	rTokenExpiresHour int
	aTokenExpiresHour int
	rTokenLength      int
	googleClientID    string
	jwtSecret         []byte
	userSrvc          UserSrvc
	refreshTokenSrvc  RefreshTokenSrvc
}

func NewAuth(
	googleClientID string,
	jwtSecret string,
	userSrvc UserSrvc,
	refreshTokenSrvc RefreshTokenSrvc,
) *Auth {
	return &Auth{
		rTokenExpiresHour: 24,
		aTokenExpiresHour: 2,
		rTokenLength:      50,
		googleClientID:    googleClientID,
		jwtSecret:         []byte(jwtSecret),
		userSrvc:          userSrvc,
		refreshTokenSrvc:  refreshTokenSrvc,
	}
}

func (a *Auth) Login(ctx context.Context, email, password, ip string) (*dto.Token, error) {
	const op = "services.Auth.Login"

	user, err := a.validateCredential(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	token, err := a.createToken(ctx, user, ip)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) GoogleLogin(ctx context.Context, tokenID, ip string) (*dto.Token, error) {
	const op = "services.Auth.GoogleLogin"

	payload, err := idtoken.Validate(ctx, tokenID, a.googleClientID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, constants.ErrInvalidGoogleData)
	}
	name, ok := payload.Claims["name"].(string)
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, constants.ErrInvalidGoogleData)
	}
	avatar, ok := payload.Claims["picture"].(string)
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, constants.ErrInvalidGoogleData)
	}

	user, err := a.userSrvc.GetOrCreate(ctx, email, name, avatar)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	token, err := a.createToken(ctx, user, ip)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) DecodeAToken(ctx context.Context, aToken string) (*dto.Claims, error) {
	const op = "services.Auth.DecodeAToken"

	token, err := jwt.ParseWithClaims(aToken, &dto.Claims{}, a.getSigningKey)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	claims, ok := token.Claims.(*dto.Claims)
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, constants.ErrInvalidClaimsType)
	}

	return claims, nil
}

func (a *Auth) Refresh(ctx context.Context, aToken, rToken, ip string) (*dto.Token, error) {
	const op = "services.Auth.Refresh"

	claims, err := a.DecodeAToken(ctx, aToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user, err := a.userSrvc.GetByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, constants.ErrCannotLogin)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if claims.IP != ip {
		// TODO: send email
		log.Printf("sending email to %s", user.Email)
	}

	err = a.validateRToken(ctx, user, claims.RefreshTokenID, rToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	token, err := a.createToken(ctx, user, ip)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) validateCredential(ctx context.Context, email, password string) (*models.User, error) {
	const op = "services.Auth.validateCredential"

	user, err := a.userSrvc.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, constants.ErrInvalidCredentials)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, fmt.Errorf("%s: %w", op, constants.ErrInvalidCredentials)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (a *Auth) createRToken(ctx context.Context, user *models.User, ip string) (*models.RefreshToken, string, error) {
	const op = "services.Auth.createRToken"

	err := a.refreshTokenSrvc.Delete(ctx, user)
	if err != nil && !errors.Is(err, constants.ErrNotFound) {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	random := make([]byte, a.rTokenLength)
	_, err = rand.Read(random)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	rToken := base64.StdEncoding.EncodeToString(random)
	hash, err := bcrypt.GenerateFromPassword([]byte(rToken), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	refreshToken, err := a.refreshTokenSrvc.Create(
		ctx,
		user,
		ip,
		string(hash),
		time.Now().Add(time.Duration(a.rTokenExpiresHour)*time.Hour),
	)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	return refreshToken, rToken, nil
}

func (a *Auth) generateAToken(user *models.User, refreshToken *models.RefreshToken, ip string) (string, error) {
	const op = "services.Auth.generateAToken"

	claims := dto.NewClaims(user, refreshToken, ip, a.aTokenExpiresHour)
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(a.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) getSigningKey(token *jwt.Token) (interface{}, error) {
	const op = "srv.Auth.getSigningKey"

	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, constants.ErrInvalidSigningMethod)
	}

	return a.jwtSecret, nil
}

func (a *Auth) validateRToken(ctx context.Context, user *models.User, refreshTokenID, rToken string) error {
	const op = "services.Auth.validateRToken"

	refreshToken, err := a.refreshTokenSrvc.GetByID(ctx, user, refreshTokenID)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return fmt.Errorf("%s: %w", op, constants.ErrTokensMismatch)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return fmt.Errorf("%s: %w", op, constants.ErrRTokenExpired)
	}

	err = bcrypt.CompareHashAndPassword([]byte(refreshToken.Hash), []byte(rToken))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return fmt.Errorf("%s: %w", op, constants.ErrInvalidRToken)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *Auth) createToken(ctx context.Context, user *models.User, ip string) (*dto.Token, error) {
	const op = "services.Auth.createToken"

	refreshToken, rToken, err := a.createRToken(ctx, user, ip)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	aToken, err := a.generateAToken(user, refreshToken, ip)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	token := dto.NewToken(aToken, rToken)

	return token, nil
}
