package cmd

import (
	"fmt"
	"os"

	"github.com/NBVTien/vtdict/internal/config"
	"github.com/NBVTien/vtdict/internal/dictionary"
	"github.com/NBVTien/vtdict/internal/storage"
	"github.com/NBVTien/vtdict/internal/translate"
	"github.com/NBVTien/vtdict/internal/ui"
)

var (
	flagNoTranslate bool
	flagLang        string
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&flagNoTranslate, "no-translate", false, "skip translation (overrides config)")
	rootCmd.PersistentFlags().StringVarP(&flagLang, "lang", "l", "", "target language code, e.g. vi, fr, ja (overrides config)")
	rootCmd.PersistentFlags().StringVarP(&searchWord, "search", "s", "", "look up a word (works on any subcommand)")
}

func lookupWord(word string) error {
	cfg := config.Get()

	// flag overrides config
	doTranslate := cfg.Translate
	if flagNoTranslate {
		doTranslate = false
	}
	lang := cfg.Lang
	if flagLang != "" {
		lang = flagLang
	}

	results, err := dictionary.Lookup(word)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// API returns duplicate entries for same word — show only first
	if len(results) > 1 {
		results = results[:1]
	}

	translation := ""
	if doTranslate {
		def := dictionary.FirstDefinition(results)
		if def != "" {
			translation, _ = translate.Translate(def, lang)
		}
	}

	ui.RenderLookup(results, translation, ui.RenderOpts{
		Phonetic: cfg.Phonetic,
		Examples: cfg.Examples,
		POS:      cfg.POS,
	})

	firstDef := dictionary.FirstDefinition(results)
	if err := storage.Save(word, firstDef, translation); err != nil {
		fmt.Fprintf(os.Stderr, "warn: could not save to history: %v\n", err)
	}

	return nil
}
