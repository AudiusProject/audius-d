package infra

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	ami          = "ami-05fb0b8c1424f266b" // Ubuntu, 22.04 LTS, amd64 jammy image build on 2023-12-07
	instanceType = "t2.medium"
	volumeSize   = 100
)

func CreateEC2Instance(ctx *pulumi.Context, instanceName string) (*ec2.Instance, string, error) {
	privateKeyFilePath, publicKeyPem, err := EnsureRSAKeyPair(instanceName)
	if err != nil {
		return nil, "", fmt.Errorf("unable to ensure RSA key pair: %w", err)
	}

	keyPair, err := ec2.NewKeyPair(ctx, fmt.Sprintf("%s-keypair", instanceName), &ec2.KeyPairArgs{
		PublicKey: pulumi.String(publicKeyPem),
	})
	if err != nil {
		return nil, privateKeyFilePath, fmt.Errorf("unable to create key pair: %w", err)
	}

	userData := `#!/bin/bash
# install system level deps
sudo apt update
sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt update
sudo apt install -y docker-ce
sudo usermod -aG docker $USER && newgrp docker

# install github cli (can be removed when audius-d repo no longer private)
# https://github.com/cli/cli/blob/trunk/docs/install_linux.md#debian-ubuntu-linux-raspberry-pi-os-apt
type -p curl >/dev/null || (sudo apt update && sudo apt install curl -y)
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg \
&& sudo chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg \
&& echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null \
&& sudo apt update \
&& sudo apt install gh -y

# download latest audius-d release
# todo mount secret
# echo <github_pat_xxxx> | gh auth login --with-token
# gh release download -R https://github.com/AudiusProject/audius-d --clobber --output ./audius-ctl --pattern audius-ctl-x86 && sudo mv ./audius-ctl /usr/local/bin/audius-ctl && sudo chmod +x /usr/local/bin/audius-ctl
`

	instance, err := ec2.NewInstance(ctx, fmt.Sprintf("%s-ec2-instance", instanceName), &ec2.InstanceArgs{
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
	})
	if err != nil {
		return nil, privateKeyFilePath, fmt.Errorf("unable to create EC2 instance: %w", err)
	}

	return instance, privateKeyFilePath, nil
}
