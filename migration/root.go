package migration

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	MigrateCmd = &cobra.Command{
		Use: "migrate <audius-docker-compose-path>",
		Short: `
		Nondestructive migration of audius-docker-compose to an audius-d context.
		Does not affect current selected context.
		`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := MigrateAudiusDockerCompose(args[0]); err != nil {
				log.Fatal("audius-docker-compose migration failed: ", err)
			}
			log.Println("audius-docker-compose migration successful 🎉")
		},
	}
)
