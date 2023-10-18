package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

//go:embed audius.conf
var confExample string
var confFilePath string
var imageTag string
var localImage bool
var port int
var tlsPort int

func main() {
	flag.StringVar(&confFilePath, "c", "", "Path to the .conf file")
	flag.StringVar(&imageTag, "t", "dev", "docker image tag to use when turning up")
	flag.BoolVar(&localImage, "local", false, "when specified, will use docker image from local repository")
	flag.IntVar(&port, "port", 80, "specify a custom http port")
	flag.IntVar(&tlsPort, "tls", 443, "specify a custom https port")

	if !regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`).MatchString(imageTag) {
		exitWithError("Invalid image tag:", imageTag)
	}
	cmdName := "up"
	if len(os.Args) > 1 {
		cmdName = os.Args[1]
	}
	flag.Parse()

	dc := ConnectDocker()

	switch cmdName {
	case "down":
		DownAll(dc)
	default:
		runUp(dc)
	}
}

func determineNodeType() string {
	nodeType := "discovery-provider"

	if confFilePath == "" {
		if usr, err := user.Current(); err != nil {
			exitWithError("Error retrieving current user:", err)
		} else {
			confFilePath = filepath.Join(usr.HomeDir, ".audius", "audius.conf")
		}
	} else {
		if absPath, err := filepath.Abs(confFilePath); err != nil {
			exitWithError("Error creating absolute path to config file:", err)
		} else {
			confFilePath = absPath
		}
	}

	if _, err := os.Stat(confFilePath); os.IsNotExist(err) {
		exitWithError("Config not found at provided location:", confFilePath, confExample)
	}

	file, err := os.Open(confFilePath)
	if err != nil {
		exitWithError("Error opening config file:", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "creatorNodeEndpoint") {
			nodeType = "creator-node"
			break
		}
	}
	if err := scanner.Err(); err != nil {
		exitWithError("Error reading config file:", err)
	}

	return nodeType
}

func runUp(dc *client.Client) {
	nodeType := determineNodeType()
	nodeConf := &container.Config{}

	mounts := []mount.Mount{}

	// separate config per node type
	if nodeType == "creator-node" {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: "/var/k8s/mediorum",
			Target: "/var/k8s/mediorum",
		})
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: "/var/k8s/creator-node-backend",
			Target: "/var/k8s/creator-node-backend",
		})
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: "/var/k8s/creator-node-db",
			Target: "/var/k8s/creator-node-db",
		})
	}

	if nodeType == "discovery-provider" {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: "/var/k8s/discovery-provider-db",
			Target: "/var/k8s/discovery-provider-db",
		})
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: "/var/k8s/discovery-provider-chain",
			Target: "/var/k8s/discovery-provider-chain",
		})
	}

	Run(dc, "audius/dot-slash", imageTag, nodeType, nodeConf, mounts)
}

func exitWithError(msg ...interface{}) {
	fmt.Println(msg...)
	os.Exit(1)
}
