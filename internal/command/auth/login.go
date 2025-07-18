package auth

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/a0dotrun/a0ctl/internal/cli"
	"github.com/a0dotrun/a0ctl/internal/settings"

	"github.com/a0dotrun/a0ctl/internal/api"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var LOGIN_HTML string

const (
	a0DefaultBaseURL = "https://api.a0.tech"
)

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
	}
	return cmd
}

func login(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	state := randString(32)
	callbackServer, err := authCallbackServer(state)
	if err != nil {
		return suggestHeadless(cmd, err)
	}

	url, err := authURL(callbackServer.Port, "", state)
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

	if !api.IsJWTTokenValid(jwt) {
		return errors.New("invalid token")
	}

	client, err := api.AuthedClient()
	if err != nil {
		return err
	}

	user, err := client.Users.GetUser()
	if err != nil {
		return err
	}

	config, err := settings.ReadSettings()
	if err != nil {
		return fmt.Errorf("failed to read settings: %w", err)
	}

	config.SetToken(jwt)
	if err := settings.TryToPersistChanges(); err != nil {
		return fmt.Errorf("%w\nIf the issue persists, set your token to the %s environment variable instead", err, cli.Emph(settings.EnvAccessToken))
	}

	config.SetUsername(user.Username)

	fmt.Printf("✔  Success! Logged in as %s\n", user.Username)

	return nil
}

func randString(n int) string {
	var runes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}

func authURL(port int, path, state string) (string, error) {
	base, err := url.Parse(getA0Url())
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
	tmpl, err := template.New("login.html").Parse(LOGIN_HTML)
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
			"assetsURL": getA0Url(),
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

func getA0Url() string {
	config, _ := settings.ReadSettings() // ok to ignore, we'll fallback to default
	url := config.GetBaseURL()
	if url == "" {
		url = a0DefaultBaseURL
	}
	return url
}
