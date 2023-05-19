// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package omnictl

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"text/tabwriter"
	"time"

	"github.com/siderolabs/go-api-signature/pkg/pgp"
	"github.com/spf13/cobra"

	pkgaccess "github.com/siderolabs/omni-client/pkg/access"
	"github.com/siderolabs/omni-client/pkg/client"
	"github.com/siderolabs/omni-client/pkg/omnictl/internal/access"
)

var (
	serviceAccountCreateFlags struct {
		scopes []string

		useUserScopes bool
		ttl           time.Duration
	}

	serviceAccountRenewFlags struct {
		ttl time.Duration
	}

	// serviceAccountCmd represents the serviceaccount command.
	serviceAccountCmd = &cobra.Command{
		Use:     "serviceaccount",
		Aliases: []string{"sa"},
		Short:   "Manage service accounts",
	}

	serviceAccountCreateCmd = &cobra.Command{
		Use:     "create <name>",
		Aliases: []string{"c"},
		Short:   "Create a service account",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			return access.WithClient(func(ctx context.Context, client *client.Client) error {
				comment := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

				key, err := pgp.GenerateKey(name, comment, name, serviceAccountCreateFlags.ttl)
				if err != nil {
					return err
				}

				armoredPublicKey, err := key.ArmorPublic()
				if err != nil {
					return err
				}

				publicKeyID, err := client.Management().CreateServiceAccount(ctx, armoredPublicKey, serviceAccountCreateFlags.scopes, serviceAccountCreateFlags.useUserScopes)
				if err != nil {
					return err
				}

				encodedKey, err := encodeServiceAccountKey(name, key)
				if err != nil {
					return err
				}

				fmt.Printf("Created service account %q with public key ID %q\n", name, publicKeyID)
				fmt.Printf("\n")
				fmt.Printf("Set the following environment variables to use the service account:\n")
				fmt.Printf("OMNI_ENDPOINT=%s\n", client.Endpoint())
				fmt.Printf("OMNI_SERVICE_ACCOUNT_KEY=%s\n", encodedKey)
				fmt.Printf("\n")
				fmt.Printf("Note: Store the service account key securely, it will not be displayed again\n")

				return nil
			})
		},
	}

	serviceAccountRenewCmd = &cobra.Command{
		Use:     "renew <name>",
		Aliases: []string{"r"},
		Short:   "Renew a service account by registering a new public key to it",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			return access.WithClient(func(ctx context.Context, client *client.Client) error {
				comment := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

				key, err := pgp.GenerateKey(name, comment, name, serviceAccountRenewFlags.ttl)
				if err != nil {
					return err
				}

				armoredPublicKey, err := key.ArmorPublic()
				if err != nil {
					return err
				}

				publicKeyID, err := client.Management().RenewServiceAccount(ctx, name, armoredPublicKey)
				if err != nil {
					return err
				}

				encodedKey, err := encodeServiceAccountKey(name, key)
				if err != nil {
					return err
				}

				fmt.Printf("Renewed service account %q by adding a public key with ID %q\n", name, publicKeyID)
				fmt.Printf("\n")
				fmt.Printf("Set the following environment variables to use the service account:\n")
				fmt.Printf("OMNI_ENDPOINT=%s\n", client.Endpoint())
				fmt.Printf("OMNI_SERVICE_ACCOUNT_KEY=%s\n", encodedKey)
				fmt.Printf("\n")
				fmt.Printf("Note: Store the service account key securely, it will not be displayed again\n")

				return nil
			})
		},
	}

	serviceAccountListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List service accounts",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return access.WithClient(func(ctx context.Context, client *client.Client) error {
				serviceAccounts, err := client.Management().ListServiceAccounts(ctx)
				if err != nil {
					return err
				}

				writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

				fmt.Fprintf(writer, "NAME\tSCOPES\tPUBLIC KEY ID\tEXPIRATION\n")

				for _, sa := range serviceAccounts {
					for i, publicKey := range sa.PgpPublicKeys {
						if i == 0 {
							fmt.Fprintf(writer, "%s\t%q\t%s\t%s\n", sa.Name, sa.Scopes, publicKey.Id, publicKey.Expiration.AsTime().String())
						} else {
							fmt.Fprintf(writer, "\t\t%s\t%s\n", publicKey.Id, publicKey.Expiration.AsTime().String())
						}
					}
				}

				return writer.Flush()
			})
		},
	}

	serviceAccountDestroyCmd = &cobra.Command{
		Use:     "destroy <name>",
		Aliases: []string{"d"},
		Short:   "Destroy a service account",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			return access.WithClient(func(ctx context.Context, client *client.Client) error {
				err := client.Management().DestroyServiceAccount(ctx, name)
				if err != nil {
					return fmt.Errorf("failed to destroy service account: %w", err)
				}

				fmt.Printf("destroyed service account: %s\n", name)

				return nil
			})
		},
	}
)

// encodeServiceAccountKey encodes a service account key to a base64-encoded JSON.
func encodeServiceAccountKey(name string, key *pgp.Key) (string, error) {
	armoredPrivateKey, err := key.Armor()
	if err != nil {
		return "", fmt.Errorf("failed to armor private key: %w", err)
	}

	saKey := pkgaccess.ServiceAccountKey{
		Name:   name,
		PGPKey: armoredPrivateKey,
	}

	saKeyJSON, err := json.Marshal(saKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(saKeyJSON), nil
}

func init() {
	RootCmd.AddCommand(serviceAccountCmd)

	serviceAccountCmd.AddCommand(serviceAccountCreateCmd)
	serviceAccountCmd.AddCommand(serviceAccountListCmd)
	serviceAccountCmd.AddCommand(serviceAccountDestroyCmd)
	serviceAccountCmd.AddCommand(serviceAccountRenewCmd)

	serviceAccountCreateCmd.Flags().DurationVarP(&serviceAccountCreateFlags.ttl, "ttl", "t", 365*24*time.Hour, "TTL for the service account key")
	serviceAccountCreateCmd.Flags().StringSliceVarP(&serviceAccountCreateFlags.scopes, "scopes", "s", nil, "scopes of the service account. only used when --use-user-scopes=false")
	serviceAccountCreateCmd.Flags().BoolVarP(&serviceAccountCreateFlags.useUserScopes, "use-user-scopes", "u", true, "use the scopes of the creating user. if true, --scopes is ignored")

	serviceAccountRenewCmd.Flags().DurationVarP(&serviceAccountRenewFlags.ttl, "ttl", "t", 365*24*time.Hour, "TTL for the service account key")
}
