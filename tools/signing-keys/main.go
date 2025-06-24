package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Keys struct {
	Keys []*Key `json:"keys"`
}

type Key struct {
	PublicKey string `json:"publicKey"`
}

func main() {

	var insertInto string
	flag.StringVar(&insertInto, "insert", "", "Insert the keys into the specified .env file")

	flag.Parse()
	publicKeyFiles := flag.Args()
	if len(publicKeyFiles) == 0 {
		fmt.Println("No public key files provided.")
		return
	}

	publicKeys := &Keys{
		Keys: []*Key{},
	}
	for _, publicKeyFile := range publicKeyFiles {
		data, err := os.ReadFile(publicKeyFile)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", publicKeyFile, err)
			continue
		}
		publicKey := string(data)
		publicKeys.Keys = append(publicKeys.Keys, &Key{
			PublicKey: publicKey,
		})
	}

	if len(publicKeys.Keys) == 0 {
		fmt.Println("No valid public keys found.")
		return
	}

	bytes, err := json.Marshal(publicKeys)
	if err != nil {
		fmt.Printf("Error marshaling keys to JSON: %v\n", err)
		return
	}

	envVarLine := fmt.Sprintf(
		"BLUELINK_GITHUB_REGISTRY_SIGNING_PUBLIC_KEYS='%s'\n",
		string(bytes),
	)

	if insertInto != "" {
		envFileData, err := os.ReadFile(insertInto)
		if err != nil {
			fmt.Printf("Error reading .env file %s: %v\n", insertInto, err)
			return
		}

		lines := strings.Split(string(envFileData), "\n")

		existingLineIndex := slices.IndexFunc(lines, func(line string) bool {
			return strings.HasPrefix(
				strings.TrimSpace(line),
				"BLUELINK_GITHUB_REGISTRY_SIGNING_PUBLIC_KEYS=",
			)
		})

		if existingLineIndex != -1 {
			lines[existingLineIndex] = envVarLine
		} else {
			lines = append(lines, envVarLine)
		}

		finalLines := removeEmptyLines(lines)

		err = os.WriteFile(insertInto, []byte(strings.Join(finalLines, "\n")), 0644)
		if err != nil {
			fmt.Printf("Error writing to .env file %s: %v\n", insertInto, err)
		}
		return
	}

	fmt.Println("Add the following to your .env file:")
	fmt.Println(envVarLine)
}

func removeEmptyLines(lines []string) []string {
	var nonEmptyLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}
	return nonEmptyLines
}
