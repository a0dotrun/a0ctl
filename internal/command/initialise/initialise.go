package initialise

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/a0dotrun/a0ctl/internal/cli"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	const (
		use   = "init"
		short = "Initializes the server locally and provisions the server if it does not already exist"
	)
	cmd := &cobra.Command{
		Use:               use,
		Short:             short,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cli.NoFilesArg,
		RunE:              initFn,
	}
	return cmd
}

type Config struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Title       string            `yaml:"title"`
	Env         map[string]string `yaml:"env"`
}

func readLine(reader *bufio.Reader) string {
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func initFn(cmd *cobra.Command, _ []string) error {
	cmd.SilenceUsage = true

	reader := bufio.NewReader(os.Stdin)
	cfg := Config{}

	fmt.Print("Enter name: ")
	cfg.Name = readLine(reader)

	fmt.Print("Enter description: ")
	cfg.Description = readLine(reader)

	fmt.Print("Enter title: ")
	cfg.Title = readLine(reader)

	// Env key-value pairs
	cfg.Env = make(map[string]string)
	fmt.Println("Enter environment variables (key=value), blank line to finish:")
	for {
		line := readLine(reader)
		if line == "" {
			break
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			cfg.Env[parts[0]] = parts[1]
		} else {
			fmt.Println("Invalid format, use key=value")
		}
	}

	// Marshal to YAML
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		fmt.Println("Error marshaling YAML:", err)
		return err
	}

	// Save to file
	file, err := os.Create("config.yaml")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	fmt.Println("config.yaml created successfully!")
	return nil
}
