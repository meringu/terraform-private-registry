package utils

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GitHubGlientForInstallation returns a github client for the installation
func GitHubGlientForInstallation(ctx context.Context, baseURL *url.URL, appID, id int64, privateKey []byte) (*github.Client, error) {
	claims := &jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute).Unix(),
		Issuer:    strconv.FormatInt(appID, 10),
	}

	bearer := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, err
	}
	ss, err := bearer.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("could not sign jwt: %s", err)
	}

	ghclient := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ss},
	)))
	ghclient.BaseURL = baseURL

	token, _, err := ghclient.Apps.CreateInstallationToken(ctx, id)
	if err != nil {
		return nil, err
	}

	ghclient = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.GetToken()},
	)))
	ghclient.BaseURL = baseURL

	return ghclient, nil
}
