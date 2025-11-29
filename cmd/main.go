package main

import (
	"encoding/pem"
	"flag"
	"fmt"
	"os"

	"github.com/loicsikidi/test-hybrid-release/internal/version"
	"go.step.sm/crypto/minica"
	"go.step.sm/crypto/pemutil"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: awesomecli [version|generate]")
		os.Exit(1)
	}

	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	output := generateCmd.String("output", "", "Output file for the generated bundle (default: stdout)")

	switch subcmd := os.Args[1]; subcmd {
	case "version":
		fmt.Println(version.Get())
	case "generate":
		generateCmd.Parse(os.Args[2:])
		if err := createBundle(*output); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating bundle: %v\n", err)
			os.Exit(1)
		}
		if *output == "" {
			fmt.Fprintln(os.Stderr, "Bundle generated successfully to stdout 🚀")
		} else {
			fmt.Fprintf(os.Stderr, "Bundle generated successfully to %s 🚀\n", *output)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand '%s'. Expected 'version' or 'generate'\n", subcmd)
		os.Exit(1)
	}
}

func createBundle(output string) error {
	localCA, err := minica.New()
	if err != nil {
		return fmt.Errorf("error creating CA: %w", err)
	}

	if output == "" {
		// Write to stdout
		block, err := pemutil.Serialize(localCA.Root)
		if err != nil {
			return fmt.Errorf("error serializing certificate: %w", err)
		}
		_, err = os.Stdout.Write(pem.EncodeToMemory(block))
		return err
	}

	// Write to file
	_, err = pemutil.Serialize(localCA.Root, pemutil.ToFile(output, 0644))
	return err
}
