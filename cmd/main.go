package main

import (
	"flag"
	"fmt"
	"os"

	"go.step.sm/crypto/minica"
	"go.step.sm/crypto/pemutil"
)

func main() {
	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	output := generateCmd.String("output", "tpm-ca-certificates.pem", "Output file for the generated bundle")

	switch subcmd := os.Args[1]; subcmd {
	case "generate":
		generateCmd.Parse(os.Args[2:])
		if err := createBundle(*output); err != nil {
			fmt.Printf("Error creating bundle: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Bundle generate successfully 🚀")
	default:
		fmt.Println("Unknown subcommand. Expected 'generate'")
		os.Exit(1)
	}
}

func createBundle(output string) error {
	localCA, err := minica.New()
	if err != nil {
		panic(err)
	}
	_, err = pemutil.Serialize(localCA.Root, pemutil.ToFile(output, 0644))
	return err
}
