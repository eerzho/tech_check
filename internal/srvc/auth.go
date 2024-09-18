package srvc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"tech_check/internal/def"
	"tech_check/internal/dto"
	"tech_check/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	rTokenExpiresHour int
	aTokenExpiresHour int
	rTokenLength      int
	jwtSecret         []byte
	userSrvc          UserSrvc
	refreshTokenSrvc  RefreshTokenSrvc
}

func NewAuth(jwtSecret string, userSrvc UserSrvc, refreshTokenSrvc RefreshTokenSrvc) *Auth {
	return &Auth{
		rTokenExpiresHour: 24,
		aTokenExpiresHour: 2,
		rTokenLength:      50,
		jwtSecret:         []byte(jwtSecret),
		userSrvc:          userSrvc,
		refreshTokenSrvc:  refreshTokenSrvc,
	}
}

func (a *Auth) Login(ctx context.Context, email, password, ip string) (*dto.Token, error) {
	const op = "srvc.Auth.Login"

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

func (a *Auth) DecodeAToken(ctx context.Context, aToken string) (*dto.Claims, error) {
	const op = "srvc.Auth.DecodeAToken"

	token, err := jwt.ParseWithClaims(aToken, &dto.Claims{}, a.getSigningKey)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	claims, ok := token.Claims.(*dto.Claims)
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, def.ErrInvalidClaimsType)
	}

	return claims, nil
}

func (a *Auth) Refresh(ctx context.Context, aToken, rToken, ip string) (*dto.Token, error) {
	const op = "srvc.Auth.Refresh"

	claims, err := a.DecodeAToken(ctx, aToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user, err := a.userSrvc.GetByID(ctx, claims.UserID.Hex())
	if err != nil {
		if errors.Is(err, def.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, def.ErrCannotLogin)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if claims.IP != ip {
		// TODO: send email
		log.Printf("sending email to %s", user.Email)
	}

	err = a.validateRToken(ctx, user, claims.RefreshTokenID.Hex(), rToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	token, err := a.createToken(ctx, user, ip)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) validateCredential(ctx context.Context, email, password string) (*model.User, error) {
	const op = "srvc.Auth.validateCredential"

	user, err := a.userSrvc.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, def.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, def.ErrInvalidCredentials)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, fmt.Errorf("%s: %w", op, def.ErrInvalidCredentials)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (a *Auth) createRTokenByUser(ctx context.Context, user *model.User, ip string) (*model.RefreshToken, string, error) {
	const op = "srvc.Auth.createRTokenByUser"

	err := a.refreshTokenSrvc.DeleteByUser(ctx, user)
	if err != nil && !errors.Is(err, def.ErrNotFound) {
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

	refreshToken, err := a.refreshTokenSrvc.CreateByUser(
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

func (a *Auth) generateAToken(user *model.User, refreshToken *model.RefreshToken, ip string) (string, error) {
	const op = "srvc.Auth.generateAToken"

	claims := dto.Claims{
		IP:             ip,
		UserID:         user.ID,
		RefreshTokenID: refreshToken.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(a.aTokenExpiresHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

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
		return nil, fmt.Errorf("%s: %w", op, def.ErrInvalidSigningMethod)
	}

	return a.jwtSecret, nil
}

func (a *Auth) validateRToken(ctx context.Context, user *model.User, refreshTokenID, rToken string) error {
	const op = "srvc.Auth.validateRToken"

	refreshToken, err := a.refreshTokenSrvc.GetByUserAndID(ctx, user, refreshTokenID)
	if err != nil {
		if errors.Is(err, def.ErrNotFound) {
			return fmt.Errorf("%s: %w", op, def.ErrTokensMismatch)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return fmt.Errorf("%s: %w", op, def.ErrRTokenExpired)
	}

	err = bcrypt.CompareHashAndPassword([]byte(refreshToken.Hash), []byte(rToken))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return fmt.Errorf("%s: %w", op, def.ErrInvalidRToken)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *Auth) createToken(ctx context.Context, user *model.User, ip string) (*dto.Token, error) {
	const op = "srvc.Auth.createToken"

	refreshToken, rToken, err := a.createRTokenByUser(ctx, user, ip)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	aToken, err := a.generateAToken(user, refreshToken, ip)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.Token{
		AToken: aToken,
		RToken: rToken,
	}, nil
}
