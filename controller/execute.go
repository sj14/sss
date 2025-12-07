package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"slices"

	"github.com/sj14/sss/util"
)

type ExecuteInput struct {
	StdoutWriter io.Writer
	StdErrWriter io.Writer
	ConfigPath   string
	ProfileName  string
	ControllerConfig
}

func Execute(ctx context.Context, cfg ExecuteInput, fn func(ctrl *Controller) error) error {
	config, err := LoadConfig(cfg.ConfigPath)
	// Do not return an error when the config file does not exist,
	// as the tool should be usable wihout a config file.
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("failed loading config: %w", err)
	}

	profile, ok := config.Profiles[cfg.ProfileName]
	if !ok && cfg.ProfileName != "default" {
		fmt.Printf("profile %q not found, available profiles:\n", cfg.ProfileName)

		keys := slices.Collect(maps.Keys(config.Profiles))
		slices.Sort(keys)

		for _, key := range keys {
			fmt.Println(key)
		}

		os.Exit(1)
	}

	util.SetIfNotZero(&profile.Endpoint, cfg.Profile.Endpoint)
	util.SetIfNotZero(&profile.Region, cfg.Profile.Region)
	util.SetIfNotZero(&profile.PathStyle, cfg.Profile.PathStyle)
	util.SetIfNotZero(&profile.AccessKey, cfg.Profile.AccessKey)
	util.SetIfNotZero(&profile.SecretKey, cfg.Profile.SecretKey)
	util.SetIfNotZero(&profile.Insecure, cfg.Profile.Insecure)
	util.SetIfNotZero(&profile.ReadOnly, cfg.Profile.ReadOnly)
	util.SetIfNotZero(&profile.SNI, cfg.Profile.SNI)

	var bandwidth uint64 = 0
	// {
	// 	bwStr := cfg.Bandwidth
	// 	if bwStr != "" {
	// 		bandwidth, err = humanize.ParseBytes(bwStr)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	ctrl, err := New(
		ctx,
		cfg.StdoutWriter,
		cfg.StdErrWriter,
		ControllerConfig{
			Profile:   profile,
			Bandwidth: bandwidth,
			Verbosity: cfg.Verbosity,
			Headers:   cfg.Headers,

			// non-global flag
			DryRun: cfg.DryRun,
		})
	if err != nil {
		return err
	}
	return fn(ctrl)
}
