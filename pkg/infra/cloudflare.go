package infra

import (
	"fmt"

	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/pulumi/pulumi-cloudflare/sdk/v3/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func cloudflareCredentialsValid(networkConfig *conf.NetworkConfig) bool {
	if networkConfig != nil && networkConfig.Infra != nil {
		return networkConfig.Infra.CloudflareAPIKey != "" && networkConfig.Infra.CloudflareZoneId != ""
	}
	return false
}

func cloudflareAuthProvider(pCtx *pulumi.Context) (*cloudflare.Provider, error) {
	if cloudflareCredentialsValid(&confCtxConfig.Network) {
		provider, err := cloudflare.NewProvider(pCtx, "cloudflareProvider", &cloudflare.ProviderArgs{
			ApiToken: pulumi.StringPtr(confCtxConfig.Network.Infra.CloudflareAPIKey),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create cloudflare provider: %w", err)
		}
		return provider, nil
	}
	return nil, fmt.Errorf("invalid CloudflareCredentials")
}

func CloudflareAddDNSRecord(pCtx *pulumi.Context, provider *cloudflare.Provider, zoneId string, name string, publicIp pulumi.StringOutput) error {

	record, err := cloudflare.NewRecord(pCtx, fmt.Sprintf("cf-record-%s", name), &cloudflare.RecordArgs{
		Name:    pulumi.String(name),
		Proxied: pulumi.Bool(true),
		Ttl:     pulumi.Int(1), // Set TTL to automatic (required for proxied)
		Type:    pulumi.String("A"),
		Value:   publicIp,
		ZoneId:  pulumi.String(zoneId),
	}, pulumi.Provider(provider))
	if err != nil {
		return fmt.Errorf("failed to create cloudflare record: %w", err)
	}

	pCtx.Export("cloudflareRecordHostname", record.Hostname)
	pCtx.Export("cloudflareRecordValue", record.Value)

	return nil
}
