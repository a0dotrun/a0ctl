package initialize

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// readmeTxt contains the content for the README.txt file created in the .a0 directory.
// Modify the text as needed to provide relevant information to the user.
const readmeTxt = `> Why do I have a folder named ".a0" in my project?
The ".a0" folder is created when you link a directory to a a0 project.

> What does the "project.json" file contain?
The "project.json" file contains:
- The Name of the a0 project.
- The Region where the project is hosted.

> Should I commit the ".a0" folder?
No, you should not share the ".a0" folder with anyone.
`

func saveInitConfig(projectName, region string) error {
	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create .a0 directory in current directory
	a0Dir := filepath.Join(currentDir, ".a0")
	if err := os.MkdirAll(a0Dir, 0755); err != nil {
		return fmt.Errorf("failed to create .a0 directory: %w", err)
	}

	// Create project.json
	config := struct {
		ProjectName string `json:"name"`
		Region      string `json:"region"`
	}{
		ProjectName: projectName,
		Region:      region,
	}

	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filepath.Join(a0Dir, "project.json"), configJSON, 0644); err != nil {
		return fmt.Errorf("failed to write project.json: %w", err)
	}

	// Create README.txt
	if err := os.WriteFile(filepath.Join(a0Dir, "README.txt"), []byte(readmeTxt), 0644); err != nil {
		return fmt.Errorf("failed to write README.txt: %w", err)
	}

	fmt.Printf("Configuration saved in %s\n", a0Dir)
	return nil
}

var (
	projectName string
	region      string
)

func New() *cobra.Command {
	const (
		short = "Initialize a0ctl configuration"
		long  = "Initializes the a0ctl configuration with project name, region, etc."
	)
	cmd := &cobra.Command{
		Use:   "init",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create a new Huh form
			form := huh.NewForm(
				huh.NewGroup(
					// Ask for Project Name
					huh.NewInput().
						Title("What is your project name?").
						Value(&projectName). // Store the input in the projectName variable
						Validate(func(str string) error {
							if str == "" {
								return fmt.Errorf("project name cannot be empty")
							}
							if strings.Contains(str, " ") {
								return fmt.Errorf("project name cannot contain spaces")
							}
							return nil
						}).
						Description("Enter a unique name for your project."),

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
						Description("Select the primary region for your project resources."),
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
				if err == huh.ErrUserAborted { // Specific error for user abort (Ctrl+C)
					fmt.Fprintln(os.Stderr, "Initialization cancelled by user.")
				} else {
					fmt.Fprintln(os.Stderr, "Error during initialization:", err)
				}
				return err // Or return nil if user cancellation isn't considered a CLI error
			}

			// // If using confirmation:
			// if !confirmed {
			//  fmt.Println("Initialization aborted by user.")
			//  return nil
			// }

			// At this point, 'projectName' and 'region' variables are populated.
			// You would typically save these to a configuration file (e.g., a0ctl.yaml)
			saveInitConfig(projectName, region)

			return nil
		},
	}
	return cmd
}
