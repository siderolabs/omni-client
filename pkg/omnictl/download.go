// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package omnictl

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/siderolabs/gen/slices"
	"github.com/siderolabs/go-api-signature/pkg/message"
	keyClient "github.com/siderolabs/go-api-signature/pkg/pgp/client"
	"github.com/spf13/cobra"

	"github.com/siderolabs/omni-client/pkg/client"
	"github.com/siderolabs/omni-client/pkg/omni/resources/omni"
	"github.com/siderolabs/omni-client/pkg/omnictl/config"
	"github.com/siderolabs/omni-client/pkg/omnictl/internal/access"
)

// downloadFlags represents the `download` command flags.
type downloadFlags struct {
	architecture string

	output string
}

var downloadCmdFlags downloadFlags

func init() {
	downloadCmd.Flags().StringVar(&downloadCmdFlags.architecture, "arch", "amd64", "Image architecture to download (amd64, arm64)")
	downloadCmd.Flags().StringVar(&downloadCmdFlags.output, "output", ".", "Output file or directory, defaults to current working directory")

	RootCmd.AddCommand(downloadCmd)
}

// downloadCmd represents the download command.
var downloadCmd = &cobra.Command{
	Use:   "download <image name>",
	Short: "Download installer media",
	Long: `This command downloads installer media from the server

It accepts one argument, which is the name of the image to download. Name can be one of the following:
     
     * iso - downloads the latest ISO image
     * AWS AMI (amd64), Vultr (arm64), Raspberry Pi 4 Model B - full image name
     * oracle, aws, vmware - platform name
     * rockpi_4, rock64 - board name

To get the full list of available images, look at the output of the following command:
    omnictl get installationmedia -o yaml

The download command tries to match the passed string in this order:

    * name
    * platform
    * board

By default it will download amd64 image if there are multiple images available for the same name.

For example, to download the latest ISO image for arm64, run:

    omnictl download iso --arch amd64

To download the latest Vultr image, run:

    omnictl download "vultr"

To download the latest Radxa ROCK PI 4 image, run:

    omnictl download "rockpi_4"
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return access.WithClient(func(ctx context.Context, client *client.Client) error {
			if args[0] == "" {
				return fmt.Errorf("image name is required")
			}

			output, err := filepath.Abs(downloadCmdFlags.output)
			if err != nil {
				return err
			}

			err = makePath(output)
			if err != nil {
				return err
			}

			image, err := findImage(ctx, client, args[0], downloadCmdFlags.architecture)
			if err != nil {
				return err
			}

			return downloadImageTo(ctx, client, image, output)
		})
	},
	ValidArgsFunction: downloadCompletion,
}

func findImage(ctx context.Context, client *client.Client, name, arch string) (*omni.InstallationMedia, error) {
	result, err := filterMedia(ctx, client, func(val *omni.InstallationMedia) (*omni.InstallationMedia, bool) {
		spec := val.TypedSpec().Value

		if strings.EqualFold(name, "iso") {
			return val, spec.Type == "iso"
		}

		return val, strings.EqualFold(spec.Name, name) ||
			strings.EqualFold(spec.Platform, name) ||
			strings.EqualFold(spec.Board, name)
	})
	if err != nil {
		return nil, err
	}

	if len(result) > 1 {
		result = slices.FilterInPlace(result, func(val *omni.InstallationMedia) bool {
			return val.TypedSpec().Value.Architecture == arch
		})
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no image found for %q", name)
	} else if len(result) > 1 {
		names := slices.Map(result, func(val *omni.InstallationMedia) string {
			return val.TypedSpec().Value.Filename
		})

		return nil, fmt.Errorf("multiple images found:\n  %s", strings.Join(names, "\n  "))
	}

	return result[0], nil
}

func downloadImageTo(ctx context.Context, client *client.Client, media *omni.InstallationMedia, output string) error {
	req, err := createRequest(ctx, client, media)
	if err != nil {
		return err
	}

	err = signRequest(req)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer checkCloser(resp.Body)

	filename := media.TypedSpec().Value.Filename

	dest := output
	if filepath.Ext(output) == "" {
		dest = filepath.Join(output, filename)
	}

	fmt.Printf("Downloading %s to %s\n", filename, dest)

	err = downloadResponseTo(dest, resp)
	if err != nil {
		return err
	}

	fmt.Printf("Downloaded %s to %s\n", filename, dest)

	return nil
}

func filterMedia[T any](ctx context.Context, client *client.Client, check func(value *omni.InstallationMedia) (T, bool)) ([]T, error) {
	media, err := safe.StateListAll[*omni.InstallationMedia](
		ctx,
		client.Omni().State(),
	)
	if err != nil {
		return nil, err
	}

	var result []T

	for it := safe.IteratorFromList(media); it.Next(); {
		if val, ok := check(it.Value()); ok {
			result = append(result, val)
		}
	}

	return result, nil
}

func createRequest(ctx context.Context, client *client.Client, image *omni.InstallationMedia) (*http.Request, error) {
	u, err := url.Parse(client.Endpoint())
	if err != nil {
		return nil, err
	}

	u.Scheme = "https"

	u.Path, err = url.JoinPath(u.Path, "image", image.Metadata().ID())
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, err
}

func signRequest(req *http.Request) error {
	contextName, configCtx, err := currentConfigCtx()
	if err != nil {
		return err
	}

	email := configCtx.Auth.SideroV1.Identity
	provider := keyClient.NewKeyProvider("omni/keys")

	key, err := provider.ReadValidKey(contextName, email)
	if err != nil {
		return err
	}

	msg, err := message.NewHTTP(req)
	if err != nil {
		return err
	}

	return msg.Sign(email, key)
}

func currentConfigCtx() (name string, ctx *config.Context, err error) {
	conf, err := config.Current()
	if err != nil {
		return "", nil, err
	}

	contextName := conf.Context
	if access.CmdFlags.Context != "" {
		contextName = access.CmdFlags.Context
	}

	configCtx, err := conf.GetContext(contextName)
	if err != nil {
		return "", nil, err
	}

	return contextName, configCtx, nil
}

func downloadResponseTo(dest string, resp *http.Response) error {
	f, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer checkCloser(f)

	_, err = io.Copy(f, resp.Body)

	return err
}

func checkCloser(c io.Closer) {
	if err := c.Close(); err != nil {
		fmt.Printf("error closing: %v", err)
	}
}

func downloadCompletion(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var results []string

	err := access.WithClient(
		func(ctx context.Context, client *client.Client) error {
			res, err := filterMedia(ctx, client, func(value *omni.InstallationMedia) (string, bool) {
				spec := value.TypedSpec().Value
				if downloadCmdFlags.architecture != spec.Architecture {
					return "", false
				}

				name := spec.Name
				if toComplete == "" || strings.Contains(name, toComplete) {
					return name, true
				}

				return "", false
			})
			if err != nil {
				return err
			}

			results = res

			return nil
		},
	)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	return dedupInplace(results), cobra.ShellCompDirectiveNoFileComp
}

func dedupInplace(results []string) []string {
	seen := make(map[string]struct{}, len(results))
	j := 0

	for _, r := range results {
		if _, ok := seen[r]; !ok {
			seen[r] = struct{}{}
			results[j] = r
			j++
		}
	}

	return results[:j]
}

func makePath(path string) error {
	if filepath.Ext(path) != "" {
		ok, err := checkPath(path)
		if err != nil {
			return err
		}

		if ok {
			return fmt.Errorf("destination %s already exists", path)
		}

		path = filepath.Dir(path)
	}

	ok, err := checkPath(path)
	if err != nil {
		return err
	}

	if !ok {
		if dirErr := os.MkdirAll(path, 0o755); dirErr != nil {
			return err
		}
	}

	return nil
}

func checkPath(path string) (bool, error) {
	_, err := os.Stat(path)

	switch {
	case os.IsNotExist(err):
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}
