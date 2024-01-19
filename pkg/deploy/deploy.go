package deploy

import (
	"fmt"
	"os"
	"sync"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func DeployContext() {
	fmt.Println("DeployContext")

	// Create a WaitGroup
	var wg sync.WaitGroup

	// Add 1 to the WaitGroup counter
	wg.Add(1)

	pulumi.Run(func(ctx *pulumi.Context) error {

		defer wg.Done() // Decrease the WaitGroup counter when the function exits

		fmt.Println("foo")

		_, err := os.ReadFile("./ssh/id_rsa_aws.pub")
		if err != nil {
			return err
		}

		return nil
	})

	// Wait for the pulumi.Run goroutine to finish
	wg.Wait()
}
