package cmd

import (
	"fmt"
	"os"

	"github.com/NBVTien/vtdict/internal/ai"
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

	doTranslate := cfg.Translate
	if flagNoTranslate {
		doTranslate = false
	}
	lang := cfg.Lang
	if flagLang != "" {
		lang = flagLang
	}

	opts := ui.RenderOpts{
		Phonetic: cfg.Phonetic,
		Examples: cfg.Examples,
		POS:      cfg.POS,
	}

	// check SQLite cache first
	if cached, err := storage.Get(word); err == nil && cached != nil {
		if err := storage.Save(word, cached.Definition, cached.Translation); err != nil {
			fmt.Fprintf(os.Stderr, "warn: could not update history: %v\n", err)
		}
		ui.RenderCached(cached.Definition, cached.Translation, opts)
		return nil
	}

	results, dictErr := dictionary.Lookup(word)

	// AI fallback: word not found in dictionary
	if dictErr != nil {
		if cfg.AIFallback {
			def, err := ai.LookupWord(word)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: word not found, AI fallback failed: %v\n", err)
				os.Exit(1)
			}
			translation := ""
			if doTranslate && def.Definition != "" {
				translation, _ = translate.Translate(def.Definition, lang)
			}
			ui.RenderAILookup(def, translation, opts)
			if err := storage.Save(word, def.Definition, translation); err != nil {
				fmt.Fprintf(os.Stderr, "warn: could not save to history: %v\n", err)
			}
			return nil
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", dictErr)
		os.Exit(1)
	}

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

	ui.RenderLookup(results, translation, opts)

	firstDef := dictionary.FirstDefinition(results)
	if err := storage.Save(word, firstDef, translation); err != nil {
		fmt.Fprintf(os.Stderr, "warn: could not save to history: %v\n", err)
	}

	return nil
}
