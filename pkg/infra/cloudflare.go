package infra

import (
	"fmt"

	"github.com/pulumi/pulumi-cloudflare/sdk/v3/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func ConfigureCloudflare(pCtx *pulumi.Context, instanceName string, publicIp pulumi.StringOutput, apiKey string, zoneId string) error {
	provider, err := cloudflare.NewProvider(pCtx, "cloudflareProvider", &cloudflare.ProviderArgs{
		ApiToken: pulumi.StringPtr(apiKey),
	})
	if err != nil {
		return fmt.Errorf("failed to create cloudflare provider: %w", err)
	}

	record, err := cloudflare.NewRecord(pCtx, fmt.Sprintf("cf-record-%s", instanceName), &cloudflare.RecordArgs{
		Name:    pulumi.String(instanceName),
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
