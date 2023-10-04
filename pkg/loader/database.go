package loader

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Name     string
}

func NewDatabaseConfig(user string, password string, host string, name string) *DatabaseConfig {
	return &DatabaseConfig{
		User:     user,
		Password: password,
		Host:     host,
		Name:     name,
	}
}

func CopyDatabase(target *DatabaseConfig, source *DatabaseConfig) error {
	tempDir := os.TempDir()
	dumpFileName := filepath.Join(tempDir, "dump.sql")
	// Dump the remote database to a file
	dumpCmd := exec.Command("mysqldump", "-u"+source.User, "-p"+source.Password, "-h"+source.Host, source.Name, "> "+dumpFileName)
	if err := dumpCmd.Run(); err != nil {
		fmt.Println("Error dumping remote database:", err)
		return err
	}

	fmt.Println("Remote database dumped to", dumpFileName)

	// Import the dump into the local database
	importCmd := exec.Command("mysql", "-u"+target.User, "-p"+target.Password, "-h"+source.Host, source.Name, target.Name, "< "+dumpFileName)

	if err := importCmd.Run(); err != nil {
		fmt.Println("Error importing data to database:", err)
		return err
	}

	fmt.Println("Dump imported into local database")

	// Optional: You can remove the dump file after importing
	if err := os.Remove(dumpFileName); err != nil {
		fmt.Println("Error removing dump file:", err)
		return err
	}
	return nil
}
