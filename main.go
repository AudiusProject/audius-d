package main

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

//go:embed audius.conf
var confExample string
var confFilePath string
var imageTag string
var localImage bool
var port int
var tlsPort int

func main() {
	dc := ConnectDocker()
	Run(dc, "audius/dot-slash", "dev", "creator-node")
	// flag.StringVar(&confFilePath, "c", "", "Path to the .conf file")
	// flag.StringVar(&imageTag, "t", "dev", "docker image tag to use when turning up")
	// flag.BoolVar(&localImage, "local", false, "when specified, will use docker image from local repository")
	// flag.IntVar(&port, "port", 80, "specify a custom http port")
	// flag.IntVar(&tlsPort, "tls", 443, "specify a custom https port")

	// if !regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`).MatchString(imageTag) {
	// 	exitWithError("Invalid image tag:", imageTag)
	// }
	// cmdName := "up"
	// if len(os.Args) > 1 {
	// 	cmdName = os.Args[1]
	// }
	// flag.Parse()

	// switch cmdName {
	// case "down":
	// 	runDown()
	// default:
	// 	runUp(checkConfigFile())
	// }
}

func checkConfigFile() string {
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

func runUp(nodeType string) {
	ensureDirectory("/tmp/dind")

	ctx := context.Background()
	docker := ConnectDocker()

	audiusDotSlashImageName := "audius/dot-slash:" + imageTag

	if !localImage {
		reader, err := docker.ImagePull(ctx, audiusDotSlashImageName, types.ImagePullOptions{})
		if err != nil {
			exitWithError("Error pulling image:", err)
		}
		reader.Close()
	}

	// //volumeFlag := ""
	// if confFilePath != "" {
	// 	volumeFlag = fmt.Sprintf("-v %s:/root/audius-docker-compose/%s/override.env", confFilePath, nodeType)
	// }

	resp, err := docker.ContainerCreate(ctx, &container.Config{
		Image:   audiusDotSlashImageName,
		Volumes: map[string]struct{}{},
	}, nil, nil, nil, "creator-node")

	if err != nil {
		exitWithError("Creating creator-node container failed:", err)
	}

	if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		exitWithError(err)
	}

	fmt.Println("creator-node started")

	// var cmd string
	// baseCmd := fmt.Sprintf(`docker run --privileged -d -v /tmp/dind:/var/lib/docker %s -p %d:80 -p %d:443`, volumeFlag, port, tlsPort)

	// if nodeType == "creator-node" {
	// 	cmd = fmt.Sprintf(baseCmd + ` \
	//     --name creator-node \
	//     -v /var/k8s/mediorum:/var/k8s/mediorum \
	//     -v /var/k8s/creator-node-backend:/var/k8s/creator-node-backend \
	//     -v /var/k8s/creator-node-db:/var/k8s/creator-node-db \
	//     audius/dot-slash:` + imageTag)
	// } else {
	// 	cmd = fmt.Sprintf(baseCmd + ` \
	//     --name discovery-provider \
	//     -v /var/k8s/discovery-provider-db:/var/k8s/discovery-provider-db \
	//     -v /var/k8s/discovery-provider-chain:/var/k8s/discovery-provider-chain \
	//     audius/dot-slash:` + imageTag)
	// }

	// execCmd := fmt.Sprintf(`docker exec %s sh -c "while ! docker ps &> /dev/null; do echo 'starting up' && sleep 1; done && cd %s && docker compose up -d"`, nodeType, nodeType)

	// if err := runCommand("/bin/sh", "-c", cmd+" && "+execCmd); err != nil {
	// 	exitWithError("Error executing command:", err)
	// }
}

func runDown() {
	runCommand("docker", "rm", "-f", "creator-node", "discovery-provider")
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func exitWithError(msg ...interface{}) {
	fmt.Println(msg...)
	os.Exit(1)
}
