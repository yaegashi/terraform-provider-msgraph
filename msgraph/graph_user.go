package msgraph

import (
	"fmt"

	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

const userSelect = "id,userPrincipalName,mailNickname,displayName,companyName,surname,givenName,otherMails,accountEnabled"

func (g *graph) userCreate(user *msgraph.User) (*msgraph.User, error) {
	user, err := g.cli.Users().Request().Add(g.ctx, user)
	if err != nil {
		return nil, err
	}
	return g.userRead(*user.ID)
}

func (g *graph) userRead(id string) (*msgraph.User, error) {
	req := g.cli.Users().ID(id).Request()
	req.Select(userSelect)
	return req.Get(g.ctx)
}

func (g *graph) userReadByMailNickname(nick string) (*msgraph.User, error) {
	req := g.cli.Users().Request()
	req.Filter(fmt.Sprintf("mailNickname eq '%s'", nick))
	req.Select(userSelect)
	users, err := req.Get(g.ctx)
	if err != nil {
		return nil, err
	}
	if len(users) != 1 {
		return nil, fmt.Errorf("found %d users for mail_nickname = %q", len(users), nick)
	}
	return &users[0], nil
}

func (g *graph) userUpdate(id string, user *msgraph.User) (*msgraph.User, error) {
	err := g.cli.Users().ID(id).Request().Update(g.ctx, user)
	if err != nil {
		return nil, err
	}
	return g.userRead(id)
}

func (g *graph) userDelete(id string) error {
	return g.cli.Users().ID(id).Request().Delete(g.ctx)
}
