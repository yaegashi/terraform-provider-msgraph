package msgraph

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	P "github.com/yaegashi/msgraph.go/ptr"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func resourceUserResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_principal_name":                {Type: schema.TypeString, Required: true},
			"display_name":                       {Type: schema.TypeString, Required: true},
			"given_name":                         {Type: schema.TypeString, Optional: true},
			"surname":                            {Type: schema.TypeString, Optional: true},
			"mail_nickname":                      {Type: schema.TypeString, Required: true},
			"mail":                               {Type: schema.TypeString, Computed: true},
			"other_mails":                        {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"account_enabled":                    {Type: schema.TypeBool, Required: true},
			"password":                           {Type: schema.TypeString, Computed: true, Optional: true, Sensitive: true},
			"force_change_password_next_sign_in": {Type: schema.TypeBool, Computed: true, Optional: true},
		},
	}
}

type resourceUser struct {
	graph    *graph
	resource *schema.ResourceData
}

func newResourceUser(d *schema.ResourceData, meta interface{}) *resourceUser {
	return &resourceUser{
		graph:    newGraph(meta),
		resource: d,
	}
}

func resourceUserCreate(d *schema.ResourceData, meta interface{}) error {
	return newResourceUser(d, meta).create()
}

func resourceUserRead(d *schema.ResourceData, meta interface{}) error {
	return newResourceUser(d, meta).read()
}

func resourceUserUpdate(d *schema.ResourceData, meta interface{}) error {
	return newResourceUser(d, meta).update()
}

func resourceUserDelete(d *schema.ResourceData, meta interface{}) error {
	return newResourceUser(d, meta).delete()
}

func (r *resourceUser) graphSet(user *msgraph.User) {
	r.resource.Set("user_principal_name", user.UserPrincipalName)
	r.resource.Set("display_name", user.DisplayName)
	r.resource.Set("given_name", user.GivenName)
	r.resource.Set("surname", user.Surname)
	r.resource.Set("mail_nickname", user.MailNickname)
	r.resource.Set("mail", user.Mail)
	r.resource.Set("other_mails", user.OtherMails)
	r.resource.Set("account_enabled", user.AccountEnabled)
	if user.PasswordProfile != nil {
		r.resource.Set("password", user.PasswordProfile.Password)
	}
}

func (r *resourceUser) graphGet() *msgraph.User {
	user := &msgraph.User{}
	if val, ok := r.resource.GetOkExists("user_principal_name"); ok {
		user.UserPrincipalName = P.CastString(val)
	}
	if val, ok := r.resource.GetOkExists("display_name"); ok {
		user.DisplayName = P.CastString(val)
	}
	if val, ok := r.resource.GetOkExists("given_name"); ok {
		user.GivenName = P.CastString(val)
	}
	if val, ok := r.resource.GetOkExists("surname"); ok {
		user.Surname = P.CastString(val)
	}
	if val, ok := r.resource.GetOkExists("mail_nickname"); ok {
		user.MailNickname = P.CastString(val)
	}
	for _, val := range r.resource.Get("other_mails").([]interface{}) {
		user.OtherMails = append(user.OtherMails, val.(string))
	}
	if val, ok := r.resource.GetOkExists("account_enabled"); ok {
		user.AccountEnabled = P.CastBool(val)
	}
	if r.resource.IsNewResource() || r.resource.HasChange("password") {
		user.PasswordProfile = &msgraph.PasswordProfile{}
		if val, ok := r.resource.GetOkExists("password"); ok {
			user.PasswordProfile.Password = P.CastString(val)
		} else {
			user.PasswordProfile.Password = P.String(uuid.New().String())
		}
		if val, ok := r.resource.GetOkExists("force_change_password_next_sign_in"); ok {
			user.PasswordProfile.ForceChangePasswordNextSignIn = P.CastBool(val)
		}
	}
	return user
}

func (r *resourceUser) create() error {
	newUser := r.graphGet()
	user, err := r.graph.userCreate(newUser)
	if err != nil {
		return err
	}
	user.PasswordProfile = newUser.PasswordProfile
	r.resource.SetId(*user.ID)
	r.graphSet(user)
	return nil
}

func (r *resourceUser) read() error {
	user, err := r.graph.userRead(r.resource.Id())
	if err != nil {
		r.resource.SetId("")
		if errRes, ok := err.(*msgraph.ErrorResponse); ok {
			if errRes.StatusCode() == http.StatusNotFound {
				return nil
			}
		}
		return err
	}
	r.graphSet(user)
	return nil
}

func (r *resourceUser) update() error {
	newUser := r.graphGet()
	user, err := r.graph.userUpdate(r.resource.Id(), newUser)
	if err != nil {
		return err
	}
	user.PasswordProfile = newUser.PasswordProfile
	r.graphSet(user)
	return nil
}

func (r *resourceUser) delete() error {
	r.graph.userDelete(r.resource.Id())
	r.resource.SetId("")
	return nil
}
