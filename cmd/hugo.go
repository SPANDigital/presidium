package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/hugo"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/proxy"
	"github.com/SPANDigital/presidium-hugo/pkg/log"
	"github.com/spf13/cobra"
)

var (
	// Server command flags
	useProxy   bool
	proxyPort  int
	disableMiddleware bool

	// hugoCommand wraps hugo into Presidium.  This allows you to run hugo
	// in Presidium, and makes it easier to debug etc.
	// All arguments and flags are passed through to Hugo unchanged.
	hugoCommand = &cobra.Command{
		Use:                "hugo",
		Short:              "Runs hugo with full access to all Hugo commands and flags",
		DisableFlagParsing: true, // Don't parse flags - let Hugo handle them
		Run: func(cmd *cobra.Command, args []string) {
			hugoService := hugo.New()
			err := hugoService.Execute(args...)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
		},
	}

	// serverCommand provides a direct way to run 'hugo server'
	serverCommand = &cobra.Command{
		Use:   "server",
		Short: "Start the Hugo development server",
		Long:  "Start the Hugo development server with live reload and other development features.\nBy default, runs with a proxy layer that provides Caddy-style URL rewriting.",
		Run: func(cmd *cobra.Command, args []string) {
			// Check if user wants direct Hugo server or proxy
			if disableMiddleware {
				// Direct Hugo server without proxy
				hugoService := hugo.New()
				err := hugoService.Execute(append([]string{"server"}, args...)...)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
				return
			}

			// Start with proxy middleware (default behavior)
			proxyServer := proxy.NewServer(proxyPort)
			
			// Set up graceful shutdown
			// Handle interrupt signals
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
			
			// Start proxy server in goroutine
			errChan := make(chan error, 1)
			go func() {
				if err := proxyServer.Start(append([]string{"server"}, args...)); err != nil {
					errChan <- err
				}
			}()

			// Wait for interrupt or error
			select {
			case <-sigChan:
				log.Info("Shutting down gracefully...")
				shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer shutdownCancel()
				
				if err := proxyServer.Shutdown(shutdownCtx); err != nil {
					log.Error(err)
					os.Exit(1)
				}
				log.Info("Server stopped")
			case err := <-errChan:
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			}
		},
	}

	// newCommand provides a direct way to run 'hugo new'
	newCommand = &cobra.Command{
		Use:   "new",
		Short: "Create new content for your Hugo site",
		Long:  "Create new content for your Hugo site, such as posts, pages, etc.",
		Run: func(cmd *cobra.Command, args []string) {
			hugoService := hugo.New()
			err := hugoService.Execute(append([]string{"new"}, args...)...)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
		},
	}

	// versionCommand provides a direct way to run 'hugo version'
	hugoVersionCommand = &cobra.Command{
		Use:   "hugo-version",
		Short: "Print the version number of Hugo",
		Long:  "Print the version number of Hugo that Presidium is using",
		Run: func(cmd *cobra.Command, args []string) {
			hugoService := hugo.New()
			err := hugoService.Execute(append([]string{"version"}, args...)...)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	// Add server command flags
	serverCommand.Flags().BoolVar(&disableMiddleware, "no-proxy", false, "Disable proxy middleware and run Hugo server directly")
	serverCommand.Flags().IntVar(&proxyPort, "port", 3131, "Port for the proxy server (Hugo will run on an auto-selected port)")
	
	rootCmd.AddCommand(hugoCommand)
	rootCmd.AddCommand(serverCommand)
	rootCmd.AddCommand(newCommand)
	rootCmd.AddCommand(hugoVersionCommand)
}
