package confighelper

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"os"
	"syscall"
)

func LoadConfig(ctx context.Context, path string, v interface{}) error {
	var err error
	logger := zerolog.Ctx(ctx)
	logger.Debug().Str("file", path).Msg("start loading config")
	if err = CheckPathExist(ctx, path, 4, ""); err != nil {
		logger.Error().Err(err).Msg("config file path is invalid")
		return err
	}

	viper.SetConfigFile(path)
	err = viper.ReadInConfig()
	if err != nil {
		logger.Error().Err(err).Send()
		return err
	}
	err = viper.Unmarshal(v)
	if err != nil {
		logger.Error().Err(err).Send()
		return err
	}

	logger.Debug().Msg("load config success")
	return nil
}

// minPermission:
// 4 -> only check if can read
// 4 + 2 = 6 -> check if can read and write
func CheckPathExist(ctx context.Context, path string, permission int, desc string) error {
	// first check if exist
	logger := zerolog.Ctx(ctx)
	if _, err := os.Stat(path); err != nil {
		logger.Warn().Err(err).Msg(desc)
		if os.IsNotExist(err) {
			// todo
		} else {
		}
		return err
	}

	// then check if can read or read/write
	var bit uint32 = syscall.O_RDWR
	if permission < 6 {
		bit = syscall.O_RDONLY
	}

	err := syscall.Access(path, bit)
	if err != nil {
		logger.Warn().Err(err).Msg(desc)
		return err
	}

	return nil
}
