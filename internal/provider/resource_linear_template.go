package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLinearTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLinearTemplateCreate,
		ReadContext:   resourceLinearTemplateRead,
		UpdateContext: resourceLinearTemplateUpdate,
		DeleteContext: resourceLinearTemplateDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the template",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the template",
			},
			"template_data": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The template data in JSON format",
			},
			"team_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the team this template belongs to",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "issue",
				Description: "The type of the template",
			},
			"sort_order": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "The sort order of the template",
			},
		},
	}
}

func resourceLinearTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*LinearClient)

	name := d.Get("name").(string)
	teamID := d.Get("team_id").(string)
	templateData := d.Get("template_data").(string)
	templateType := d.Get("type").(string)

	// Prepare variables for GraphQL mutation
	variables := map[string]interface{}{
		"templateCreateInput": map[string]interface{}{
			"name":         name,
			"teamId":       teamID,
			"templateData": templateData,
			"type":         templateType,
		},
	}

	// Add optional fields if present
	if description, ok := d.GetOk("description"); ok {
		variables["templateCreateInput"].(map[string]interface{})["description"] = description.(string)
	}

	if sortOrder, ok := d.GetOk("sort_order"); ok {
		variables["templateCreateInput"].(map[string]interface{})["sortOrder"] = sortOrder.(float64)
	}

	mutation := `
		mutation TemplateCreate($templateCreateInput: TemplateCreateInput!) {
			templateCreate(input: $templateCreateInput) {
				success
				template {
					id
				}
			}
		}
	`

	data, err := client.GraphQLRequest(mutation, variables)
	if err != nil {
		return diag.FromErr(err)
	}

	// Extract the template ID from the response
	var response struct {
		TemplateCreate struct {
			Success  bool
			Template struct {
				ID string
			}
		}
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return diag.FromErr(err)
	}

	if !response.TemplateCreate.Success {
		return diag.Errorf("Failed to create Linear template")
	}

	d.SetId(response.TemplateCreate.Template.ID)

	return resourceLinearTemplateRead(ctx, d, m)
}

func resourceLinearTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*LinearClient)
	var diags diag.Diagnostics

	templateID := d.Id()

	query := `
		query Template($templateId: String!) {
			template(id: $templateId) {
				id
				name
				description
				templateData
				type
				sortOrder
				team {
					id
				}
			}
		}
	`

	variables := map[string]interface{}{
		"templateId": templateID,
	}

	data, err := client.GraphQLRequest(query, variables)
	if err != nil {
		return diag.FromErr(err)
	}

	var response struct {
		Template struct {
			ID           string  `json:"id"`
			Name         string  `json:"name"`
			Description  string  `json:"description"`
			TemplateData string  `json:"templateData"`
			Type         string  `json:"type"`
			SortOrder    float64 `json:"sortOrder"`
			Team         struct {
				ID string `json:"id"`
			} `json:"team"`
		} `json:"template"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return diag.FromErr(err)
	}

	if response.Template.ID == "" {
		d.SetId("")
		return diags
	}

	d.Set("name", response.Template.Name)
	d.Set("description", response.Template.Description)
	d.Set("template_data", response.Template.TemplateData)
	d.Set("team_id", response.Template.Team.ID)
	d.Set("type", response.Template.Type)
	d.Set("sort_order", response.Template.SortOrder)

	return diags
}

func resourceLinearTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*LinearClient)

	templateID := d.Id()

	// Prepare variables for GraphQL mutation
	updateInput := map[string]interface{}{
		"id": templateID,
	}

	if d.HasChange("name") {
		updateInput["name"] = d.Get("name").(string)
	}

	if d.HasChange("description") {
		updateInput["description"] = d.Get("description").(string)
	}

	if d.HasChange("template_data") {
		updateInput["templateData"] = d.Get("template_data").(string)
	}

	if d.HasChange("team_id") {
		updateInput["teamId"] = d.Get("team_id").(string)
	}

	if d.HasChange("type") {
		updateInput["type"] = d.Get("type").(string)
	}

	if d.HasChange("sort_order") {
		updateInput["sortOrder"] = d.Get("sort_order").(float64)
	}

	variables := map[string]interface{}{
		"templateUpdateInput": updateInput,
	}

	mutation := `
		mutation TemplateUpdate($templateUpdateInput: TemplateUpdateInput!) {
			templateUpdate(input: $templateUpdateInput) {
				success
			}
		}
	`

	data, err := client.GraphQLRequest(mutation, variables)
	if err != nil {
		return diag.FromErr(err)
	}

	var response struct {
		TemplateUpdate struct {
			Success bool
		}
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return diag.FromErr(err)
	}

	if !response.TemplateUpdate.Success {
		return diag.Errorf("Failed to update Linear template")
	}

	return resourceLinearTemplateRead(ctx, d, m)
}

func resourceLinearTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*LinearClient)
	var diags diag.Diagnostics

	templateID := d.Id()

	variables := map[string]interface{}{
		"templateDeleteInput": map[string]interface{}{
			"id": templateID,
		},
	}

	mutation := `
		mutation TemplateDelete($templateDeleteInput: TemplateDeleteInput!) {
			templateDelete(input: $templateDeleteInput) {
				success
			}
		}
	`

	data, err := client.GraphQLRequest(mutation, variables)
	if err != nil {
		return diag.FromErr(err)
	}

	var response struct {
		TemplateDelete struct {
			Success bool
		}
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return diag.FromErr(err)
	}

	if !response.TemplateDelete.Success {
		return diag.Errorf("Failed to delete Linear template")
	}

	d.SetId("")

	return diags
} 