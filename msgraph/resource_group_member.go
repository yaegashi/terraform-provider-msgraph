package msgraph

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func resourceGroupMemberResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupMemberCreate,
		Read:   resourceGroupMemberRead,
		Delete: resourceGroupMemberDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_id":  {Type: schema.TypeString, Required: true, ForceNew: true},
			"member_id": {Type: schema.TypeString, Required: true, ForceNew: true},
		},
	}
}

type resourceGroupMember struct {
	graph    *graph
	resource *schema.ResourceData
}

func newResourceGroupMember(r *schema.ResourceData, m interface{}) *resourceGroupMember {
	return &resourceGroupMember{
		graph:    newGraph(m),
		resource: r,
	}
}

func resourceGroupMemberCreate(r *schema.ResourceData, m interface{}) error {
	return newResourceGroupMember(r, m).create()
}

func resourceGroupMemberRead(r *schema.ResourceData, m interface{}) error {
	return newResourceGroupMember(r, m).read()
}

func resourceGroupMemberDelete(r *schema.ResourceData, m interface{}) error {
	return newResourceGroupMember(r, m).delete()
}

func (r *resourceGroupMember) create() error {
	groupID := r.resource.Get("group_id").(string)
	memberID := r.resource.Get("member_id").(string)
	reqObj := map[string]interface{}{
		"@odata.id": r.graph.cli.DirectoryObjects().ID(memberID).Request().URL(),
	}
	req := r.graph.cli.Groups().ID(groupID).Members().Request()
	err := req.JSONRequest(r.graph.ctx, "POST", "/$ref", reqObj, nil)
	if err != nil {
		return err
	}
	r.resource.SetId(fmt.Sprintf("%s:%s", groupID, memberID))
	return nil
}

func (r *resourceGroupMember) read() error {
	id := r.resource.Id()
	s := strings.Split(id, ":")
	if len(s) != 2 {
		return fmt.Errorf("Unable to parse ID: %s", id)
	}
	groupID := s[0]
	memberID := s[1]
	req := r.graph.cli.Groups().ID(groupID).Members().ID(memberID).Request()
	err := req.JSONRequest(r.graph.ctx, "GET", "/$ref", nil, nil)
	if err != nil {
		r.resource.SetId("")
		if errRes, ok := err.(*msgraph.ErrorResponse); ok {
			if errRes.StatusCode() == http.StatusNotFound {
				return nil
			}
		}
		return err
	}
	r.resource.Set("group_id", groupID)
	r.resource.Set("member_id", memberID)
	return nil
}

func (r *resourceGroupMember) delete() error {
	groupID := r.resource.Get("group_id").(string)
	memberID := r.resource.Get("member_id").(string)
	req := r.graph.cli.Groups().ID(groupID).Members().ID(memberID).Request()
	err := req.JSONRequest(r.graph.ctx, "DELETE", "/$ref", nil, nil)
	if err != nil {
		return err
	}
	r.resource.SetId("")
	return nil
}
