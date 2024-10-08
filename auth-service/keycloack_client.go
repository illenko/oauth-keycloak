package main

import (
	"context"
	"log"
	"os"

	"github.com/Nerzal/gocloak/v13"
	"github.com/joho/godotenv"
)

type KeycloakAdminClientService struct {
	client   *gocloak.GoCloak
	realm    string
	clientId string
	secret   string
}

func NewKeycloakAdminClientService() *KeycloakAdminClientService {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	client := gocloak.NewClient(os.Getenv("KEYCLOAK_URL"))
	return &KeycloakAdminClientService{
		client:   client,
		realm:    os.Getenv("KEYCLOAK_REALM"),
		clientId: os.Getenv("KEYCLOAK_CLIENT_ID"),
		secret:   os.Getenv("KEYCLOAK_CLIENT_SECRET"),
	}
}

func (s *KeycloakAdminClientService) LoginUser(request LoginRequest) JwtResponse {
	ctx := context.Background()
	token, err := s.client.Login(ctx, s.clientId, s.secret, s.realm, request.Username, request.Password)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	return JwtResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
}

func (s *KeycloakAdminClientService) RefreshToken(refreshToken string) (JwtResponse, error) {
	ctx := context.Background()
	token, err := s.client.RefreshToken(ctx, refreshToken, s.clientId, s.secret, s.realm)
	if err != nil {
		return JwtResponse{}, err
	}
	return JwtResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (s *KeycloakAdminClientService) ValidateToken(token string) bool {
	ctx := context.Background()
	rptResult, err := s.client.RetrospectToken(ctx, token, s.clientId, s.secret, s.realm)
	if err != nil || !*rptResult.Active {
		return false
	}
	return true
}
