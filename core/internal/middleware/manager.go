package middleware

import (
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

type Manager struct {
	clerkSecret string
	userClient  *user.Client
}

func NewManager(clerkSecret string, cfg *clerk.ClientConfig) *Manager {
	userClient := user.NewClient(cfg)
	return &Manager{
		clerkSecret: clerkSecret,
		userClient:  userClient,
	}
}
