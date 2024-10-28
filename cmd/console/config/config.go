package config

import (
	"os"

	"github.com/spf13/cobra"
)

type Certificate struct {
	Path    string
	KeyPath string
}

type Configuration struct {
	Debug   bool
	Betting struct {
		URL         string
		Certificate Certificate
	}
	Balance           float64
	CallbackServerURL string
	DataBetGQLURL     string
}

// nolint:lll // configuration flags
func LoadConfig(cmd *cobra.Command) Configuration {
	cfg := Configuration{}
	flags := cmd.PersistentFlags()

	flags.BoolVarP(&cfg.Debug, "debug", "d", false, "enable debug mode")
	flags.Float64VarP(&cfg.Balance, "balance", "b", 1000, "Balance of the player")

	flags.StringVarP(&cfg.CallbackServerURL, "callback-url", "u", "http://127.0.0.1:3000/databet", "Callback server URL")

	flags.StringVarP(&cfg.DataBetGQLURL, "databet-gql-url", "g", "https://betting-public-gql-stage-betting.ginsp.net/graphql", "DATA.BET gql server URL")

	flags.StringVarP(&cfg.Betting.URL, "betting-url", "", "https://betting-public-stage-betting.ginsp.net", "Betting server URL")
	flags.StringVarP(&cfg.Betting.Certificate.Path, "betting-certificate-path", "", "./databetstage.crt", "Path to the betting .crt file")
	flags.StringVarP(&cfg.Betting.Certificate.KeyPath, "betting-certificate-key", "", "./databetstage.key", "Path to the betting .key file")

	err := cmd.ParseFlags(os.Args)
	if err != nil {
		panic(err)
	}

	return cfg
}
