package cmd

import (
	"context"
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	entsql "github.com/facebookincubator/ent/dialect/sql"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/meringu/terraform-private-registry/internal/ent"
	"github.com/meringu/terraform-private-registry/internal/handlers"
	"github.com/meringu/terraform-private-registry/internal/handlers/middleware"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "starts the Terraform private module registry",
	Long:  "starts the Terraform private module registry",
	Run: func(cmd *cobra.Command, args []string) {
		level, err := logrus.ParseLevel(logLevel)
		if err != nil {
			logrus.Fatalf("%v", err)
		}
		logrus.SetLevel(level)

		validDriver := false
		for _, driver := range sql.Drivers() {
			if driver == dbDriverName {
				validDriver = true
			}
		}
		if !validDriver {
			logrus.Fatalf("`%s` not a supported database driver. Supported drivers: `%s`", dbDriverName, strings.Join(sql.Drivers(), "`, `"))
		}

		parsedGitHubURL, err := url.Parse(gitHubURL)
		if err != nil {
			logrus.Fatalf("Failed to parse github-url: %v", err)
		}

		gitHubWebhookSecretBytes := []byte(gitHubWebhookSecret)
		if gitHubWebhookSecretPath != "" {
			gitHubWebhookSecretBytes, err = ioutil.ReadFile(gitHubPrivateKeyPath)
			if err != nil {
				logrus.Fatalf("Failed to read github-webhook-secret: %v", err)
			}
		}
		if len(gitHubWebhookSecretBytes) == 0 {
			logrus.Fatal("github-webhook-secret must be specified")
		}

		gitHubPrivateKey, err := ioutil.ReadFile(gitHubPrivateKeyPath)
		if err != nil {
			logrus.Fatalf("Failed to read GitHub private key: %s", err)
		}

		drv, err := entsql.Open(dbDriverName, dbDSN)
		if err != nil {
			logrus.Fatalf("failed opening connection to DB: %v", err)
		}
		// Get the underlying sql.DB object of the driver.
		db := drv.DB()
		db.SetMaxIdleConns(10)
		db.SetMaxOpenConns(100)
		db.SetConnMaxLifetime(time.Hour)
		defer db.Close()

		entClient := ent.NewClient(ent.Driver(drv))

		logrus.Info("Migrating DB schema")

		if err := entClient.Schema.Create(context.Background()); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}

		mux := mux.NewRouter()
		mux.Use(middleware.LoggerMiddleware)

		mux.Handle("/.well-known/terraform.json", handlers.DiscoveryHandler()).Methods("GET", "HEAD")
		mux.Handle("/ping", handlers.PingHandler()).Methods("GET", "HEAD")

		mux.Handle("/v1/search", handlers.NotImplementedHandler()).Methods("GET", "HEAD")
		mux.Handle("/v1/modules", handlers.ModulesListHandler(entClient)).Methods("GET", "HEAD")
		mux.Handle("/v1/modules/{namespace}", handlers.ModulesListHandler(entClient)).Methods("GET", "HEAD")
		mux.Handle("/v1/modules/{namespace}/{name}", handlers.ModulesListHandler(entClient)).Methods("GET", "HEAD")
		mux.Handle("/v1/modules/{namespace}/{name}/{provider}", handlers.ModulesGetVersionHandler(entClient)).Methods("GET", "HEAD")
		mux.Handle("/v1/modules/{namespace}/{name}/{provider}/download", handlers.ModulesDownloadVersionHandler(entClient, parsedGitHubURL, gitHubPrivateKey)).Methods("GET", "HEAD")
		mux.Handle("/v1/modules/{namespace}/{name}/{provider}/versions", handlers.ModulesListVersionsHandler(entClient)).Methods("GET", "HEAD")
		mux.Handle("/v1/modules/{namespace}/{name}/{provider}/{version}", handlers.ModulesGetVersionHandler(entClient)).Methods("GET", "HEAD")
		mux.Handle("/v1/modules/{namespace}/{name}/{provider}/{version}/download", handlers.ModulesDownloadVersionHandler(entClient, parsedGitHubURL, gitHubPrivateKey)).Methods("GET", "HEAD")

		mux.Handle("/githubapp/events", handlers.GitHubEventHandler(entClient, gitHubWebhookSecretBytes, parsedGitHubURL, gitHubPrivateKey)).Methods("POST")

		mux.PathPrefix("/").Handler(handlers.NotFoundHandler())

		server := &http.Server{
			Addr:         bind,
			Handler:      mux,
			WriteTimeout: 1800 * time.Second,
			ReadTimeout:  60 * time.Second,
		}
		logrus.Infof("Starting server on port %s", bind)

		logrus.Error(server.ListenAndServe())
	},
}

var bind string
var logLevel string
var dbDriverName string
var dbDSN string
var gitHubWebhookSecret string
var gitHubWebhookSecretPath string
var gitHubURL string
var gitHubPrivateKeyPath string

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&bind, "bind", "", ":80", "address to bind to")
	serverCmd.Flags().StringVarP(&logLevel, "log-level", "l", logrus.InfoLevel.String(), "log level")
	serverCmd.Flags().StringVarP(&dbDriverName, "db-driver", "", "sqlite3", "database driver name")
	serverCmd.Flags().StringVarP(&dbDSN, "db-dsn", "", "file:ent?mode=memory&cache=shared&_fk=1", "database data source name")
	serverCmd.Flags().StringVarP(&gitHubWebhookSecret, "github-webhook-secret", "", "", "secret for GitHub event validation")
	serverCmd.Flags().StringVarP(&gitHubWebhookSecretPath, "github-webhook-secret-path", "", "", "path to secret for GitHub event validation")
	serverCmd.Flags().StringVarP(&gitHubURL, "github-url", "", "https://api.github.com/", "URL for GitHub")
	serverCmd.Flags().StringVarP(&gitHubPrivateKeyPath, "github-private-key-path", "", "", "path to the GitHub App private key")

	serverCmd.MarkFlagRequired("github-private-key-path")
}
