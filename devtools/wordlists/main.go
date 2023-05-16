// Tool wordlists generates all wordlists to be embedded in the machine-name tool for generating
// localized names.
package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/PlanktoScope/machine-name/internal/wordlists/sources"
	"gopkg.in/yaml.v3"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Println("Resetting the directory of generated wordlists...")
	outputDir := filepath.Join(cwd, "internal", "wordlists", "generated")
	if err = cleanGenerated(outputDir); err != nil {
		panic(err)
	}

	fmt.Println("Generating wordlists from sources...")
	if err = generateAll(sources.WordlistsFS, outputDir); err != nil {
		panic(err)
	}

	fmt.Println("Done!")
}

func cleanGenerated(dir string) error {
	file, err := os.Open(dir)
	if err != nil {
		return errors.Wrap(err, "couldn't open directory containing generated wordlists")
	}
	dirEntries, err := file.ReadDir(0)
	if err != nil {
		return errors.Wrap(err, "couldn't list subdirectories of the generated wordlist directory")
	}
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			continue
		}
		subdir := filepath.Join(dir, entry.Name())
		fmt.Printf("  Deleting '%s'...\n", subdir)
		if err := os.RemoveAll(subdir); err != nil {
			return errors.Wrapf(err, "couldn't remove subdirectory %s", subdir)
		}
	}
	return nil
}

func generateAll(sourcesDir fs.FS, generatedDir string) error {
	dirEntries, err := fs.ReadDir(sourcesDir, ".")
	if err != nil {
		return errors.Wrap(err, "couldn't scan directories for wordlist sources")
	}
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			continue
		}
		language := entry.Name()
		fmt.Printf("Generating word lists for language '%s'...\n", language)
		sourcesDir, err := fs.Sub(sourcesDir, language)
		if err != nil {
			return errors.Wrapf(err, "couldn't open wordlist source directory '%s'", language)
		}
		if err := generate(sourcesDir, filepath.Join(generatedDir, language)); err != nil {
			return errors.Wrapf(
				err, "couldn't generate word lists for language %s in %s",
				language, filepath.Join(generatedDir, language),
			)
		}
	}
	return nil
}

func generate(sources fs.FS, generatedDir string) error {
	config, err := loadGenerationConfig(sources)
	if err != nil {
		return errors.Wrap(err, "couldn't load config file for wordlist generation")
	}
	fmt.Println("  Generating first wordlist...")
	firstWordList, err := makeWordList(sources, config.First)
	if err != nil {
		return errors.Wrap(err, "couldn't make first word list")
	}
	fmt.Printf("    Generated %d words!\n", len(firstWordList))
	fmt.Printf("    Saving to %s/first.txt...\n", generatedDir)
	if err = saveWordList(generatedDir, firstWordList, "first.txt"); err != nil {
		return errors.Wrap(err, "couldn't save generated first word list")
	}
	fmt.Println("  Generating second wordlist...")
	secondWordList, err := makeWordList(sources, config.Second)
	if err != nil {
		return errors.Wrap(err, "couldn't make second word list")
	}
	fmt.Printf("    Generated %d words!\n", len(secondWordList))
	fmt.Printf("    Saving to %s/second.txt...\n", generatedDir)
	if err = saveWordList(generatedDir, secondWordList, "second.txt"); err != nil {
		return errors.Wrap(err, "couldn't save generated second word list")
	}
	const uint32max = 4294967295
	digits := math.Ceil(math.Log10(
		uint32max / float64(len(firstWordList)) / float64(len(secondWordList)),
	))
	fmt.Printf(
		"  A 32-bit serial number will require a %.f-digit number at the end of the name.\n", digits,
	)
	return nil
}

func loadGenerationConfig(sources fs.FS) (GenerationConfig, error) {
	contents, err := fs.ReadFile(sources, "config.yml")
	if err != nil {
		return GenerationConfig{}, errors.Wrap(err, "couldn't open config file")
	}
	config := GenerationConfig{}
	if err = yaml.Unmarshal(contents, &config); err != nil {
		return GenerationConfig{}, errors.Wrap(err, "couldn't parse generation configuration")
	}
	return config, nil
}

func makeWordList(sources fs.FS, spec WordlistSpec) ([]string, error) {
	sourceWords, err := loadWords(sources, spec.Sources)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't load source words")
	}
	filterWords, err := loadWords(sources, spec.Filters)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't load filter words")
	}
	capacity := len(sourceWords)
	if len(filterWords) < capacity {
		capacity = len(filterWords)
	}
	intersection := make([]string, 0, capacity)
	for word := range sourceWords {
		if _, ok := filterWords[word]; !ok {
			continue
		}
		intersection = append(intersection, word)
	}
	sort.Strings(intersection)
	return intersection, nil
}

func loadWords(sources fs.FS, files []string) (map[string]struct{}, error) {
	allWords := make(map[string]struct{})
	for _, file := range files {
		contents, err := fs.ReadFile(sources, file)
		if err != nil {
			return nil, errors.Wrapf(err, "couldn't read file '%s'", file)
		}
		fileWords := strings.Split(string(contents), "\n")
		for _, word := range fileWords {
			if word == "" {
				continue
			}
			if strings.HasPrefix(word, "# ") {
				continue
			}
			allWords[word] = struct{}{}
		}
	}
	return allWords, nil
}

func saveWordList(outputDir string, wordList []string, filename string) error {
	const perm = 0o777 // allow everything
	os.MkdirAll(outputDir, perm)
	outputFile := filepath.Join(outputDir, filename)
	file, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return errors.Wrapf(err, "couldn't open file '%s' for writing", outputFile)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, word := range wordList {
		if _, err := writer.WriteString(fmt.Sprintf("%s\n", word)); err != nil {
			return errors.Wrapf(err, "couldn't write word '%s' to file '%s'", word, outputFile)
		}
	}
	writer.Flush()
	return nil
}
