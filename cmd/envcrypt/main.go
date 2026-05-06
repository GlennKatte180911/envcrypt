// Command envcrypt provides a CLI for encrypting and managing .env files
// for teams using symmetric keys.
//
// Usage:
//
//	envcrypt <command> [flags]
//
// Commands:
//
//	lock    Encrypt a .env file into a vault
//	unlock  Decrypt a vault file into a .env file
//	keygen  Generate a new encryption key file
//	diff    Show differences between two .env files
//	audit   Audit a .env file for common issues
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envcrypt/internal/vault"
)

const usage = `envcrypt — encrypt and manage .env files for teams

Usage:
  envcrypt <command> [flags]

Commands:
  lock    Encrypt a .env file and write to a vault file
  unlock  Decrypt a vault file and write to a .env file
  keygen  Generate a new encryption key file

Flags:
  -env    Path to the .env file       (default: .env)
  -vault  Path to the vault file      (default: .env.vault)
  -key    Path to the key file        (default: ~/.envcrypt/key)
  -help   Show this help message

Examples:
  envcrypt keygen
  envcrypt lock   -env .env -vault .env.vault
  envcrypt unlock -vault .env.vault -env .env.decrypted
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "lock":
		runLock(args)
	case "unlock":
		runUnlock(args)
	case "keygen":
		runKeygen(args)
	case "-help", "--help", "help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "envcrypt: unknown command %q\n\n", cmd)
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}
}

func runLock(args []string) {
	fs := flag.NewFlagSet("lock", flag.ExitOnError)
	envPath := fs.String("env", ".env", "path to the .env file")
	vaultPath := fs.String("vault", ".env.vault", "path to the output vault file")
	keyPath := fs.String("key", vault.DefaultKeyPath(), "path to the key file")
	_ = fs.Parse(args)

	key, err := vault.LoadKeyFile(*keyPath)
	if err != nil {
		fatalf("load key: %v", err)
	}

	v := vault.New(*vaultPath, key)
	data, err := os.ReadFile(*envPath)
	if err != nil {
		fatalf("read env file: %v", err)
	}

	if err := v.Lock(string(data)); err != nil {
		fatalf("lock: %v", err)
	}
	fmt.Printf("locked %s → %s\n", *envPath, *vaultPath)
}

func runUnlock(args []string) {
	fs := flag.NewFlagSet("unlock", flag.ExitOnError)
	vaultPath := fs.String("vault", ".env.vault", "path to the vault file")
	envPath := fs.String("env", ".env.decrypted", "path to the output .env file")
	keyPath := fs.String("key", vault.DefaultKeyPath(), "path to the key file")
	_ = fs.Parse(args)

	key, err := vault.LoadKeyFile(*keyPath)
	if err != nil {
		fatalf("load key: %v", err)
	}

	v := vault.New(*vaultPath, key)
	plaintext, err := v.Unlock()
	if err != nil {
		fatalf("unlock: %v", err)
	}

	if err := os.WriteFile(*envPath, []byte(plaintext), 0o600); err != nil {
		fatalf("write env file: %v", err)
	}
	fmt.Printf("unlocked %s → %s\n", *vaultPath, *envPath)
}

func runKeygen(args []string) {
	fs := flag.NewFlagSet("keygen", flag.ExitOnError)
	keyPath := fs.String("key", vault.DefaultKeyPath(), "path to write the generated key")
	_ = fs.Parse(args)

	if err := vault.GenerateKeyFile(*keyPath); err != nil {
		fatalf("keygen: %v", err)
	}
	fmt.Printf("key written to %s\n", *keyPath)
}

func fatalf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "envcrypt: "+format+"\n", a...)
	os.Exit(1)
}
