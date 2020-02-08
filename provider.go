package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/yaegashi/msgraph.go/msauth"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
	"golang.org/x/oauth2"
)

const (
	defaultTenantID       = "common"
	defaultClientID       = "82492584-8587-4e7d-ad48-19546ce8238f"
	defaultTokenCachePath = "token_cache.json"
)

func provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Description: "Tenant ID",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSGRAPH_TENANT_ID", defaultTenantID),
			},
			"client_id": {
				Type:        schema.TypeString,
				Description: "Client ID",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSGRAPH_CLIENT_ID", defaultClientID),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Description: "Client Secret",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSGRAPH_CLIENT_SECRET", ""),
			},
			"token_cache_path": {
				Type:        schema.TypeString,
				Description: "Token Cache Path",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSGRAPH_TOKEN_CACHE_PATH", defaultTokenCachePath),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"msgraph_user": dataUserResource(),
		},
		ConfigureFunc: configureFunc,
	}
}

func configureFunc(data *schema.ResourceData) (interface{}, error) {
	tenantID := data.Get("tenant_id").(string)
	clientID := data.Get("client_id").(string)
	clientSecret := data.Get("client_secret").(string)
	tokenCachePath := data.Get("token_cache_path").(string)
	ctx := context.Background()
	m := msauth.NewManager()
	var ts oauth2.TokenSource
	var err error
	if len(clientSecret) > 0 {
		ts, err = m.ClientCredentialsGrant(ctx, tenantID, clientID, clientSecret, []string{msauth.DefaultMSGraphScope})
		if err != nil {
			return nil, err
		}
	} else {
		m.LoadFile(tokenCachePath)
		ts, err = m.DeviceAuthorizationGrant(ctx, tenantID, clientID, []string{"offline_access", msauth.DefaultMSGraphScope}, func(dc *msauth.DeviceCode) error { log.Println(dc.Message); return nil })
		if err != nil {
			return nil, err
		}
		m.SaveFile(tokenCachePath)
	}
	httpClient := oauth2.NewClient(ctx, ts)
	graphClient := msgraph.NewClient(httpClient)
	return graphClient, nil
}
