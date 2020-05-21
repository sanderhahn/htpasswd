package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"hash"
	"time"
)

// Secret contains the secret configuration
type Secret struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

var errInvalidJWTSecret = errors.New("Invalid JWT secret")

func (s *Secret) getHash() (h func() hash.Hash, err error) {
	switch s.Type {
	case "HS256":
		h = sha256.New
	case "HS384":
		h = sha512.New384
	case "HS512":
		h = sha512.New
	default:
		err = errInvalidJWTSecret
	}
	return h, err
}

// Header is the jwt header
type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// Payload is the jwt payload
type Payload struct {
	Nbf    int64   `json:"nbf"`
	Exp    int64   `json:"exp"`
	Claims *Claims `json:"https://hasura.io/jwt/claims"`
}

// Claims are the Hasura claims
type Claims struct {
	Username     string   `json:"x-hasura-user-id"`
	DefaultRole  string   `json:"x-hasura-default-role"`
	AllowedRoles []string `json:"x-hasura-allowed-roles"`
}

// SignPayload signs the payload
func (s *Secret) SignPayload(payload *Payload) (signed string, err error) {
	header := &Header{
		Alg: s.Type,
		Typ: "JWT",
	}
	jsonHeader, err := json.Marshal(header)
	if err != nil {
		return
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return
	}
	hash, err := s.getHash()
	if err != nil {
		return
	}
	h := hmac.New(hash, []byte(s.Key))
	part := base64.RawURLEncoding.EncodeToString([]byte(jsonHeader)) + "." +
		base64.RawURLEncoding.EncodeToString([]byte(jsonPayload))
	h.Write([]byte(part))
	signed = part + "." + base64.RawURLEncoding.EncodeToString([]byte(h.Sum(nil)))
	return
}

// Sign username and role
func (s *Secret) Sign(expiresIn time.Duration, username, role string) (string, error) {
	now := time.Now().UTC()
	payload := &Payload{
		Nbf: now.Unix(),
		Exp: now.Add(expiresIn).Unix(),
		Claims: &Claims{
			Username:     username,
			DefaultRole:  role,
			AllowedRoles: []string{role},
		},
	}
	buf, err := s.SignPayload(payload)
	return buf, err
}
