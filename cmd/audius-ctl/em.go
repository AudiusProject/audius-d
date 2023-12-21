package main

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/AudiusProject/audius-d/pkg/acdc"
	"github.com/AudiusProject/audius-d/pkg/hashes"
	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

var emCmd *cobra.Command

const (
	keyfileLocation = ".audius_developer_app_private_key"
)

func init() {

	var actorId string
	var action string
	var entityType string
	var entityId string

	emCmd = &cobra.Command{
		Use:   "em",
		Short: "Send EntityManager transactions",
		RunE: func(cmd *cobra.Command, args []string) error {

			// todo: get endpoint from env (stage / prod / other)
			const ACDC_ENDPOINT = `https://acdc-gateway.staging.audius.co`
			client, err := ethclient.Dial(ACDC_ENDPOINT)
			if err != nil {
				return logger.Error(err)
			}

			privateKey, err := crypto.LoadECDSA(keyfileLocation)
			if err != nil {
				return logger.Error("invalid keyfile: ", err)
			}

			actorIdInt, err := hashes.MaybeDecode(actorId)
			if err != nil {
				return logger.Error("invalid --actor", actorIdInt, err)
			}

			entityIdInt, err := hashes.MaybeDecode(entityId)
			if err != nil {
				return logger.Error("invalid --id", actorIdInt, err)
			}

			tx, err := acdc.SendEmTx(client, privateKey, acdc.EmArgs{
				UserID:     int64(actorIdInt),
				Action:     action,
				EntityType: entityType,
				EntityID:   int64(entityIdInt),
			})

			if err != nil {
				return logger.Error("failed to send tx", "err", err)
			} else {
				logger.Info("sent tx", "txhash", tx.Hash().Hex())
			}
			return nil
		},
	}

	emCmd.Flags().StringVar(&actorId, "actor", "", "user id performing the action")
	emCmd.Flags().StringVar(&action, "action", "", "verb to perform: Repost / Save / Follow")
	emCmd.Flags().StringVar(&entityType, "type", "", "entity type: Track / Playlist / User")
	emCmd.Flags().StringVar(&entityId, "id", "", "entity id")

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
				logger.Infof("Reading keyfile %s", keyfileLocation)

				privateKey, err := crypto.LoadECDSA(keyfileLocation)
				if err != nil {
					return err
				}

				logger.Out("\n== private ==")
				logger.Out("private key \t0x%x\n\n", crypto.FromECDSA(privateKey))

				logger.Out("\n== public ==")
				logger.Out("public key \t0x%x\n", crypto.FromECDSAPub(&privateKey.PublicKey))
				logger.Out("public key compressed \t0x%x\n", crypto.CompressPubkey(&privateKey.PublicKey))
				logger.Out("address \t%s\n\n", crypto.PubkeyToAddress(privateKey.PublicKey).Hex())

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
				logger.Infof("Writing keyfile %s", keyfileLocation)
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
				logger.Out(asHex)
				return nil
			},
		},
	)
	emCmd.AddCommand(keyCmd)

}
