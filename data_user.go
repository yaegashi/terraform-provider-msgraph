package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	P "github.com/yaegashi/msgraph.go/ptr"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func dataUserResource() *schema.Resource {
	return &schema.Resource{
		Create: dataUserCreate,
		Read:   dataUserRead,
		Update: dataUserUpdate,
		Delete: dataUserDelete,
		Schema: map[string]*schema.Schema{
			"user_principal_name": &schema.Schema{Type: schema.TypeString, Required: true},
			"display_name":        &schema.Schema{Type: schema.TypeString, Optional: true},
			"given_name":          &schema.Schema{Type: schema.TypeString, Optional: true},
			"surname":             &schema.Schema{Type: schema.TypeString, Optional: true},
			"mail_nickname":       &schema.Schema{Type: schema.TypeString, Optional: true},
			"mail":                &schema.Schema{Type: schema.TypeString, Computed: true},
			"other_mails":         &schema.Schema{Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
		},
	}
}

type dataUser struct {
	ctx   context.Context
	graph *msgraph.GraphServiceRequestBuilder
	data  *schema.ResourceData
}

func newDataUser(d *schema.ResourceData, m interface{}) *dataUser {
	return &dataUser{
		ctx:   context.Background(),
		graph: m.(*msgraph.GraphServiceRequestBuilder),
		data:  d,
	}
}

func dataUserCreate(d *schema.ResourceData, m interface{}) error {
	return newDataUser(d, m).dataCreate()
}

func dataUserRead(d *schema.ResourceData, m interface{}) error {
	return newDataUser(d, m).dataRead()
}

func dataUserUpdate(d *schema.ResourceData, m interface{}) error {
	return newDataUser(d, m).dataUpdate()
}

func dataUserDelete(d *schema.ResourceData, m interface{}) error {
	return newDataUser(d, m).dataDelete()
}

func (d *dataUser) dataGraphSet(user *msgraph.User) {
	d.data.Set("user_principal_name", user.UserPrincipalName)
	d.data.Set("display_name", user.DisplayName)
	d.data.Set("given_name", user.GivenName)
	d.data.Set("surname", user.Surname)
	d.data.Set("mail_nickname", user.MailNickname)
	d.data.Set("mail", user.Mail)
	d.data.Set("other_mails", user.OtherMails)
}

func (d *dataUser) dataGraphGet() *msgraph.User {
	user := &msgraph.User{
		UserPrincipalName: P.String(d.data.Get("user_principal_name").(string)),
		DisplayName:       P.String(d.data.Get("display_name").(string)),
		GivenName:         P.String(d.data.Get("given_name").(string)),
		Surname:           P.String(d.data.Get("surname").(string)),
		MailNickname:      P.String(d.data.Get("mail_nickname").(string)),
	}
	for _, i := range d.data.Get("other_mails").([]interface{}) {
		user.OtherMails = append(user.OtherMails, i.(string))
	}
	return user
}

func (d *dataUser) dataCreate() error {
	newUser := d.dataGraphGet()
	newUser.AccountEnabled = P.Bool(true)
	newUser.PasswordProfile = &msgraph.PasswordProfile{
		ForceChangePasswordNextSignIn: P.Bool(false),
		Password:                      P.String(uuid.New().String()), // XXX: random password
	}
	user, err := d.graphCreate(newUser)
	if err != nil {
		return err
	}
	d.data.SetId(*user.ID)
	d.dataGraphSet(user)
	return nil
}

func (d *dataUser) dataRead() error {
	user, err := d.graphRead(d.data.Id())
	if err != nil {
		d.data.SetId("")
		if errRes, ok := err.(*msgraph.ErrorResponse); ok {
			if errRes.StatusCode() == http.StatusNotFound {
				return nil
			}
		}
		return err
	}
	d.dataGraphSet(user)
	return nil
}

func (d *dataUser) dataUpdate() error {
	newUser := d.dataGraphGet()
	user, err := d.graphUpdate(d.data.Id(), newUser)
	if err != nil {
		return err
	}
	d.dataGraphSet(user)
	return nil
}

func (d *dataUser) dataDelete() error {
	d.graphDelete(d.data.Id())
	d.data.SetId("")
	return nil
}

func (d *dataUser) graphCreate(user *msgraph.User) (*msgraph.User, error) {
	user, err := d.graph.Users().Request().Add(d.ctx, user)
	if err != nil {
		return nil, err
	}
	return d.graphRead(*user.ID)
}

func (d *dataUser) graphRead(id string) (*msgraph.User, error) {
	req := d.graph.Users().ID(id).Request()
	req.Select("id,userPrincipalName,mailNickname,displayName,companyName,surname,givenName,otherMails")
	return req.Get(d.ctx)
}

func (d *dataUser) graphUpdate(id string, user *msgraph.User) (*msgraph.User, error) {
	err := d.graph.Users().ID(id).Request().Update(d.ctx, user)
	if err != nil {
		return nil, err
	}
	return d.graphRead(id)
}

func (d *dataUser) graphDelete(id string) error {
	return d.graph.Users().ID(id).Request().Delete(d.ctx)
}
