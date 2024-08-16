package vault

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"dnevnik-rg.ru/config"
	"github.com/hashicorp/vault-client-go"
)

type VaultClient struct {
	Client *vault.Client
	Token  string
}

const clientTimeout = 10 * time.Second

func NewVaultClient(vaultCfg *config.VaultCfg) (*VaultClient, error) {
	client, err := vault.New(
		vault.WithAddress(
			fmt.Sprintf(
				"http://%s:%s",
				vaultCfg.Host,
				vaultCfg.Port,
			),
		),
		vault.WithRequestTimeout(clientTimeout),
	)

	if err != nil {
		return nil, err
	}

	vlt := &VaultClient{
		Client: client,
		// token будем ставить через пайплайн во время билда
		Token: os.Getenv("VAULT_DEV_ROOT_TOKEN_ID"),
	}

	if err := client.SetToken(vlt.Token); err != nil {
		return nil, err
	}

	return vlt, nil
}

func (v *VaultClient) GetVaultData(ctx context.Context, route, key string) (string, error) {
	response, err := v.Client.Read(ctx, route)
	if err != nil {
		return "", err
	}

	data, ok := response.Data["data"].(map[string]interface{})
	if !ok {
		return "", errors.New("cannot assert type from vault response")
	}

	return data[key].(string), nil
}
