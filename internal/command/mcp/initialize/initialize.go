package mcpinitialize

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/a0dotrun/a0ctl/internal/appconfig"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// readmeTxt contains the content for the README.txt file created in the .a0 directory.
const readmeTxt = `> Why do I have a folder named ".a0" in my app?
The ".a0" folder is created when you link a directory to a a0 app.

> What does the "app.json" file contain?
The "app.json" file contains:
- The Name of the a0 app.
- The Region where the app is hosted.

> Should I commit the ".a0" folder?
No, you should not share the ".a0" folder with anyone.
`

func saveInitConfig(appName, region string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create .a0 directory in current directory
	a0Dir := filepath.Join(currentDir, ".a0")
	if err := os.MkdirAll(a0Dir, 0755); err != nil {
		return fmt.Errorf("failed to create .a0 directory: %w", err)
	}

	config := appconfig.NewConfig(appName, region)
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filepath.Join(a0Dir, "app.json"), configJSON, 0644); err != nil {
		return fmt.Errorf("failed to write app.json: %w", err)
	}

	// Create README.txt
	if err := os.WriteFile(filepath.Join(a0Dir, "README.txt"), []byte(readmeTxt), 0644); err != nil {
		return fmt.Errorf("failed to write README.txt: %w", err)
	}

	fmt.Printf("Configuration saved in %s\n", a0Dir)
	return nil
}

var (
	appName string
	region  string
)

func New() *cobra.Command {
	const (
		short = "Initialize a0ctl configuration"
		long  = "Initializes the a0ctl configuration with app name, region, etc."
	)
	cmd := &cobra.Command{
		Use:   "init",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			_ = ctx // Access the context here
			// Create a new Huh form
			form := huh.NewForm(
				huh.NewGroup(
					// Ask for Project Name
					huh.NewInput().
						Title("What is your app name?").
						Value(&appName). // Store the input in the appName variable
						// TODO:@sanchitrk
						// improve the validation
						// Return types for Errors?
						Validate(func(str string) error {
							if str == "" {
								return fmt.Errorf("name cannot be empty")
							}
							if strings.Contains(str, " ") {
								return fmt.Errorf("name cannot contain spaces")
							}
							return nil
						}).
						Description("Enter a name for your app."),
					// Ask for Region
					huh.NewSelect[string]().
						Title("Choose a region:").
						Options(
							huh.NewOption("US East (N. Virginia)", "us-east-1"),
							huh.NewOption("US West (Oregon)", "us-west-2"),
							huh.NewOption("EU (Frankfurt)", "eu-central-1"),
							huh.NewOption("Asia Pacific (Mumbai)", "ap-south-1"),
							huh.NewOption("Asia Pacific (Sydney)", "ap-southeast-2"),
						).
						Value(&region). // Store the selected option in the region variable
						Description("Select the primary region for your app resources."),
				),
				// You can add more groups or fields here
				// For example, a confirmation step:
				// huh.NewConfirm().
				//  Title("Proceed with these settings?").
				//  Affirmative("Yes, initialize!").
				//  Negative("No, cancel.").
				//  Value(&confirmed), // Assuming 'var confirmed bool'
			)

			// Run the form
			err := form.Run()

			// Handle errors (e.g., user exits with Ctrl+C)
			if err != nil {
				if errors.Is(err, huh.ErrUserAborted) { // Specific error for user abort (Ctrl+C)
					_, err := fmt.Fprintln(os.Stderr, "Initialization cancelled by user.")
					if err != nil {
						return err
					}
				} else {
					_, err := fmt.Fprintln(os.Stderr, "Error during initialization:", err)
					if err != nil {
						return err
					}
				}
				return err // Or return nil if user cancellation isn't considered a CLI error
			}

			// // If using confirmation:
			// if !confirmed {
			//  fmt.Println("Initialization aborted by user.")
			//  return nil
			// }

			// At this point, 'appName' and 'region' variables are populated.
			// You would typically save these to a configuration file (e.g., a0ctl.yaml)
			err = saveInitConfig(appName, region)
			if err != nil {
				_, err := fmt.Fprintln(os.Stderr, "Error saving configuration:", err)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
	return cmd
}
