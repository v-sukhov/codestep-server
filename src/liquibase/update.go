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
Параметры подключения к БД берутся из server.conf и записываются в liquibase/liquibase.properties.
jdbcDriverFile задаёт путь к JDBC-драйверу относительно текущей рабочей директории сервера.
*/
func Update(host, port, dbname, user, password, jdbcDriverFile string) error {
	if jdbcDriverFile == "" {
		return fmt.Errorf("jdbc_driver_file is not configured")
	}

	driverClasspath, err := driverClasspathForLiquibase(jdbcDriverFile)
	if err != nil {
		return err
	}

	if err := validateFiles(jdbcDriverFile); err != nil {
		return err
	}

	if err := writeProperties(host, port, dbname, user, password, driverClasspath); err != nil {
		return fmt.Errorf("failed to write %s: %w", filepath.Join(liquibaseDir, propertiesFile), err)
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

func driverClasspathForLiquibase(jdbcDriverFile string) (string, error) {
	absDriverPath, err := filepath.Abs(jdbcDriverFile)
	if err != nil {
		return "", fmt.Errorf("invalid jdbc_driver_file path: %w", err)
	}

	absLiquibaseDir, err := filepath.Abs(liquibaseDir)
	if err != nil {
		return "", fmt.Errorf("invalid liquibase directory path: %w", err)
	}

	relPath, err := filepath.Rel(absLiquibaseDir, absDriverPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve jdbc_driver_file path for liquibase: %w", err)
	}

	return filepath.ToSlash(relPath), nil
}

func validateFiles(jdbcDriverFile string) error {
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

	if _, err := os.Stat(jdbcDriverFile); err != nil {
		return fmt.Errorf("jdbc driver file not found: %s", jdbcDriverFile)
	}

	if _, err := exec.LookPath("liquibase"); err != nil {
		return fmt.Errorf("liquibase command not found in PATH: %w", err)
	}

	return nil
}

func writeProperties(host, port, dbname, user, password, driverClasspath string) error {
	content := fmt.Sprintf(`changeLogFile: %s
url: jdbc:postgresql://%s:%s/%s
username: %s
password: %s
driver: org.postgresql.Driver
classpath: %s
`, changeLogFile, host, port, dbname, user, password, driverClasspath)

	propertiesPath := filepath.Join(liquibaseDir, propertiesFile)
	return os.WriteFile(propertiesPath, []byte(content), 0644)
}
