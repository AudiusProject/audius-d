package infra

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	ami          = "ami-05fb0b8c1424f266b" // Ubuntu, 22.04 LTS, amd64 jammy image build on 2023-12-07
	instanceType = "c5.2xlarge"
	volumeSize   = 100
)

func authProvider(pCtx *pulumi.Context) (provider *aws.Provider, err error) {
	if confCtxConfig == nil {
		return nil, fmt.Errorf("confCtxConfig is nil")
	}
	if confCtxConfig.Network.AWSAccessKeyID != "" && confCtxConfig.Network.AWSSecretAccessKey != "" && confCtxConfig.Network.AWSRegion != "" {
		provider, err = aws.NewProvider(pCtx, "aws", &aws.ProviderArgs{
			AccessKey: pulumi.String(confCtxConfig.Network.AWSAccessKeyID),
			SecretKey: pulumi.String(confCtxConfig.Network.AWSSecretAccessKey),
			Region:    pulumi.String(confCtxConfig.Network.AWSRegion),
		})
		return provider, err
	} else {
		return nil, fmt.Errorf("Invalid AWS creds")
	}
}

func CreateEC2Instance(pCtx *pulumi.Context, instanceName string) (*ec2.Instance, string, error) {

	awsProvider, err := authProvider(pCtx)
	if err != nil {
		return nil, "", fmt.Errorf("unable to authenticate aws provider: %w", err)
	}

	privateKeyFilePath, publicKeyPem, err := EnsureRSAKeyPair(instanceName)
	if err != nil {
		return nil, "", fmt.Errorf("unable to ensure RSA key pair: %w", err)
	}

	keyPair, err := ec2.NewKeyPair(pCtx, fmt.Sprintf("%s-keypair", instanceName), &ec2.KeyPairArgs{
		PublicKey: pulumi.String(publicKeyPem),
	}, pulumi.Provider(awsProvider))
	if err != nil {
		return nil, privateKeyFilePath, fmt.Errorf("unable to create key pair: %w", err)
	}

	userData := fmt.Sprintf(`#!/bin/bash
set -x
set -e

# install system level deps
sudo apt update
sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt update
sudo apt install -y docker-ce
sudo usermod -aG docker ubuntu

# install github cli (can be removed when audius-d repo no longer private)
# https://github.com/cli/cli/blob/trunk/docs/install_linux.md#debian-ubuntu-linux-raspberry-pi-os-apt
type -p curl >/dev/null || (sudo apt update && sudo apt install curl -y)
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg \
&& sudo chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg \
&& echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null \
&& sudo apt update \
&& sudo apt install gh -y

# download latest audius-d release
# TODO: remove PAT when audius-d repo public
echo "%s" | gh auth login --with-token
gh release download -R https://github.com/AudiusProject/audius-d --clobber --output ./audius-ctl --pattern audius-ctl-x86 && sudo mv ./audius-ctl /usr/local/bin/audius-ctl && sudo chmod +x /usr/local/bin/audius-ctl
# TODO: for now repo is required for devnet docker-compose. we should perhaps embed the compose file in the binary.
gh repo clone AudiusProject/audius-d /home/ubuntu/audius-d

# signal for successful completion
touch /home/ubuntu/user-data-done
`, confCtxConfig.Network.GHPat)

	instance, err := ec2.NewInstance(pCtx, fmt.Sprintf("%s-ec2-instance", instanceName), &ec2.InstanceArgs{
		Ami:          pulumi.String(ami),
		InstanceType: pulumi.String(instanceType),
		UserData:     pulumi.String(userData),
		KeyName:      keyPair.KeyName,
		Tags: pulumi.StringMap{
			"Name": pulumi.String(instanceName),
		},
		RootBlockDevice: &ec2.InstanceRootBlockDeviceArgs{
			VolumeType:          pulumi.String("gp3"),
			VolumeSize:          pulumi.Int(volumeSize),
			DeleteOnTermination: pulumi.Bool(true),
		},
	}, pulumi.Provider(awsProvider))
	if err != nil {
		return nil, privateKeyFilePath, fmt.Errorf("unable to create EC2 instance: %w", err)
	}

	pCtx.Export("instancePublicIp", instance.PublicIp)
	pCtx.Export("instancePrivateKeyFilePath", pulumi.String(privateKeyFilePath))

	return instance, privateKeyFilePath, nil
}
