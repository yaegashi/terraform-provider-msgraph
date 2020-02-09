package main

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	P "github.com/yaegashi/msgraph.go/ptr"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func dataGroupResource() *schema.Resource {
	return &schema.Resource{
		Create: dataGroupCreate,
		Read:   dataGroupRead,
		Update: dataGroupUpdate,
		Delete: dataGroupDelete,
		Schema: map[string]*schema.Schema{
			"display_name":     &schema.Schema{Type: schema.TypeString, Required: true},
			"mail_nickname":    &schema.Schema{Type: schema.TypeString, Required: true},
			"mail_enabled":     &schema.Schema{Type: schema.TypeBool, Computed: true},
			"security_enabled": &schema.Schema{Type: schema.TypeBool, Computed: true},
			"mail":             &schema.Schema{Type: schema.TypeString, Computed: true},
			"group_types":      &schema.Schema{Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"visibility":       &schema.Schema{Type: schema.TypeString, Optional: true},
		},
	}
}

type dataGroup struct {
	ctx   context.Context
	graph *msgraph.GraphServiceRequestBuilder
	data  *schema.ResourceData
}

func newDataGroup(d *schema.ResourceData, m interface{}) *dataGroup {
	return &dataGroup{
		ctx:   context.Background(),
		graph: m.(*msgraph.GraphServiceRequestBuilder),
		data:  d,
	}
}

func dataGroupCreate(d *schema.ResourceData, m interface{}) error {
	return newDataGroup(d, m).dataCreate()
}

func dataGroupRead(d *schema.ResourceData, m interface{}) error {
	return newDataGroup(d, m).dataRead()
}

func dataGroupUpdate(d *schema.ResourceData, m interface{}) error {
	return newDataGroup(d, m).dataUpdate()
}

func dataGroupDelete(d *schema.ResourceData, m interface{}) error {
	return newDataGroup(d, m).dataDelete()
}

func (d *dataGroup) dataGraphSet(group *msgraph.Group) {
	d.data.Set("display_name", group.DisplayName)
	d.data.Set("mail_nickname", group.MailNickname)
	d.data.Set("mail_enabled", group.MailEnabled)
	d.data.Set("security_enabled", group.SecurityEnabled)
	d.data.Set("mail", group.Mail)
	d.data.Set("group_types", group.GroupTypes)
	d.data.Set("visibility", group.Visibility)
}

func (d *dataGroup) dataGraphGet() *msgraph.Group {
	group := &msgraph.Group{}
	if val, ok := d.data.GetOk("display_name"); ok {
		group.DisplayName = P.CastString(val)
	}
	if val, ok := d.data.GetOk("mail_nickname"); ok {
		group.MailNickname = P.CastString(val)
	}
	if val, ok := d.data.GetOk("visibility"); ok {
		group.Visibility = P.CastString(val)
	}
	for _, val := range d.data.Get("group_types").([]interface{}) {
		group.GroupTypes = append(group.GroupTypes, val.(string))
	}
	return group
}

func (d *dataGroup) dataCreate() error {
	newGroup := d.dataGraphGet()
	newGroup.MailEnabled = P.Bool(false)
	newGroup.SecurityEnabled = P.Bool(true)
	group, err := d.graphCreate(newGroup)
	if err != nil {
		return err
	}
	d.data.SetId(*group.ID)
	d.dataGraphSet(group)
	return nil
}

func (d *dataGroup) dataRead() error {
	group, err := d.graphRead(d.data.Id())
	if err != nil {
		d.data.SetId("")
		if errRes, ok := err.(*msgraph.ErrorResponse); ok {
			if errRes.StatusCode() == http.StatusNotFound {
				return nil
			}
		}
		return err
	}
	d.dataGraphSet(group)
	return nil
}

func (d *dataGroup) dataUpdate() error {
	newGroup := d.dataGraphGet()
	group, err := d.graphUpdate(d.data.Id(), newGroup)
	if err != nil {
		return err
	}
	d.dataGraphSet(group)
	return nil
}

func (d *dataGroup) dataDelete() error {
	d.graphDelete(d.data.Id())
	d.data.SetId("")
	return nil
}

func (d *dataGroup) graphCreate(group *msgraph.Group) (*msgraph.Group, error) {
	group, err := d.graph.Groups().Request().Add(d.ctx, group)
	if err != nil {
		return nil, err
	}
	return d.graphRead(*group.ID)
}

func (d *dataGroup) graphRead(id string) (*msgraph.Group, error) {
	req := d.graph.Groups().ID(id).Request()
	req.Select("id,displayName,mailNickname,mailEnabled,securityEnabled,mail,groupTypes,visibility")
	return req.Get(d.ctx)
}

func (d *dataGroup) graphUpdate(id string, group *msgraph.Group) (*msgraph.Group, error) {
	err := d.graph.Groups().ID(id).Request().Update(d.ctx, group)
	if err != nil {
		return nil, err
	}
	return d.graphRead(id)
}

func (d *dataGroup) graphDelete(id string) error {
	return d.graph.Groups().ID(id).Request().Delete(d.ctx)
}
