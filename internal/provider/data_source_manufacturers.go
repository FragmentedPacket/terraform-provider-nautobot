package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	nb "github.com/nautobot/go-nautobot/pkg/nautobot"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &manufacturersDataSource{}
	_ datasource.DataSourceWithConfigure = &manufacturersDataSource{}
)

// NewManufacturersDataSource is a helper function to simplify the provider implementation.
func NewManufacturersDataSource() datasource.DataSource {
	return &manufacturersDataSource{}
}

// manufacturersDataSource is the data source implementation.
type manufacturersDataSource struct {
	client *apiClient
}

func (d *manufacturersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*apiClient)
}

// Metadata returns the data source type name.
func (d *manufacturersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_manufacturers"
}

type manufacturersDataSourceModel struct {
	Manufacturers []manufacturersModel `tfsdk:"manufacturers"`
}

type manufacturersApiModel struct {
	ID                 string                 `tfsdk:"id" json:"id"`
	Created            string                 `tfsdk:"created" json:"created"`
	Description        string                 `tfsdk:"description" json:"description"`
	CustomFields       map[string]interface{} `tfsdk:"custom_fields" json:"custom_fields"`
	DeviceTypeCount    int64                  `tfsdk:"devicetype_count" json:"devicetype_count"`
	Display            string                 `tfsdk:"display" json:"display"`
	InventoryItemCount int64                  `tfsdk:"inventoryitem_count" json:"inventoryitem_count"`
	LastUpdated        string                 `tfsdk:"last_updated" json:"last_updated"`
	Name               string                 `tfsdk:"name" json:"name"`
	NotesUrl           string                 `tfsdk:"notes_url" json:"notes_url"`
	PlatformCount      int64                  `tfsdk:"platform_count" json:"platform_count"`
	Slug               string                 `tfsdk:"slug" json:"slug"`
	Url                string                 `tfsdk:"url" json:"url"`
}

type manufacturersModel struct {
	ID                 types.String     `tfsdk:"id" json:"id"`
	Created            types.String     `tfsdk:"created" json:"created"`
	Description        types.String     `tfsdk:"description" json:"description"`
	CustomFields       types.ObjectType `tfsdk:"custom_fields" json:"custom_fields"`
	DeviceTypeCount    types.Int64      `tfsdk:"devicetype_count" json:"devicetype_count"`
	Display            types.String     `tfsdk:"display" json:"display"`
	InventoryItemCount types.Int64      `tfsdk:"inventoryitem_count" json:"inventoryitem_count"`
	LastUpdated        types.String     `tfsdk:"last_updated" json:"last_updated"`
	Name               types.String     `tfsdk:"name" json:"name"`
	NotesUrl           types.String     `tfsdk:"notes_url" json:"notes_url"`
	PlatformCount      types.Int64      `tfsdk:"platform_count" json:"platform_count"`
	Slug               types.String     `tfsdk:"slug" json:"slug"`
	Url                types.String     `tfsdk:"url" json:"url"`
}

type apiResponse struct {
	Next     interface{}             `json:"next"`
	Previous interface{}             `json:"previous"`
	Count    int                     `json:"count"`
	Results  []manufacturersApiModel `json:"results"`
}

// Schema defines the schema for the data source.
func (d *manufacturersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"manufacturers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"created": schema.StringAttribute{
							Description: "Manufacturer's creation date.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Manufacturer's description.",
							Optional:    true,
						},
						"custom_fields": schema.MapAttribute{
							Description: "Manufacturer custom fields.",
							ElementType: schema.ObjectAttribute,
							Optional:    true,
							Computed:    true,
						},
						"devicetype_count": schema.Int64Attribute{
							Description: "Manufacturer's device count.",
							Computed:    true,
						},
						"display": schema.StringAttribute{
							Description: "Manufacturer's display name.",
							Optional:    true,
							Computed:    true,
						},
						"id": schema.StringAttribute{
							Description: "Manufacturer's UUID.",
							Computed:    true,
						},
						"inventoryitem_count": schema.Int64Attribute{
							Description: "Manufacturer's inventory item count.",
							Computed:    true,
						},
						"last_updated": schema.StringAttribute{
							Description: "Manufacturer's last update.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Manufacturer's name.",
							Required:    true,
						},
						"notes_url": schema.StringAttribute{
							Description: "Notes for manufacturer.",
							Optional:    true,
							Computed:    true,
						},
						"platform_count": schema.Int64Attribute{
							Description: "Manufacturer's platform count.",
							Computed:    true,
						},
						"slug": schema.StringAttribute{
							Description: "Manufacturer's slug.",
							Optional:    true,
							Computed:    true,
						},
						"url": schema.StringAttribute{
							Description: "Manufacturer's URL.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *manufacturersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state manufacturersDataSourceModel

	rsp, err := d.client.Client.DcimManufacturersListWithResponse(
		ctx,
		&nb.DcimManufacturersListParams{},
	)

	if err != nil {
		resp.Diagnostics.AddError("failed to get manufacturers list", err.Error())
		return
	}

	resultsReader := bytes.NewBufferString(string(rsp.Body))
	body, _ := ioutil.ReadAll(resultsReader)

	defer rsp.HTTPResponse.Body.Close()

	var response apiResponse
	// manufacturers := make([]manufacturersModel, 0)
	err = json.Unmarshal(body, &response)

	if err != nil {
		resp.Diagnostics.AddError("Failed to serialize", err.Error())
		return
	}

	// Map response body to model
	for _, manufacturer := range response.Results {
		// fmt.Println(manufacturer.CustomFields)
		manufacturerState := manufacturersModel{
			ID:                 types.StringValue(manufacturer.ID),
			Created:            types.StringValue(manufacturer.Created),
			Description:        types.StringValue(manufacturer.Description),
			CustomFields:       types.ObjectValue(manufacturer.CustomFields),
			DeviceTypeCount:    types.Int64Value(manufacturer.DeviceTypeCount),
			Display:            types.StringValue(manufacturer.Display),
			InventoryItemCount: types.Int64Value(manufacturer.InventoryItemCount),
			LastUpdated:        types.StringValue(manufacturer.LastUpdated),
			Name:               types.StringValue(manufacturer.Name),
			NotesUrl:           types.StringValue(manufacturer.NotesUrl),
			PlatformCount:      types.Int64Value(manufacturer.PlatformCount),
			Slug:               types.StringValue(manufacturer.Slug),
			Url:                types.StringValue(manufacturer.Url),
		}
		// customFields, diags := types.ObjectValue(ctx, types.StringType, manufacturer.CustomFields)
		// resp.Diagnostics.Append(diags...)
		// manufacturerState.CustomFields = customFields

		state.Manufacturers = append(state.Manufacturers, manufacturerState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
