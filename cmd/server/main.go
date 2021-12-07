// Command server provides our primary API server.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	cloud "github.com/kmhebb/serverExample"
	utilibill "github.com/kmhebb/serverExample/API/Utilibill"
	"github.com/kmhebb/serverExample/cmd"
	"github.com/kmhebb/serverExample/instrumentation"
	"github.com/kmhebb/serverExample/instrumentation/handlers/cli"
	"github.com/kmhebb/serverExample/internal/service"
	"github.com/kmhebb/serverExample/lib/email"
	"github.com/kmhebb/serverExample/lib/email/slack"
	"github.com/kmhebb/serverExample/lib/random"
	"github.com/kmhebb/serverExample/lib/random/password"
	"github.com/kmhebb/serverExample/lib/token"
	"github.com/kmhebb/serverExample/log"
	"github.com/kmhebb/serverExample/pg"
	"github.com/kmhebb/serverExample/web"
)

const Version = "1.1.1"

func main() {
	cfg := cloud.Config{}
	if err := cfg.Load(os.Args[1:]); err != nil {
		panic(err)
	}
	fmt.Printf("config: %#+v\n", cfg)
	if err := Run(cfg); err != nil {
		panic(err)
	}
}

func Run(cfg cloud.Config) error {
	ctx := context.Background()
	start := time.Now()

	logger := log.NewLogger()

	if cfg.CLI {
		log.SetFormat(log.CLI)
	}
	if cfg.Debug {
		log.SetLevel(log.DebugLevel)
	}

	srv := web.NewServer(cfg.Addr)

	// Next we enable our monitoring provider
	if cfg.CLI {
		instrumentation.SetHandler(cli.New(os.Stderr))
	}

	// We add the guard middleware to prevent Heroku reporting a flood of H18
	// (Request Interrupted) errors resulting from unclosed request bodies.
	srv.Use(web.H18)

	// ts := tags.
	// 	New("env", cfg.Environ).
	// 	Add("version", Version)

	// reporter := ReporterCompat{
	// 	Service: reporting.New(
	// 		reporting.WithRelease(fmt.Sprintf("cloud@v%s", Version)),
	// 		reporting.WithEnvironment(cfg.Environ),
	// 	),
	// }
	var emails email.Service
	switch cfg.Environ {
	case "local":
		token.SetSigningKey(cfg.SigningKey)
		password.Register(random.Passphrase)
		emails = slack.New(cfg.SlackToken, logger)
		utilibill.SetCredentials(cfg.UBusername, cfg.UBpwd)
	case "staging":
		//logger.Log(ctx, log.Info, "Setting signing key")
		token.SetSigningKey(cfg.SigningKey)
		password.Register(random.Passphrase)
		emails = slack.New(cfg.SlackToken, logger)
		utilibill.SetCredentials(cfg.UBusername, cfg.UBpwd)
	case "production":
		token.SetSigningKey(cfg.SigningKey)
		password.Register(random.Passphrase)
		emails = slack.New(cfg.SlackToken, logger)
		utilibill.SetCredentials(cfg.UBusername, cfg.UBpwd)

		// //logger.Log(ctx, log.Info, "Setting signing key")
		// token.SetSigningKey(cfg.SigningKey)
		// password.Register(random.Passphrase)
		// emails = sendgrid.NewService(
		// 	cfg.SendGridKey,
		// 	cfg.SendGridFrom,
		// 	cfg.SendGridEmail,
		// 	cfg.SendGridBaseUrl,
		// 	// reporter,
		//  )
	default:
		// logger.Log(ctx, log.Info, "Setting static signing key")
		token.SetSigningKey("test")
		emails = email.NewNoOpService()
	}

	// Now we connect to our database
	db, err := pg.NewDatabase(ctx, cfg.DatabaseURL)
	if err != nil {
		return cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "failed to connect to database",
			Cause:   err,
		}) //fmt.Errorf("Run: %w", err)
	}

	// We add a circuit breaker that will cause the server to start serving 503s
	// again if we lose connection to the database. The server state will also
	// move to "live" so we pick it up on monitoring.

	srv.AddCircuitBreakers(func() error {
		return db.Ping(ctx)
	})

	// srv.AddCircuitBreakers(func() error {
	// 	return emails.TestConnection()
	// })

	us := service.UserService{
		DB: db,
		L:  logger,
		Em: emails,
	}
	cmd.RegisterUserRoutes(srv, us)

	ds := service.NPDataService{
		DB: db,
		L:  logger,
		Em: emails,
	}
	cmd.RegisterDataServiceRoutes(srv, ds)

	// Finally we're ready to start accepting requests
	log.Info("listening", log.Fields{
		"addr": cfg.Addr, "port": ":8080",
		"startup": time.Since(start),
	})
	srv.Listen()

	// We wait for SIGINTs and SIGTERMs so we can shutdown the server. Graceful
	// shutdown still TODO.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs // Wait for a signal before proceeding.

	// Finally we manually shut down our server before exiting
	return srv.Stop()
}
