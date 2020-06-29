package msgraph

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func dataUserResource() *schema.Resource {
	return &schema.Resource{
		Read: dataUserRead,
		Schema: map[string]*schema.Schema{
			"id":                  {Type: schema.TypeString, ValidateFunc: validation.IsUUID, Computed: true, Optional: true, ConflictsWith: []string{"user_principal_name", "mail_nickname"}},
			"user_principal_name": {Type: schema.TypeString, Computed: true, Optional: true, ConflictsWith: []string{"id", "mail_nickname"}},
			"display_name":        {Type: schema.TypeString, Computed: true},
			"given_name":          {Type: schema.TypeString, Computed: true},
			"surname":             {Type: schema.TypeString, Computed: true},
			"mail_nickname":       {Type: schema.TypeString, Computed: true, Optional: true, ConflictsWith: []string{"id", "user_principal_name"}},
			"mail":                {Type: schema.TypeString, Computed: true},
			"other_mails":         {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"account_enabled":     {Type: schema.TypeBool, Computed: true},
		},
	}
}

type dataUser struct {
	graph *graph
	data  *schema.ResourceData
}

func newDataUser(d *schema.ResourceData, m interface{}) *dataUser {
	return &dataUser{
		graph: newGraph(m),
		data:  d,
	}
}

func dataUserRead(d *schema.ResourceData, m interface{}) error {
	return newDataUser(d, m).read()
}

func (d *dataUser) graphSet(user *msgraph.User) {
	d.data.Set("user_principal_name", user.UserPrincipalName)
	d.data.Set("display_name", user.DisplayName)
	d.data.Set("given_name", user.GivenName)
	d.data.Set("surname", user.Surname)
	d.data.Set("mail_nickname", user.MailNickname)
	d.data.Set("mail", user.Mail)
	d.data.Set("other_mails", user.OtherMails)
	d.data.Set("account_enabled", user.AccountEnabled)
}

func (d *dataUser) read() error {
	var (
		user *msgraph.User
		err  error
	)
	if id, ok := d.data.GetOkExists("id"); ok {
		user, err = d.graph.userRead(id.(string))
	} else if upn, ok := d.data.GetOkExists("user_principal_name"); ok {
		user, err = d.graph.userRead(upn.(string))
	} else if nick, ok := d.data.GetOkExists("mail_nickname"); ok {
		user, err = d.graph.userReadByMailNickname(nick.(string))
	} else {
		err = fmt.Errorf("one of `id`, `user_principal_name` and `mail_nickname` must be supplied")
	}
	if err != nil {
		d.data.SetId("")
		if errRes, ok := err.(*msgraph.ErrorResponse); ok {
			if errRes.StatusCode() == http.StatusNotFound {
				return nil
			}
		}
		return err
	}
	d.data.SetId(*user.ID)
	d.graphSet(user)
	return nil
}
