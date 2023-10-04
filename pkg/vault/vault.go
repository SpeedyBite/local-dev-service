package vault

import (
	"context"
	"fmt"
	"log"
	"os"

	vault "github.com/hashicorp/vault/api"
)

type Vault struct {
	Client *vault.Client
	Config *vault.Config
}

func NewVault(server string) *Vault {
	config := vault.DefaultConfig()
	config.Address = server

	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}

	token := os.Getenv("VAULT_TOKEN")
	client.SetToken(token)

	return &Vault{
		Config: config,
		Client: client,
	}
}

func (vault *Vault) GetVaults(environment string, serviceName string) (*vault.KVSecret, error) {
	secret, err := vault.Client.KVv2(environment).Get(context.Background(), serviceName)
	if err != nil {
		return nil, fmt.Errorf("unable to read secret: %w", err)
	}
	return secret, nil
}
