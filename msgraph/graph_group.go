package msgraph

import msgraph "github.com/yaegashi/msgraph.go/v1.0"

func (g *graph) groupCreate(group *msgraph.Group) (*msgraph.Group, error) {
	group, err := g.cli.Groups().Request().Add(g.ctx, group)
	if err != nil {
		return nil, err
	}
	return g.groupRead(*group.ID)
}

func (g *graph) groupRead(id string) (*msgraph.Group, error) {
	req := g.cli.Groups().ID(id).Request()
	req.Select("id,displayName,mailNickname,mailEnabled,securityEnabled,mail,groupTypes,visibility")
	return req.Get(g.ctx)
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
