package auth

import (
	"context"
	"fmt"
	"html/template"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/a0dotrun/a0ctl/internal/cli"
	"github.com/a0dotrun/a0ctl/internal/flags"
	"github.com/a0dotrun/a0ctl/internal/settings"

	"github.com/a0dotrun/a0ctl/internal/api"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"

	_ "embed"
)

// FIXME: @sanchitrk
//
//go:embed
var loginHTML string

func newLogin() *cobra.Command {
	const (
		use   = "login"
		short = "new user logged in"
	)
	cmd := &cobra.Command{
		Use:               "login",
		Short:             "Login to the platform.",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cli.NoFilesArg,
		RunE:              login,
		PersistentPreRunE: checkEnvAuth,
	}
	return cmd
}

func exitOnValidAuth(settings *settings.Settings) {
	username := settings.GetUsername()
	if len(username) <= 0 {
		fmt.Println("✔  Success! Existing JWT still valid")
		return
	}
	fmt.Printf("Already signed in as %s. Use %s to log out of this account\n", username, cli.Emph("a0ctl auth logout"))
}

func login(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	config, err := settings.ReadSettings()
	if err != nil {
		return fmt.Errorf("could not retrieve local config: %w", err)
	}

	if api.IsJWTTokenValid(config.GetToken()) {
		exitOnValidAuth(config)
		return nil
	}

	// TODO: @sanchitrk
	// add support for headless login

	if flags.Headless() {
		return printHeadlessLoginInstructions(authURLPath)
	}

	state := randString(32)
	callbackServer, err := authCallbackServer(state)
	if err != nil {
		return suggestHeadless(cmd, err)
	}

	// Now, that we got the callback server, let's get the auth URL
	// for making the auth request

	// NOTE:: path must match the one in the auth server
	url, err := authURL(callbackServer.Port, authURLPath, state)
	if err != nil {
		return fmt.Errorf("failed to get auth URL: %w", err)
	}

	if err := browser.OpenURL(url); err != nil {
		err := fmt.Errorf("failed to open auth URL: %w", err)
		return suggestHeadless(cmd, err)
	}

	fmt.Println("Opening your browser at:")
	fmt.Println(url)
	fmt.Println("Waiting for authentication...")

	jwt, err := callbackServer.Result()
	if err != nil {
		return suggestHeadless(cmd, err)
	}

	username, err := validateToken(jwt)
	if err != nil {
		return suggestHeadless(cmd, err)
	}

	config.SetToken(jwt)
	config.SetUsername(username)

	settings.PersistChanges()

	fmt.Printf("✔  Success! Logged in as %s\n", username)

	return nil
}

func randString(n int) string {
	runes := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}

func authURL(port int, path, state string) (string, error) {
	base, err := url.Parse(settings.GetA0URL())
	if err != nil {
		return "", fmt.Errorf("error parsing auth URL: %w", err)
	}
	authURL := base.JoinPath(path)

	values := url.Values{
		"redirect": {"false"},
	}
	if port != 0 {
		values = url.Values{
			"port":     {strconv.Itoa(port)},
			"redirect": {"true"},
			"type":     {"cli"},
		}
	}
	if state != "" {
		values["state"] = []string{state}
	}
	authURL.RawQuery = values.Encode()
	return authURL.String(), nil
}

type authCallback struct {
	ch     chan string
	server *http.Server
	Port   int
}

func authCallbackServer(state string) (authCallback, error) {
	ch := make(chan string, 1)
	server, err := createCallbackServer(ch, state)
	if err != nil {
		return authCallback{}, fmt.Errorf("cannot create callback server: %w", err)
	}

	port, err := runServer(server)
	if err != nil {
		return authCallback{}, fmt.Errorf("cannot run authentication server: %w", err)
	}

	return authCallback{
		ch:     ch,
		server: server,
		Port:   port,
	}, nil
}

func (a authCallback) Result() (string, error) {
	select {
	case result := <-a.ch:
		_ = a.server.Shutdown(context.Background())
		return result, nil
	case <-time.After(5 * time.Minute):
		_ = a.server.Shutdown(context.Background())
		return "", fmt.Errorf("authentication timed out, try again")
	}
}

func createCallbackServer(ch chan string, state string) (*http.Server, error) {
	// FIXME: @sanchitrk
	// make html template with path and embed in gobuild
	tmpl, err := template.New("login.html").Parse(loginHTML)
	if err != nil {
		return nil, fmt.Errorf("could not parse login callback template: %w", err)
	}
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("state") != state {
			w.WriteHeader(400)
			return
		}

		ch <- q.Get("jwt")

		w.WriteHeader(200)
		tmpl.Execute(w, map[string]string{
			"assetsURL": settings.GetA0URL(), // FIXME: @sanchitrk: later when you make the html
		})
	})

	return &http.Server{Handler: handler}, nil
}

func runServer(server *http.Server) (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, fmt.Errorf("could not allocate port for http server: %w", err)
	}

	go func() {
		server.Serve(listener)
	}()

	return listener.Addr().(*net.TCPAddr).Port, nil
}

func suggestHeadless(cmd *cobra.Command, err error) error {
	if err == nil {
		return nil
	}
	cmdWithFlag := cmd.CommandPath() + " --headless"
	return fmt.Errorf("%w\nIf the issue persists, try running %s", err, cli.Emph(cmdWithFlag))
}

func printHeadlessLoginInstructions(path string) error {
	url, err := authURL(0, path, "")
	if err != nil {
		return err
	}
	fmt.Println("Visit the following URL to login:")
	fmt.Println(url)
	return nil
}
