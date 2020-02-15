package msgraph

import (
	"fmt"

	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

const groupSelect = "id,displayName,mailNickname,mailEnabled,securityEnabled,mail,groupTypes,visibility"

func (g *graph) groupCreate(group *msgraph.Group) (*msgraph.Group, error) {
	group, err := g.cli.Groups().Request().Add(g.ctx, group)
	if err != nil {
		return nil, err
	}
	return g.groupRead(*group.ID)
}

func (g *graph) groupRead(id string) (*msgraph.Group, error) {
	req := g.cli.Groups().ID(id).Request()
	req.Select(groupSelect)
	return req.Get(g.ctx)
}

func (g *graph) groupReadByMailNickname(nick string) (*msgraph.Group, error) {
	req := g.cli.Groups().Request()
	req.Filter(fmt.Sprintf("mailNickname eq '%s'", nick))
	req.Select(groupSelect)
	groups, err := req.Get(g.ctx)
	if err != nil {
		return nil, err
	}
	if len(groups) != 1 {
		return nil, fmt.Errorf("found %d groups for mail_nickname = %q", len(groups), nick)
	}
	return &groups[0], nil
}

func (g *graph) groupUpdate(id string, group *msgraph.Group) (*msgraph.Group, error) {
	err := g.cli.Groups().ID(id).Request().Update(g.ctx, group)
	if err != nil {
		return nil, err
	}
	return g.groupRead(id)
}

func (g *graph) groupDelete(id string) error {
	return g.cli.Groups().ID(id).Request().Delete(g.ctx)
}
