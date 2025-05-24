package initialize

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

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
			fmt.Printf("\n--- Configuration Summary ---\n")
			fmt.Printf("Project Name: %s\n", projectName)
			fmt.Printf("Region:       %s\n", region)
			fmt.Println("\nInitialization complete. You would typically save this to a config file.")
			// Example: saveConfig(projectName, region)

			return nil
		},
	}
	return cmd
}
