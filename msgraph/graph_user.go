package msgraph

import msgraph "github.com/yaegashi/msgraph.go/v1.0"

func (g *graph) userCreate(user *msgraph.User) (*msgraph.User, error) {
	user, err := g.cli.Users().Request().Add(g.ctx, user)
	if err != nil {
		return nil, err
	}
	return g.userRead(*user.ID)
}

func (g *graph) userRead(id string) (*msgraph.User, error) {
	req := g.cli.Users().ID(id).Request()
	req.Select("id,userPrincipalName,mailNickname,displayName,companyName,surname,givenName,otherMails,accountEnabled")
	return req.Get(g.ctx)
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
