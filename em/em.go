package em

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strings"

	"github.com/AudiusProject/audius-d/acdc"
	"github.com/AudiusProject/audius-d/hashes"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

var EmCmd *cobra.Command

const (
	keyfileLocation = ".audius_developer_app_private_key"
)

func init() {

	var actorId string
	var action string
	var entityType string
	var entityId string

	EmCmd = &cobra.Command{
		Use:   "em",
		Short: "Send EntityManager transactions",
		Run: func(cmd *cobra.Command, args []string) {
			// todo: get endpoint from env (stage / prod / other)
			const ACDC_ENDPOINT = `https://acdc-gateway.staging.audius.co`
			client, err := ethclient.Dial(ACDC_ENDPOINT)
			if err != nil {
				log.Fatal(err)
			}

			sendEmTx(client, actorId, action, entityType, entityId)
		},
	}

	EmCmd.PersistentFlags().StringVar(&actorId, "actor", "", "user id performing the action")
	EmCmd.PersistentFlags().StringVar(&action, "action", "", "verb to perform: Repost / Save / Follow")
	EmCmd.PersistentFlags().StringVar(&entityType, "type", "", "entity type: Track / Playlist / User")
	EmCmd.PersistentFlags().StringVar(&entityId, "id", "", "entity id")

	// commands for managing developer private key
	keyCmd := &cobra.Command{
		Use:   "key",
		Short: "manage developer app key",
	}
	keyCmd.AddCommand(
		&cobra.Command{
			Use:   "show",
			Short: "show current developer key",
			RunE: func(cmd *cobra.Command, args []string) error {
				fmt.Println("\nreading keyfile", keyfileLocation)

				privateKey, err := crypto.LoadECDSA(keyfileLocation)
				if err != nil {
					return err
				}

				fmt.Println("\n== private ==")
				fmt.Printf("private key \t0x%x\n\n", crypto.FromECDSA(privateKey))

				fmt.Println("\n== public ==")
				fmt.Printf("public key \t0x%x\n", crypto.FromECDSAPub(&privateKey.PublicKey))
				fmt.Printf("public key compressed \t0x%x\n", crypto.CompressPubkey(&privateKey.PublicKey))
				fmt.Printf("address \t%s\n\n", crypto.PubkeyToAddress(privateKey.PublicKey).Hex())

				return nil
			},
		},
		&cobra.Command{
			Use:   "set",
			Short: "set developer key",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) != 1 {
					return errors.New("expect single argument")
				}
				privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(args[0], "0x"))
				if err != nil {
					return err
				}
				fmt.Println("writing keyfile", keyfileLocation)
				return crypto.SaveECDSA(keyfileLocation, privateKey)
			},
		},
		&cobra.Command{
			Use:   "gen",
			Short: "generate + print a random key",
			RunE: func(cmd *cobra.Command, args []string) error {
				pk, err := crypto.GenerateKey()
				if err != nil {
					return err
				}
				asHex := hex.EncodeToString(crypto.FromECDSA(pk))
				fmt.Println(asHex)
				return nil
			},
		},
	)
	EmCmd.AddCommand(keyCmd)

}

func sendEmTx(client *ethclient.Client, actorIdEnc string, action string, entityType string, entityIdEnc string) error {
	privateKey, err := crypto.LoadECDSA(keyfileLocation)
	if err != nil {
		log.Fatal("invalid keyfile: ", err)
	}
	signerAddress := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	actorId, err := hashes.MaybeDecode(actorIdEnc)
	if err != nil {
		log.Fatal("invalid --actor", actorId, err)
	}

	entityId, err := hashes.MaybeDecode(entityIdEnc)
	if err != nil {
		log.Fatal("invalid --id", actorId, err)
	}

	logger := slog.With("Signer", signerAddress, "Actor", actorId, "Action", action, "EntityType", entityType, "EntityID", entityId)

	tx, err := acdc.SendEmTx(client, privateKey, acdc.EmArgs{
		UserID:     int64(actorId),
		Action:     action,
		EntityType: entityType,
		EntityID:   int64(entityId),
	})

	if err != nil {
		logger.Error("failed to send tx", "err", err)
	} else {
		logger.Info("sent tx", "txhash", tx.Hash().Hex())
	}

	return nil
}
