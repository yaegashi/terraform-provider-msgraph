package msgraph

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	P "github.com/yaegashi/msgraph.go/ptr"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func resourceApplicationResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationCreate,
		Read:   resourceApplicationRead,
		Update: resourceApplicationUpdate,
		Delete: resourceApplicationDelete,
		Schema: map[string]*schema.Schema{
			"app_id":                       &schema.Schema{Type: schema.TypeString, Computed: true},
			"display_name":                 &schema.Schema{Type: schema.TypeString, Required: true},
			"sign_in_audience":             &schema.Schema{Type: schema.TypeString, Optional: true},
			"identifier_uris":              &schema.Schema{Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"home_page_url":                &schema.Schema{Type: schema.TypeString, Optional: true},
			"logout_url":                   &schema.Schema{Type: schema.TypeString, Optional: true},
			"redirect_uris":                &schema.Schema{Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"enable_id_token_issuance":     &schema.Schema{Type: schema.TypeBool, Optional: true},
			"enable_access_token_issuance": &schema.Schema{Type: schema.TypeBool, Optional: true},
		},
	}
}

type resourceApplication struct {
	graph    *graph
	resource *schema.ResourceData
}

func newResourceApplication(r *schema.ResourceData, m interface{}) *resourceApplication {
	return &resourceApplication{
		graph:    newGraph(m),
		resource: r,
	}
}

func resourceApplicationCreate(r *schema.ResourceData, m interface{}) error {
	return newResourceApplication(r, m).create()
}

func resourceApplicationRead(r *schema.ResourceData, m interface{}) error {
	return newResourceApplication(r, m).read()
}

func resourceApplicationUpdate(r *schema.ResourceData, m interface{}) error {
	return newResourceApplication(r, m).update()
}

func resourceApplicationDelete(r *schema.ResourceData, m interface{}) error {
	return newResourceApplication(r, m).delete()
}

func (r *resourceApplication) graphSet(application *msgraph.Application) {
	r.resource.Set("app_id", application.AppID)
	r.resource.Set("display_name", application.DisplayName)
	r.resource.Set("sign_in_audience", application.SignInAudience)
	r.resource.Set("identifier_uris", application.IdentifierUris)
	r.resource.Set("home_page_url", application.Web.HomePageURL)
	r.resource.Set("logout_url", application.Web.LogoutURL)
	r.resource.Set("redirect_uris", application.Web.RedirectUris)
	r.resource.Set("enable_id_token_issuance", application.Web.ImplicitGrantSettings.EnableIDTokenIssuance)
	r.resource.Set("enable_access_token_issuance", application.Web.ImplicitGrantSettings.EnableAccessTokenIssuance)
}

func (r *resourceApplication) graphGet() *msgraph.Application {
	application := &msgraph.Application{Web: &msgraph.WebApplication{ImplicitGrantSettings: &msgraph.ImplicitGrantSettings{}}}
	if val, ok := r.resource.GetOkExists("app_id"); ok {
		application.AppID = P.CastString(val)
	}
	if val, ok := r.resource.GetOkExists("display_name"); ok {
		application.DisplayName = P.CastString(val)
	}
	if val, ok := r.resource.GetOkExists("sign_in_audience"); ok {
		application.SignInAudience = P.CastString(val)
	}
	for _, val := range r.resource.Get("identifier_uris").([]interface{}) {
		application.IdentifierUris = append(application.IdentifierUris, val.(string))
	}
	if val, ok := r.resource.GetOkExists("home_page_url"); ok {
		application.Web.HomePageURL = P.CastString(val)
	}
	if val, ok := r.resource.GetOkExists("logout_url"); ok {
		application.Web.LogoutURL = P.CastString(val)
	}
	for _, val := range r.resource.Get("redirect_uris").([]interface{}) {
		application.Web.RedirectUris = append(application.Web.RedirectUris, val.(string))
	}
	if val, ok := r.resource.GetOkExists("enable_id_token_issuance"); ok {
		application.Web.ImplicitGrantSettings.EnableIDTokenIssuance = P.CastBool(val)
	}
	if val, ok := r.resource.GetOkExists("enable_access_token_issuance"); ok {
		application.Web.ImplicitGrantSettings.EnableAccessTokenIssuance = P.CastBool(val)
	}
	return application
}

func (r *resourceApplication) create() error {
	newApplication := r.graphGet()
	application, err := r.graph.applicationCreate(newApplication)
	if err != nil {
		return err
	}
	r.resource.SetId(*application.ID)
	r.graphSet(application)
	return nil
}

func (r *resourceApplication) read() error {
	application, err := r.graph.applicationRead(r.resource.Id())
	if err != nil {
		r.resource.SetId("")
		if errRes, ok := err.(*msgraph.ErrorResponse); ok {
			if errRes.StatusCode() == http.StatusNotFound {
				return nil
			}
		}
		return err
	}
	r.graphSet(application)
	return nil
}

func (r *resourceApplication) update() error {
	newApplication := r.graphGet()
	application, err := r.graph.applicationUpdate(r.resource.Id(), newApplication)
	if err != nil {
		return err
	}
	r.graphSet(application)
	return nil
}

func (r *resourceApplication) delete() error {
	r.graph.applicationDelete(r.resource.Id())
	r.resource.SetId("")
	return nil
}
