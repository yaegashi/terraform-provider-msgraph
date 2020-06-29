package msgraph

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	P "github.com/yaegashi/msgraph.go/ptr"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func resourceGroupResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"display_name":     {Type: schema.TypeString, Required: true},
			"mail_nickname":    {Type: schema.TypeString, Required: true},
			"mail_enabled":     {Type: schema.TypeBool, Computed: true},
			"security_enabled": {Type: schema.TypeBool, Computed: true},
			"mail":             {Type: schema.TypeString, Computed: true},
			"group_types":      {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"visibility":       {Type: schema.TypeString, Optional: true},
		},
	}
}

type resourceGroup struct {
	graph    *graph
	resource *schema.ResourceData
}

func newResourceGroup(r *schema.ResourceData, m interface{}) *resourceGroup {
	return &resourceGroup{
		graph:    newGraph(m),
		resource: r,
	}
}

func resourceGroupCreate(r *schema.ResourceData, m interface{}) error {
	return newResourceGroup(r, m).create()
}

func resourceGroupRead(r *schema.ResourceData, m interface{}) error {
	return newResourceGroup(r, m).read()
}

func resourceGroupUpdate(r *schema.ResourceData, m interface{}) error {
	return newResourceGroup(r, m).update()
}

func resourceGroupDelete(r *schema.ResourceData, m interface{}) error {
	return newResourceGroup(r, m).delete()
}

func (r *resourceGroup) graphSet(group *msgraph.Group) {
	r.resource.Set("display_name", group.DisplayName)
	r.resource.Set("mail_nickname", group.MailNickname)
	r.resource.Set("mail_enabled", group.MailEnabled)
	r.resource.Set("security_enabled", group.SecurityEnabled)
	r.resource.Set("mail", group.Mail)
	r.resource.Set("group_types", group.GroupTypes)
	r.resource.Set("visibility", group.Visibility)
}

func (r *resourceGroup) graphGet() *msgraph.Group {
	group := &msgraph.Group{}
	if val, ok := r.resource.GetOk("display_name"); ok {
		group.DisplayName = P.CastString(val)
	}
	if val, ok := r.resource.GetOk("mail_nickname"); ok {
		group.MailNickname = P.CastString(val)
	}
	if val, ok := r.resource.GetOk("visibility"); ok {
		group.Visibility = P.CastString(val)
	}
	for _, val := range r.resource.Get("group_types").([]interface{}) {
		group.GroupTypes = append(group.GroupTypes, val.(string))
	}
	return group
}

func (r *resourceGroup) create() error {
	newGroup := r.graphGet()
	newGroup.MailEnabled = P.Bool(false)
	newGroup.SecurityEnabled = P.Bool(true)
	group, err := r.graph.groupCreate(newGroup)
	if err != nil {
		return err
	}
	r.resource.SetId(*group.ID)
	r.graphSet(group)
	return nil
}

func (r *resourceGroup) read() error {
	group, err := r.graph.groupRead(r.resource.Id())
	if err != nil {
		r.resource.SetId("")
		if errRes, ok := err.(*msgraph.ErrorResponse); ok {
			if errRes.StatusCode() == http.StatusNotFound {
				return nil
			}
		}
		return err
	}
	r.graphSet(group)
	return nil
}

func (r *resourceGroup) update() error {
	newGroup := r.graphGet()
	group, err := r.graph.groupUpdate(r.resource.Id(), newGroup)
	if err != nil {
		return err
	}
	r.graphSet(group)
	return nil
}

func (r *resourceGroup) delete() error {
	r.graph.groupDelete(r.resource.Id())
	r.resource.SetId("")
	return nil
}
