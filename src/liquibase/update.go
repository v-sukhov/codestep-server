package liquibase

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	liquibaseDir   = "liquibase"
	changeLogFile  = "dbchangelog.xml"
	changesetsDir  = "changesets"
	propertiesFile = "liquibase.properties"
)

/*
Update выполняет liquibase update в каталоге liquibase относительно текущей рабочей директории.
Параметры подключения к БД берутся из liquibase/liquibase.properties.
*/
func Update() error {
	if err := validateFiles(); err != nil {
		return err
	}

	cmd := exec.Command("liquibase", "update")
	cmd.Dir = liquibaseDir

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		if output.Len() > 0 {
			return fmt.Errorf("liquibase update failed: %w\n%s", err, output.String())
		}
		return fmt.Errorf("liquibase update failed: %w", err)
	}

	if output.Len() > 0 {
		log.Printf("Liquibase update completed successfully:\n%s", output.String())
	} else {
		log.Print("Liquibase update completed successfully")
	}

	return nil
}

func validateFiles() error {
	changeLogPath := filepath.Join(liquibaseDir, changeLogFile)
	if _, err := os.Stat(changeLogPath); err != nil {
		return fmt.Errorf("liquibase changelog file not found: %s", changeLogPath)
	}

	changesetsPath := filepath.Join(liquibaseDir, changesetsDir)
	entries, err := os.ReadDir(changesetsPath)
	if err != nil {
		return fmt.Errorf("liquibase changesets directory not found: %s", changesetsPath)
	}
	if len(entries) == 0 {
		return fmt.Errorf("liquibase changesets directory is empty: %s", changesetsPath)
	}

	propertiesPath := filepath.Join(liquibaseDir, propertiesFile)
	if _, err := os.Stat(propertiesPath); err != nil {
		return fmt.Errorf("liquibase properties file not found: %s", propertiesPath)
	}

	if _, err := exec.LookPath("liquibase"); err != nil {
		return fmt.Errorf("liquibase command not found in PATH: %w", err)
	}

	return nil
}
