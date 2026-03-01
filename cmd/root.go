package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/av-cli/internal/client"
	"github.com/yourusername/av-cli/internal/config"
)

// Exit codes
const (
	ExitOK      = 0
	ExitGeneral = 1
	ExitInput   = 2
	ExitAuth    = 3
	ExitNetwork = 4
)

// global flags populated before any subcommand runs
var (
	flagPretty  bool
	flagQuiet   bool
	flagTimeout time.Duration
)

// rootCmd is the top-level command.
var rootCmd = &cobra.Command{
	Use:   "av-cli",
	Short: "CLI for DMM and DUGA affiliate APIs",
	Long: `av-cli wraps the DMM API v3 and DUGA API.
Output is always JSON on stdout; status messages go to stderr.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(ExitGeneral)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&flagPretty, "pretty", false, "Pretty-print JSON output")
	rootCmd.PersistentFlags().BoolVar(&flagQuiet, "quiet", false, "Suppress stderr log output")
	rootCmd.PersistentFlags().DurationVar(&flagTimeout, "timeout", 30*time.Second, "HTTP request timeout")
}

// AddCommand registers a subcommand on the root.
func AddCommand(cmds ...*cobra.Command) {
	rootCmd.AddCommand(cmds...)
}

// LoadDMMClient loads config and returns a DMMClient. Exits on missing credentials.
func LoadDMMClient() *client.DMMClient {
	cfg, err := config.Load()
	if err != nil {
		logErr("load config: %v", err)
		os.Exit(ExitGeneral)
	}
	if cfg.DMM.APIID == "" || cfg.DMM.AffiliateID == "" {
		logErr("DMM credentials not set — use DMM_API_ID / DMM_AFFILIATE_ID env vars or ~/.config/av-cli/config.yaml")
		os.Exit(ExitInput)
	}
	return client.NewDMMClient(cfg.DMM.APIID, cfg.DMM.AffiliateID, flagTimeout)
}

// LoadDUGAClient loads config and returns a DUGAClient. Exits on missing credentials.
func LoadDUGAClient() *client.DUGAClient {
	cfg, err := config.Load()
	if err != nil {
		logErr("load config: %v", err)
		os.Exit(ExitGeneral)
	}
	if cfg.DUGA.AppID == "" || cfg.DUGA.AgentID == "" {
		logErr("DUGA credentials not set — use DUGA_APP_ID / DUGA_AGENT_ID env vars or ~/.config/av-cli/config.yaml")
		os.Exit(ExitInput)
	}
	return client.NewDUGAClient(cfg.DUGA.AppID, cfg.DUGA.AgentID, cfg.DUGA.BannerID, flagTimeout)
}

// HandleError maps errors to exit codes and terminates the process.
func HandleError(err error) {
	if err == nil {
		return
	}
	var authErr *client.AuthError
	var apiErr *client.APIError
	switch {
	case errors.As(err, &authErr):
		logErr("%v", err)
		os.Exit(ExitAuth)
	case errors.As(err, &apiErr):
		logErr("%v", err)
		os.Exit(ExitNetwork)
	default:
		logErr("%v", err)
		os.Exit(ExitGeneral)
	}
}

// Logf writes a message to stderr (unless --quiet).
func Logf(format string, args ...any) {
	logErr(format, args...)
}

func logErr(format string, args ...any) {
	if !flagQuiet {
		fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	}
}

// Pretty returns whether --pretty was set.
func Pretty() bool { return flagPretty }

// Quiet returns whether --quiet was set.
func Quiet() bool { return flagQuiet }
