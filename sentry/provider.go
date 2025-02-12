package sentry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown
}

// NewProvider returns a *schema.Provider.
func NewProvider(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"token": {
					Description: "The authentication token used to connect to Sentry. The value can be sourced from the `SENTRY_AUTH_TOKEN` environment variable.",
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SENTRY_AUTH_TOKEN", "SENTRY_TOKEN"}, nil),
					Sensitive:   true,
				},
				"base_url": {
					Description: "The target Sentry Base API URL follows the format `https://[hostname]/api/`, and this URL must end with the `/api/` path, including the trailing slash. The default value is `https://sentry.io/api/`, which proxies most requests to the appropriate region based on your organization. To avoid additional round trips, it is preferable to set the region URL explicitly. Use `https://us.sentry.io/api/` for the US region and `https://de.sentry.io/api/` for the EU region. This value is required for Sentry On-Premise deployments and can be sourced from the `SENTRY_BASE_URL` environment variable.",
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("SENTRY_BASE_URL", "https://sentry.io/api/"),
				},
			},

			ResourcesMap: map[string]*schema.Resource{
				"sentry_dashboard":                      resourceSentryDashboard(),
				"sentry_metric_alert":                   resourceSentryMetricAlert(),
				"sentry_organization_code_mapping":      resourceSentryOrganizationCodeMapping(),
				"sentry_organization_member":            resourceSentryOrganizationMember(),
				"sentry_organization_repository_github": resourceSentryOrganizationRepositoryGithub(),
				"sentry_organization":                   resourceSentryOrganization(),
				"sentry_plugin":                         resourceSentryPlugin(),
				"sentry_project":                        resourceSentryProject(),
				"sentry_team":                           resourceSentryTeam(),
			},

			DataSourcesMap: map[string]*schema.Resource{
				"sentry_dashboard":    dataSourceSentryDashboard(),
				"sentry_metric_alert": dataSourceSentryMetricAlert(),
				"sentry_organization": dataSourceSentryOrganization(),
				"sentry_team":         dataSourceSentryTeam(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		config := sentryclient.Config{
			UserAgent: p.UserAgent("terraform-provider-sentry", version),
			Token:     d.Get("token").(string),
			BaseURL:   d.Get("base_url").(string),
		}
		client, err := config.Client(ctx)

		if err != nil {
			return nil, diag.FromErr(err)
		}

		return client, nil
	}
}
