package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hasura/go-graphql-client"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
)

// initGraphQLClient initializes the graphql client.
func initGraphQLClient(client api.Client) {
	graphQLURL := fmt.Sprintf("https://api.%s/v3/graphql", client.Host)
	bearerToken := fmt.Sprintf("Bearer %s", client.AccessToken)
	api.GraphQLClient = graphql.NewClient(graphQLURL, nil).WithRequestModifier(func(req *http.Request) {
		req.Header.Set("Authorization", bearerToken)
	})
}

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			DataSourcesMap: map[string]*schema.Resource{
				"squadcast_squad":             dataSourceSquad(),
				"squadcast_service":           dataSourceService(),
				"squadcast_escalation_policy": dataSourceEscalationPolicy(),
				// "squadcast_teams": dataSourceTeams(),
				"squadcast_team":        dataSourceTeam(),
				"squadcast_team_role":   dataSourceTeamRole(),
				"squadcast_user":        dataSourceUser(),
				"squadcast_schedule":    dataSourceSchedule(),
				"squadcast_schedule_v2": dataSourceScheduleV2(),
				"squadcast_runbook":     dataSourceRunbook(),
				"squadcast_webform":     dataSourceWebform(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"squadcast_apta_config":                  resourceAPTAConfig(),
				"squadcast_custom_content_template":      resourceCustomContentTemplate(),
				"squadcast_deduplication_rules":          resourceDeduplicationRules(),
				"squadcast_deduplication_rule_v2":        resourceDeduplicationRuleV2(),
				"squadcast_delayed_notification_config":  resourceDelayedNotificationConfig(),
				"squadcast_dedup_key_overlay":            resourceDedupKeyOverlay(),
				"squadcast_escalation_policy":            resourceEscalationPolicy(),
				"squadcast_ger":                          resourceGER(),
				"squadcast_ger_ruleset":                  resourceGERRuleset(),
				"squadcast_ger_ruleset_rule":             resourceGERRulesetRule(),
				"squadcast_ger_ruleset_rules_ordering":   resourceGERRulesetRulesOrdering(),
				"squadcast_global_oncall_reminder_rules": resourceGlobalOncallReminderRules(),
				"squadcast_iag_config":                   resourceIAGConfig(),
				"squadcast_routing_rules":                resourceRoutingRules(),
				"squadcast_routing_rule_v2":              resourceRoutingRuleV2(),
				"squadcast_runbook":                      resourceRunbook(),
				"squadcast_schedule":                     resourceSchedule(),
				"squadcast_schedule_v2":                  resourceScheduleV2(),
				"squadcast_schedule_rotation_v2":         resourceScheduleRotationV2(),
				"squadcast_service_maintenance":          resourceServiceMaintenance(),
				"squadcast_service":                      resourceService(),
				"squadcast_squad":                        resourceSquad(),
				"squadcast_status_page":                  resourceStatusPage(),
				"squadcast_status_page_component":        resourceStatusPageComponent(),
				"squadcast_status_page_group":            resourceStatusPageGroup(),
				"squadcast_suppression_rules":            resourceSuppressionRules(),
				"squadcast_suppression_rule_v2":          resourceSuppressionRuleV2(),
				"squadcast_tagging_rules":                resourceTaggingRules(),
				"squadcast_tagging_rule_v2":              resourceTaggingRuleV2(),
				"squadcast_team_member":                  resourceTeamMember(),
				"squadcast_team_role":                    resourceTeamRole(),
				"squadcast_team":                         resourceTeam(),
				"squadcast_user":                         resourceUser(),
				"squadcast_slo":                          resourceSlo(),
				"squadcast_webform":                      resourceWebform(),
				"squadcast_workflow":                     resourceWorkflow(),
				"squadcast_workflow_action":              resourceWorkflowAction(),
				"squadcast_workflow_action_ordering":     resourceWorkflowActionOrdering(),
			},
			Schema: map[string]*schema.Schema{
				"region": {
					Description: "The region you are currently hosted on." +
						"Supported values are \"us\" and \"eu\"",
					Type:         schema.TypeString,
					Optional:     true,
					DefaultFunc:  schema.EnvDefaultFunc("SQUADCAST_REGION", "us"),
					ValidateFunc: validation.StringInSlice([]string{"us", "eu", "internal", "staging", "dev"}, false),
				},
				"refresh_token": {
					Description: "The refresh token, This can be created from user profile",
					Type:        schema.TypeString,
					Sensitive:   true,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("SQUADCAST_REFRESH_TOKEN", nil),
				},
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (any, diag.Diagnostics) {
	return func(ctx context.Context, rd *schema.ResourceData) (c any, diags diag.Diagnostics) {
		client := &api.Client{}
		client.UserAgent = p.UserAgent("terraform-provider-squadcast", version)

		region := rd.Get("region").(string)
		refreshToken := rd.Get("refresh_token").(string)

		if refreshToken == "" {
			refreshToken = os.Getenv("SQUADCAST_REFRESH_TOKEN")
		}
		if refreshToken == "" {
			return nil, diag.Errorf("refresh_token is required")
		}

		client.RefreshToken = refreshToken

		switch region {
		case "us":
			client.Host = "squadcast.com"
		case "eu":
			client.Host = "eu.squadcast.com"
		case "internal":
			client.Host = "squadcast.xyz"
		case "staging":
			client.Host = "squadcast.tech"
		case "dev":
			client.Host = "localhost"
		}

		if region == "dev" {
			client.BaseURLV4 = fmt.Sprintf("http://%s:8081/v4", client.Host)
			client.BaseURLV3 = fmt.Sprintf("http://%s:8081/v3", client.Host)
			client.AuthBaseURL = fmt.Sprintf("http://%s:8081/v3", client.Host)
			client.IngestionBaseURL = fmt.Sprintf("http://%s:8458", client.Host)
		} else {
			client.BaseURLV4 = fmt.Sprintf("https://api.%s/v4", client.Host)
			client.BaseURLV3 = fmt.Sprintf("https://api.%s/v3", client.Host)
			client.AuthBaseURL = fmt.Sprintf("https://api.%s/v3", client.Host)
			client.IngestionBaseURL = fmt.Sprintf("https://api.%s", client.Host)
		}

		token, err := client.GetAccessToken(ctx)
		if err != nil {
			return nil, append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "An error occurred while fetching the access token.",
				Detail:   err.Error(),
			})
		}
		client.AccessToken = token.AccessToken

		org, err := client.GetCurrentOrganization(ctx)
		if err != nil {
			return nil, append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "An error occurred while fetching the organization.",
				Detail:   err.Error(),
			})
		}
		client.OrganizationID = org.ID

		initGraphQLClient(*client)

		return client, nil
	}
}
