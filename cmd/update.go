package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"resty.dev/v3"
	"github.com/spf13/cobra"
)

const (
	githubOwner = "NBVTien"
	githubRepo  = "vtdict"
)

var currentVersion = "dev" // replaced by goreleaser at build time via ldflags

type githubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update vtdict to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Checking for updates...")

		client := resty.New()
		var release githubRelease

		_, err := client.R().
			SetResult(&release).
			SetHeader("Accept", "application/vnd.github+json").
			Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", githubOwner, githubRepo))
		if err != nil {
			return fmt.Errorf("failed to check releases: %w", err)
		}

		if release.TagName == "" {
			return fmt.Errorf("no releases found — make sure the repo has a published release")
		}

		if release.TagName == currentVersion {
			fmt.Printf("Already up to date (%s)\n", currentVersion)
			return nil
		}

		fmt.Printf("New version available: %s (current: %s)\n", release.TagName, currentVersion)

		assetName := fmt.Sprintf("vtdict_%s_%s", runtime.GOOS, runtime.GOARCH)
		if runtime.GOOS == "windows" {
			assetName += ".exe"
		}

		var downloadURL string
		for _, a := range release.Assets {
			if a.Name == assetName {
				downloadURL = a.BrowserDownloadURL
				break
			}
		}
		if downloadURL == "" {
			return fmt.Errorf("no binary found for %s/%s in release %s\nasset expected: %s",
				runtime.GOOS, runtime.GOARCH, release.TagName, assetName)
		}

		fmt.Printf("Downloading %s...\n", assetName)

		httpResp, err := http.Get(downloadURL)
		if err != nil {
			return fmt.Errorf("download failed: %w", err)
		}
		defer httpResp.Body.Close()

		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("cannot find current executable: %w", err)
		}

		// write to temp file first, then atomically replace
		tmp := exe + ".new"
		f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return fmt.Errorf("cannot write update: %w", err)
		}

		if _, err := io.Copy(f, httpResp.Body); err != nil {
			f.Close()
			os.Remove(tmp)
			return fmt.Errorf("download failed: %w", err)
		}
		f.Close()

		if err := os.Rename(tmp, exe); err != nil {
			os.Remove(tmp)
			return fmt.Errorf("failed to replace binary: %w", err)
		}

		fmt.Printf("Updated to %s\n", release.TagName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
