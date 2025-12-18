package cli

import (
	"context"
	"fmt"
	"io"
	"maps"
	"slices"

	"github.com/alecthomas/kong"
	"github.com/sj14/sss/controller"
	"github.com/sj14/sss/util"
)

func Exec(ctx context.Context, outWriter, errWriter io.Writer, buildInfo util.BuildInfo) error {
	cli := CLI{
		Version: Version{
			info: buildInfo,
		},
	}

	kctx := kong.Parse(
		&cli,
		kong.Name("sss"),
		kong.DefaultEnvars("SSS"),
	)

	config, err := controller.LoadConfig(cli.Config)
	if err != nil {
		return err
	}

	profile, ok := config.Profiles[cli.Profile]
	if !ok && (cli.Config != "" || cli.Profile != "default") {
		fmt.Fprintf(errWriter, "available profiles:\n")

		keys := slices.Collect(maps.Keys(config.Profiles))
		slices.Sort(keys)

		for _, key := range keys {
			fmt.Fprintf(errWriter, "  %s\n", key)
		}

		return fmt.Errorf("profile not found: %q", cli.Profile)
	}

	util.SetIfNotZero(&profile.Endpoint, cli.Endpoint)
	util.SetIfNotZero(&profile.Region, cli.Region)
	util.SetIfNotZero(&profile.PathStyle, cli.PathStyle)
	util.SetIfNotZero(&profile.AccessKey, cli.AccessKey)
	util.SetIfNotZero(&profile.SecretKey, cli.SecretKey)
	util.SetIfNotZero(&profile.Insecure, cli.Insecure)
	util.SetIfNotZero(&profile.ReadOnly, cli.ReadOnly)
	util.SetIfNotZero(&profile.SNI, cli.SNI)
	util.SetIfNotZero(&profile.Network, cli.Network)
	util.SetIfNotZero(&profile.Bandwidth, cli.Bandwidth)

	dryRun := isFlagSet(kctx.Selected().Flags, "dry-run")

	ctrl, err := controller.New(
		ctx,
		controller.ControllerConfig{
			OutWriter: outWriter,
			ErrWriter: errWriter,
			Profile:   profile,
			Verbosity: cli.Verbosity,
			Headers:   cli.Headers,
			Params:    cli.Params,
			DryRun:    dryRun,
			BuildInfo: buildInfo,
		})
	if err != nil {
		return err
	}

	err = kctx.Run(cli, ctrl, config)
	if err != nil {
		return err
	}

	return nil
}
