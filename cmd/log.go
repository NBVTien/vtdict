package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/NBVTien/vtdict/internal/storage"
	"github.com/NBVTien/vtdict/internal/ui"
	"github.com/spf13/cobra"
)

var (
	logLimit    int
	clearFlag   bool
	clearAll    bool
	clearBefore string
	clearWord   string
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show or clear lookup history",
	Example: `  vtdict log
  vtdict log -n 50
  vtdict log --clear              (interactive picker)
  vtdict log --clear --all        (wipe everything)
  vtdict log --clear --before 7d
  vtdict log --clear --word hello`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if clearWord != "" {
			found, err := storage.ClearWord(clearWord)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			if found {
				fmt.Printf("Removed \"%s\" from history.\n", clearWord)
			} else {
				fmt.Printf("\"%s\" not in history.\n", clearWord)
			}
			return nil
		}

		if clearBefore != "" {
			d, err := parseDuration(clearBefore)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid duration %q — use e.g. 7d, 2w, 1m\n", clearBefore)
				os.Exit(1)
			}
			n, err := storage.ClearBefore(d)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Removed %d entries older than %s.\n", n, clearBefore)
			return nil
		}

		if clearAll {
			if !clearFlag {
				fmt.Fprintf(os.Stderr, "use --clear --all to wipe history\n")
				os.Exit(1)
			}
			if err := storage.ClearHistory(); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("History cleared.")
			return nil
		}

		if clearFlag {
			entries, err := storage.GetHistory(1000)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			words, err := ui.PickWordsToDelete(entries)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			if len(words) == 0 {
				fmt.Println("Cancelled.")
				return nil
			}
			n, err := storage.ClearWords(words)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Removed %d word(s).\n", n)
			return nil
		}

		entries, err := storage.GetHistory(logLimit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		ui.RenderHistory(entries)
		return nil
	},
}

func init() {
	logCmd.Flags().IntVarP(&logLimit, "limit", "n", 20, "number of entries to show")
	logCmd.Flags().BoolVar(&clearFlag, "clear", false, "interactive picker to remove words")
	logCmd.Flags().BoolVar(&clearAll, "all", false, "with --clear: wipe entire history")
	logCmd.Flags().StringVar(&clearBefore, "before", "", "with --clear: remove entries older than duration (e.g. 7d, 2w, 1m)")
	logCmd.Flags().StringVar(&clearWord, "word", "", "with --clear: remove a specific word")
	rootCmd.AddCommand(logCmd)
}

func parseDuration(s string) (time.Duration, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if len(s) < 2 {
		return 0, fmt.Errorf("too short")
	}
	suffix := s[len(s)-1]
	numStr := s[:len(s)-1]
	n, err := strconv.Atoi(numStr)
	if err != nil {
		return time.ParseDuration(s)
	}
	switch suffix {
	case 'd':
		return time.Duration(n) * 24 * time.Hour, nil
	case 'w':
		return time.Duration(n) * 7 * 24 * time.Hour, nil
	case 'm':
		return time.Duration(n) * 30 * 24 * time.Hour, nil
	case 'h':
		return time.Duration(n) * time.Hour, nil
	}
	return time.ParseDuration(s)
}
