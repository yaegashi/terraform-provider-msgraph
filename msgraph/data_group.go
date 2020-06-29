package msgraph

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func dataGroupResource() *schema.Resource {
	return &schema.Resource{
		Read: dataGroupRead,
		Schema: map[string]*schema.Schema{
			"id":               {Type: schema.TypeString, ValidateFunc: validation.IsUUID, Computed: true, Optional: true},
			"display_name":     {Type: schema.TypeString, Computed: true},
			"mail_nickname":    {Type: schema.TypeString, Computed: true, Optional: true},
			"mail_enabled":     {Type: schema.TypeBool, Computed: true},
			"security_enabled": {Type: schema.TypeBool, Computed: true},
			"mail":             {Type: schema.TypeString, Computed: true},
			"group_types":      {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"visibility":       {Type: schema.TypeString, Computed: true},
		},
	}
}

type dataGroup struct {
	graph *graph
	data  *schema.ResourceData
}

func newDataGroup(d *schema.ResourceData, m interface{}) *dataGroup {
	return &dataGroup{
		graph: newGraph(m),
		data:  d,
	}
}

func dataGroupRead(d *schema.ResourceData, m interface{}) error {
	return newDataGroup(d, m).read()
}

func (d *dataGroup) graphSet(group *msgraph.Group) {
	d.data.Set("display_name", group.DisplayName)
	d.data.Set("mail_nickname", group.MailNickname)
	d.data.Set("mail_enabled", group.MailEnabled)
	d.data.Set("security_enabled", group.SecurityEnabled)
	d.data.Set("mail", group.Mail)
	d.data.Set("group_types", group.GroupTypes)
	d.data.Set("visibility", group.Visibility)
}

func (d *dataGroup) read() error {
	var (
		group *msgraph.Group
		err   error
	)
	if id, ok := d.data.GetOkExists("id"); ok {
		group, err = d.graph.groupRead(id.(string))
	} else if nick, ok := d.data.GetOkExists("mail_nickname"); ok {
		group, err = d.graph.groupReadByMailNickname(nick.(string))
	} else {
		err = fmt.Errorf("one of `id` and `mail_nickname` must be supplied")
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
	d.data.SetId(*group.ID)
	d.graphSet(group)
	return nil
}
