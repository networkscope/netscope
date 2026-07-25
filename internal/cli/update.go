package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const githubAPI = "https://api.github.com/repos/networkscope/netscope/releases/latest"

type release struct {
	TagName string `json:"tag_name"`
}

func checkForUpdate() {
	if version == "" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubAPI, nil)
	if err != nil {
		return
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "netscope-auto-update")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var rel release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return
	}

	latest := strings.TrimPrefix(rel.TagName, "v")
	current := strings.TrimPrefix(version, "v")
	if latest != "" && latest != current {
		fmt.Fprintf(os.Stderr, "A new version is available: %s (current: %s)\n", rel.TagName, version)
		fmt.Fprintf(os.Stderr, "Download from: https://github.com/networkscope/netscope/releases/latest\n")
	}
}

func autoUpdatePreRun(cmd *cobra.Command, args []string) {
	checkForUpdate()
}
