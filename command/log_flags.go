package command

import (
	"flag"
	"os"
	"strings"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/posener/complete"
)

// logFlags are the 'log' related flags that can be shared across commands.
type logFlags struct {
	flagLogLevel          string
	flagLogFormat         string
	flagLogFile           string
	flagLogRotateBytes    string
	flagLogRotateDuration string
	flagLogRotateMaxFiles string
}

type provider = func(key string) (string, bool)

type valuesProvider struct {
	flagProvider   provider
	envVarProvider provider
}

// addLogFlags will add the set of 'log' related flags to a flag set.
func (f *FlagSet) addLogFlags(l *logFlags) {
	if l == nil {
		l = &logFlags{}
	}

	f.StringVar(&StringVar{
		Name:       flagNameLogLevel,
		Target:     &l.flagLogLevel,
		Default:    notSetValue,
		EnvVar:     EnvVaultLogLevel,
		Completion: complete.PredictSet("trace", "debug", "info", "warn", "error"),
		Usage: "Log verbosity level. Supported values (in order of detail) are " +
			"\"trace\", \"debug\", \"info\", \"warn\", and \"error\".",
	})

	f.StringVar(&StringVar{
		Name:       flagNameLogFormat,
		Target:     &l.flagLogFormat,
		Default:    notSetValue,
		EnvVar:     EnvVaultLogFormat,
		Completion: complete.PredictSet("standard", "json"),
		Usage:      `Log format. Supported values are "standard" and "json".`,
	})

	f.StringVar(&StringVar{
		Name:   flagNameLogFile,
		Target: &l.flagLogFile,
		EnvVar: EnvVaultLogFile,
		Usage:  "Path to the log file that Vault should use for logging",
	})

	f.StringVar(&StringVar{
		Name:   flagNameLogRotateBytes,
		Target: &l.flagLogRotateBytes,
		EnvVar: EnvVaultLogRotateBytes,
		Usage: "Number of bytes that should be written to a log before it needs to be rotated. " +
			"Unless specified, there is no limit to the number of bytes that can be written to a log file",
	})

	f.StringVar(&StringVar{
		Name:   flagNameLogRotateDuration,
		Target: &l.flagLogRotateDuration,
		EnvVar: EnvVaultLogRotateDuration,
		Usage: "The maximum duration a log should be written to before it needs to be rotated. " +
			"Must be a duration value such as 30s",
	})

	f.StringVar(&StringVar{
		Name:   flagNameLogRotateMaxFiles,
		Target: &l.flagLogRotateMaxFiles,
		EnvVar: EnvVaultLogRotateMaxFiles,
		Usage:  "The maximum number of older log file archives to keep",
	})
}

// getFlagValue will attempt to find the flag with the corresponding key and return the value
// along with a bool representing whether of not the flag had been found/set.
func getFlagValue(fs *FlagSets, key string) (string, bool) {
	var result string
	var isFlagSet bool

	fs.Visit(func(f *flag.Flag) {
		if f.Name == key {
			result = f.Value.String()
			isFlagSet = true
		}
	})

	return result, isFlagSet
}

// getAggregatedConfigValue uses the provided keys to check CLI flags and environment variables for values that may be
// used to override any specified configuration. If nothing can be found the 'fallback' (default) value will be provided.
func (p *valuesProvider) getAggregatedConfigValue(flagKey, envVarKey, current, fallback string) string {
	var result string
	current = strings.TrimSpace(current)

	flg, flgFound := p.flagProvider(flagKey)
	env, envFound := p.envVarProvider(envVarKey)

	switch {
	case flgFound:
		result = flg
	case envFound:
		// Use value from env var
		result = env
	case current != "":
		// Use value from config
		result = current
	default:
		// Use the default value
		result = fallback
	}

	return result
}

// updateLogConfig will accept a shared config and specifically attempt to update the 'log' related config keys.
// For each 'log' key we aggregate file config/env vars and CLI flags to select the one with the highest precedence.
// This method mutates the config object passed into it.
func (f *FlagSets) updateLogConfig(config *configutil.SharedConfig) {
	p := &valuesProvider{
		flagProvider:   func(key string) (string, bool) { return getFlagValue(f, key) },
		envVarProvider: os.LookupEnv,
	}

	config.LogLevel = p.getAggregatedConfigValue(flagNameLogLevel, EnvVaultLogLevel, config.LogLevel, "info")
	config.LogFormat = p.getAggregatedConfigValue(flagNameLogFormat, EnvVaultLogFormat, config.LogFormat, "")
	config.LogFile = p.getAggregatedConfigValue(flagNameLogFile, EnvVaultLogFile, config.LogFile, "")
	config.LogRotateDuration = p.getAggregatedConfigValue(flagNameLogRotateDuration, EnvVaultLogRotateDuration, config.LogRotateDuration, "")
	config.LogRotateBytes = p.getAggregatedConfigValue(flagNameLogRotateBytes, EnvVaultLogRotateBytes, config.LogRotateBytes, "")
	config.LogRotateMaxFiles = p.getAggregatedConfigValue(flagNameLogRotateMaxFiles, EnvVaultLogRotateMaxFiles, config.LogRotateMaxFiles, "")
}
