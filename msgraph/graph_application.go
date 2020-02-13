package msgraph

import msgraph "github.com/yaegashi/msgraph.go/v1.0"

func (g *graph) applicationCreate(application *msgraph.Application) (*msgraph.Application, error) {
	application, err := g.cli.Applications().Request().Add(g.ctx, application)
	if err != nil {
		return nil, err
	}
	return g.applicationRead(*application.ID)
}

func (g *graph) applicationRead(id string) (*msgraph.Application, error) {
	req := g.cli.Applications().ID(id).Request()
	//req.Select("id,displayName,mailNickname,mailEnabled,securityEnabled,mail,groupTypes,visibility")
	return req.Get(g.ctx)
}

func (g *graph) applicationUpdate(id string, application *msgraph.Application) (*msgraph.Application, error) {
	err := g.cli.Applications().ID(id).Request().Update(g.ctx, application)
	if err != nil {
		return nil, err
	}
	return g.applicationRead(id)
}

func (g *graph) applicationDelete(id string) error {
	return g.cli.Applications().ID(id).Request().Delete(g.ctx)
}

func (g *graph) applicationAddPassword(id string, reqObj *msgraph.ApplicationAddPasswordRequestParameter) (*msgraph.PasswordCredential, error) {
	return g.cli.Applications().ID(id).AddPassword(reqObj).Request().Post(g.ctx)
}

func (g *graph) applicationRemovePassword(id string, reqObj *msgraph.ApplicationRemovePasswordRequestParameter) error {
	return g.cli.Applications().ID(id).RemovePassword(reqObj).Request().Post(g.ctx)
}
