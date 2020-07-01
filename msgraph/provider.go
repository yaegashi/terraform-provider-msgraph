package msgraph

import (
	"context"
	"fmt"
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

// Provider returns msgraph provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Description: "Tenant ID",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_TENANT_ID", defaultTenantID),
			},
			"client_id": {
				Type:        schema.TypeString,
				Description: "Client ID",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_ID", defaultClientID),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Description: "Client Secret",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET", ""),
			},
			"token_cache_path": {
				Type:        schema.TypeString,
				Description: "Token Cache Path",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_TOKEN_CACHE_PATH", defaultTokenCachePath),
			},
			"console_device_path": {
				Type:        schema.TypeString,
				Description: "Console Device Path",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CONSOLE_DEVICE_PATH", ""),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"msgraph_user":  dataUserResource(),
			"msgraph_group": dataGroupResource(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"msgraph_user":                 resourceUserResource(),
			"msgraph_group":                resourceGroupResource(),
			"msgraph_group_member":         resourceGroupMemberResource(),
			"msgraph_application":          resourceApplicationResource(),
			"msgraph_application_password": resourceApplicationPasswordResource(),
		},
		ConfigureFunc: configureFunc,
	}
}

func configureFunc(r *schema.ResourceData) (interface{}, error) {
	tenantID := r.Get("tenant_id").(string)
	clientID := r.Get("client_id").(string)
	clientSecret := r.Get("client_secret").(string)
	tokenCachePath := r.Get("token_cache_path").(string)
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
		ts, err = m.DeviceAuthorizationGrant(ctx, tenantID, clientID, []string{"offline_access", msauth.DefaultMSGraphScope},
			func(dc *msauth.DeviceCode) error {
				log.Println(dc.Message)
				con, err := openConsole(r.Get("console_device_path").(string))
				if err == nil {
					fmt.Fprintln(con, dc.Message)
					con.Close()
				} else {
					log.Printf("Failed to open console: %s", err)
				}
				return nil
			})
		if err != nil {
			return nil, err
		}
		m.SaveFile(tokenCachePath)
	}
	httpClient := oauth2.NewClient(ctx, ts)
	graphClient := msgraph.NewClient(httpClient)
	return graphClient, nil
}
