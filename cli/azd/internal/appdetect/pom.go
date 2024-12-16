package appdetect

import (
	"archive/zip"
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func getMavenProjectOfEffectivePom(pomPath string) (mavenProject, error) {
	if !commandExistsInPath("java") {
		log.Println("java not exist, skip get dependencies by 'mvn help:effective-pom'.")
	}
	mvn, err := getMvnCommand(pomPath)
	if err != nil {
		return mavenProject{}, err
	}
	cmd := exec.Command(mvn, "help:effective-pom", "-f", pomPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mavenProject{}, err
	}
	effectivePom, err := getEffectivePomFromConsoleOutput(string(output))
	if err != nil {
		return mavenProject{}, err
	}
	var project mavenProject
	if err := xml.Unmarshal([]byte(effectivePom), &project); err != nil {
		return mavenProject{}, fmt.Errorf("parsing xml: %w", err)
	}
	return project, nil
}

func commandExistsInPath(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func getMvnCommand(pomPath string) (string, error) {
	if commandExistsInPath("mvn") {
		return "mvn", nil
	}
	mvnwCommand, err := getMvnwCommandInProject(pomPath)
	if err == nil {
		return mvnwCommand, nil
	}
	return getDownloadedMvnCommand()
}

func getMvnwCommandInProject(pomPath string) (string, error) {
	mvnwCommand := "mvnw"
	dir := filepath.Dir(pomPath)
	for {
		commandPath := filepath.Join(dir, mvnwCommand)
		if fileExists(commandPath) {
			return commandPath, nil
		}
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}
		dir = parentDir
	}
	return "", fmt.Errorf("failed to find mvnw command in project")
}

func fileExists(path string) bool {
	if path == "" {
		return false
	}
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
}

func getDownloadedMvnCommand() (string, error) {
	mavenVersion := "3.9.9"
	mavenCommand, err := getAzdMvnCommand(mavenVersion)
	if err != nil {
		return "", err
	}
	if fileExists(mavenCommand) {
		log.Println("Skip downloading maven because it already exists.")
		return mavenCommand, nil
	}
	log.Println("Downloading maven")
	mavenDir, err := getAzdMvnDir()
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(mavenDir); os.IsNotExist(err) {
		err = os.Mkdir(mavenDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("unable to create directory: %w", err)
		}
	}
	mavenURL := fmt.Sprintf("https://repo.maven.apache.org/maven2/org/apache/maven/apache-maven/"+
		"%s/apache-maven-%s-bin.zip", mavenVersion, mavenVersion)
	mavenFile := fmt.Sprintf("maven-wrapper-%s-bin.zip", mavenVersion)
	wrapperPath := filepath.Join(mavenDir, mavenFile)
	err = downloadFile(wrapperPath, mavenURL)
	if err != nil {
		return "", err
	}
	err = unzip(wrapperPath, mavenDir)
	if err != nil {
		return "", fmt.Errorf("failed to unzip maven bin.zip: %w", err)
	}
	return mavenCommand, nil
}

func getAzdMvnDir() (string, error) {
	azdMvnFolderName := "azd-maven"
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get user home directory: %w", err)
	}
	return filepath.Join(userHome, azdMvnFolderName), nil
}

func getAzdMvnCommand(mavenVersion string) (string, error) {
	mavenDir, err := getAzdMvnDir()
	if err != nil {
		return "", err
	}
	azdMvnCommand := filepath.Join(mavenDir, "apache-maven-"+mavenVersion, "bin", "mvn")
	return azdMvnCommand, nil
}

func downloadFile(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func unzip(src string, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, f := range reader.File {
		filePath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			err := os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				return err
			}

			outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			_, err = io.Copy(outFile, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getEffectivePomFromConsoleOutput(consoleOutput string) (string, error) {
	var effectivePom strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(consoleOutput))
	inProject := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "<project") {
			inProject = true
		} else if strings.HasPrefix(strings.TrimSpace(line), "</project>") {
			effectivePom.WriteString(line)
			break
		}
		if inProject {
			effectivePom.WriteString(line)
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to scan console output. %w", err)
	}
	return effectivePom.String(), nil
}
