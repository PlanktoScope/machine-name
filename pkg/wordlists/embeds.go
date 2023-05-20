// Package wordlists provides embedded wordlists which can be used for constructing names.
package wordlists

import (
	"embed"
	"io/fs"
	"strings"

	"github.com/pkg/errors"
)

//go:embed *
var FS embed.FS

func ListLanguages(f fs.FS) (map[string]struct{}, error) {
	dirEntries, err := fs.ReadDir(f, ".")
	if err != nil {
		return nil, errors.Wrap(err, "couldn't list subdirectories of the directory of wordlists")
	}
	languages := make(map[string]struct{})
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			continue
		}
		languages[entry.Name()] = struct{}{}
	}
	return languages, nil
}

func Load(f fs.FS, lang string) (first []string, second []string, err error) {
	lfs, err := fs.Sub(f, lang)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "couldn't open subdirectory '%s'", lang)
	}
	if first, err = loadWordlist(lfs, "first.txt"); err != nil {
		return nil, nil, errors.Wrapf(err, "couldn't load '%s/first.txt'", lang)
	}
	if second, err = loadWordlist(lfs, "second.txt"); err != nil {
		return nil, nil, errors.Wrapf(err, "couldn't load '%s/second.txt'", lang)
	}
	return first, second, nil
}

func loadWordlist(lang fs.FS, file string) ([]string, error) {
	contents, err := fs.ReadFile(lang, file)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't read file '%s'", file)
	}
	fileWords := strings.Split(string(contents), "\n")
	words := make([]string, 0, len(fileWords))
	for _, word := range fileWords {
		if word == "" {
			continue
		}
		if strings.HasPrefix(word, "# ") {
			continue
		}
		words = append(words, word)
	}
	return words, nil
}
