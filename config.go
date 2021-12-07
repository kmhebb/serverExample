package cloud

import (
	"github.com/kmhebb/serverExample/lib/flag"
	"github.com/kmhebb/serverExample/lib/os"
)

type Config struct {
	Environ         string
	Addr            string
	CLI             bool
	DatabaseURL     string
	DatabaseVersion int
	Debug           bool
	Shout           bool
	SigningKey      string
	SlackToken      string
	SendGridKey     string
	SendGridFrom    string
	SendGridEmail   string
	SendGridBaseUrl string
	UBusername      string
	UBpwd           string
}

func (cfg *Config) Load(args []string) error {
	var printHelp bool

	fs := flag.NewFlagSet("server")
	fs.StringVar(
		&cfg.Environ,
		"e",
		"cloud_ENV",
		os.GetStringEnv("cloud_ENV"),
		"The environment",
	)
	fs.StringVar(
		&cfg.Addr,
		"addr",
		"cloud_SERVER_ADDRESS",
		os.GetStringEnv("cloud_SERVER_ADDRESS"),
		"The server address",
	)
	fs.StringVar(
		&cfg.SigningKey,
		"sk",
		"cloud_SIGNING_KEY",
		os.GetStringEnv("cloud_SIGNING_KEY"),
		"The signing key",
	)
	fs.StringVar(
		&cfg.SlackToken,
		"st",
		"cloud_SLACK_BOT_TOKEN",
		os.GetStringEnv("cloud_SLACK_BOT_TOKEN"),
		"The token for slack-bot",
	)
	fs.StringVar(
		&cfg.UBusername,
		"ubn",
		"cloud_UB_USERNAME",
		os.GetStringEnv("cloud_UB_USERNAME"),
		"The username for the utilibill api",
	)
	fs.StringVar(
		&cfg.UBpwd,
		"ubp",
		"cloud_UB_PW",
		os.GetStringEnv("cloud_UB_PW"),
		"The password for the utilibill api",
	)
	// fs.StringVar(
	// 	&cfg.SendGridKey,
	// 	"sgk",
	// 	"SendGridKey",
	// 	os.GetStringEnv("SENDGRID_KEY"),
	// 	"The key for Sendgrid",
	// )
	// fs.StringVar(
	// 	&cfg.SendGridFrom,
	// 	"sgf",
	// 	"SendGridFromName",
	// 	os.GetStringEnv("SENDGRID_FROM_NAME"),
	// 	"The from name for Sendgrid",
	// )
	// fs.StringVar(
	// 	&cfg.SendGridEmail,
	// 	"sge",
	// 	"SendGridFromEmail",
	// 	os.GetStringEnv("SENDGRID_FROM_EMAIL"),
	// 	"The from email for Sendgrid",
	// )
	// fs.StringVar(
	// 	&cfg.SendGridBaseUrl,
	// 	"sgu",
	// 	"SendGridBaseURL",
	// 	os.GetStringEnv("SENDGRID_BASE_URL"),
	// 	"The base url for Sendgrid",
	// )
	fs.BoolVar(
		&printHelp,
		"h",
		"help",
		false,
		"Print help information",
	)
	fs.StringVar(
		&cfg.DatabaseURL,
		"db",
		"database-url",
		os.GetStringEnv("cloud_DATABASE_URL"),
		"The URL of a database to connect to",
	)
	fs.BoolVar(
		&cfg.Debug,
		"d",
		"debug",
		true,
		"Enable debug logging",
	)
	fs.BoolVar(
		&cfg.Shout,
		"",
		"shout",
		true,
		"Enable shouty messages",
	)
	// fs.IntVar(
	// 	&cfg.DatabaseVersion,
	// 	"",
	// 	"database-version",
	// 	os.GetIntEnv("DATABASE_VERSION"),
	// 	"The database version to use",
	// )
	fs.Parse(args)

	// if printHelp {
	// 	fmt.Println(strings.TrimSpace(usage))
	// 	os.Exit(2) // TODO: This isn't my favorite way to do this
	// }

	return nil
}

// const usage = `
// Usage: server --database-url DATABASE [OPTION]...
// Connect to DATABASE and run the server.

//   --help                Print this help and exit
//   --debug               Enable verbose logging
//   --shout               Use the shouting greeter service
//   --database-version    Set the target database version
// `
