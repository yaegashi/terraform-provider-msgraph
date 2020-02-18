package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/yaegashi/msgraph.go/jsonx"
	"github.com/yaegashi/msgraph.go/msauth"
	P "github.com/yaegashi/msgraph.go/ptr"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
	"golang.org/x/oauth2"
)

func dump(o interface{}) {
	enc := jsonx.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(o)
}

func main() {
	var tenantDomain, tenantID, clientID, clientSecret string
	var userCount, groupCount int
	var clean bool
	flag.StringVar(&tenantDomain, "tenant-domain", "l0wdev.onmicrosoft.com", "Tenant Domain")
	flag.StringVar(&tenantID, "tenant-id", os.Getenv("ARM_TENANT_ID"), "Tenant ID")
	flag.StringVar(&clientID, "client-id", os.Getenv("ARM_CLIENT_ID"), "Client ID")
	flag.StringVar(&clientSecret, "client-secret", os.Getenv("ARM_CLIENT_SECRET"), "Client Secret")
	flag.IntVar(&userCount, "user-count", 100, "User count")
	flag.IntVar(&groupCount, "group-count", 10, "Group count")
	flag.BoolVar(&clean, "clean", false, "Clean")
	flag.Parse()

	ctx := context.Background()
	m := msauth.NewManager()
	scopes := []string{msauth.DefaultMSGraphScope}
	ts, err := m.ClientCredentialsGrant(ctx, tenantID, clientID, clientSecret, scopes)
	if err != nil {
		log.Fatal(err)
	}

	httpClient := oauth2.NewClient(ctx, ts)
	graphClient := msgraph.NewClient(httpClient)

	if clean {
		for i := 0; i < userCount; i++ {
			r := graphClient.Users().Request()
			r.Filter(fmt.Sprintf("mailNickname eq 'msgraphgouser%d'", i))
			users, _ := r.Get(ctx)
			for _, user := range users {
				log.Printf("Delete user %d", i)
				graphClient.Users().ID(*user.ID).Request().Delete(ctx)
			}
		}
		for i := 0; i < groupCount; i++ {
			r := graphClient.Groups().Request()
			r.Filter(fmt.Sprintf("mailNickname eq 'msgraphgogroup%d'", i))
			groups, _ := r.Get(ctx)
			for _, group := range groups {
				log.Printf("Delete group %d", i)
				graphClient.Groups().ID(*group.ID).Request().Delete(ctx)
			}
		}
		return
	}

	users := make([]*msgraph.User, userCount)
	for i := 0; i < userCount; i++ {
		user := &msgraph.User{
			AccountEnabled:    P.Bool(true),
			UserPrincipalName: P.String(fmt.Sprintf("msgraphgouser%d@%s", i, tenantDomain)),
			MailNickname:      P.String(fmt.Sprintf("msgraphgouser%d", i)),
			DisplayName:       P.String(fmt.Sprintf("msgraphgo user %d", i)),
			PasswordProfile: &msgraph.PasswordProfile{
				ForceChangePasswordNextSignIn: P.Bool(false),
				Password:                      P.String(uuid.New().String()),
			},
		}
		log.Printf("Create user %d", i)
		users[i], err = graphClient.Users().Request().Add(ctx, user)
		if err != nil {
			log.Fatal(err)
		}
	}

	groups := make([]*msgraph.Group, groupCount)
	for i := 0; i < groupCount; i++ {
		group := &msgraph.Group{
			SecurityEnabled: P.Bool(true),
			MailEnabled:     P.Bool(false),
			MailNickname:    P.String(fmt.Sprintf("msgraphgogroup%d", i)),
			DisplayName:     P.String(fmt.Sprintf("msgraphgo group %d", i)),
		}
		log.Printf("Create group %d", i)
		groups[i], err = graphClient.Groups().Request().Add(ctx, group)
		if err != nil {
			log.Fatal(err)
		}
	}

	for i := 0; i < userCount; i++ {
		member := users[i]
		group := groups[i%groupCount]
		log.Printf("Add user %d to group %d", i, i%groupCount)
		reqObj := map[string]interface{}{
			"@odata.id": graphClient.DirectoryObjects().ID(*member.ID).Request().URL(),
		}
		r := graphClient.Groups().ID(*group.ID).Members().Request()
		err := r.JSONRequest(ctx, "POST", "/$ref", reqObj, nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	for i := 0; i < groupCount-1; i++ {
		member := groups[i+1]
		group := groups[i]
		log.Printf("Add group %d to group %d", i+1, i)
		reqObj := map[string]interface{}{
			"@odata.id": graphClient.DirectoryObjects().ID(*member.ID).Request().URL(),
		}
		r := graphClient.Groups().ID(*group.ID).Members().Request()
		err := r.JSONRequest(ctx, "POST", "/$ref", reqObj, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}
