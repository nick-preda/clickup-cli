package list

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
	"github.com/triptechtravel/clickup-cli/internal/api"
	"github.com/triptechtravel/clickup-cli/pkg/cmdutil"
)

// NewCmdListMove returns the "list move" command.
func NewCmdListMove(f *cmdutil.Factory) *cobra.Command {
	var folderID string

	cmd := &cobra.Command{
		Use:   "move <list-id>",
		Short: "Move a list into a folder",
		Long: `Move a ClickUp list into a folder (or back to folderless).

Use --folder to specify the target folder ID. The list retains all its
tasks, statuses, and settings — only its parent container changes.`,
		Example: `  # Move a list into a folder
  clickup list move 901522563411 --folder 901515444234

  # Move multiple lists into a folder
  clickup list move 123 456 789 --folder 901515444234`,
		Args:    cobra.MinimumNArgs(1),
		PreRunE: cmdutil.NeedsAuth(f),
		RunE: func(cmd *cobra.Command, args []string) error {
			if folderID == "" {
				return fmt.Errorf("--folder is required")
			}

			client, err := f.ApiClient()
			if err != nil {
				return err
			}

			ios := f.IOStreams

			for _, listID := range args {
				err := moveList(cmd.Context(), client, listID, folderID)
				if err != nil {
					fmt.Fprintf(ios.ErrOut, "✗ Failed to move list %s: %v\n", listID, err)
					continue
				}
				fmt.Fprintf(ios.Out, "✓ Moved list %s into folder %s\n", listID, folderID)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&folderID, "folder", "", "Target folder ID (required)")
	_ = cmd.MarkFlagRequired("folder")

	return cmd
}

func moveList(ctx context.Context, client *api.Client, listID, folderID string) error {
	apiURL := fmt.Sprintf("https://api.clickup.com/api/v2/list/%s", url.PathEscape(listID))

	body := fmt.Sprintf(`{"folder_id":"%s"}`, folderID)
	req, err := http.NewRequestWithContext(ctx, "PUT", apiURL, strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.DoRequest(req)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		var errResp struct {
			Err string `json:"err"`
		}
		if json.Unmarshal(respBody, &errResp) == nil && errResp.Err != "" {
			return fmt.Errorf("%s", errResp.Err)
		}
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	return nil
}
