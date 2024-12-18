package appdetect

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// pom represents the top-level structure of a Maven POM file.
type pom struct {
	XmlName              xml.Name             `xml:"project"`
	Parent               parent               `xml:"parent"`
	Modules              []string             `xml:"modules>module"` // Capture the modules
	Properties           Properties           `xml:"properties"`
	Dependencies         []dependency         `xml:"dependencies>dependency"`
	DependencyManagement dependencyManagement `xml:"dependencyManagement"`
	Build                build                `xml:"build"`
	path                 string               // todo: add 'pom.xml' in the path.
	propertyMap          map[string]string
}

// Parent represents the parent POM if this project is a module.
type parent struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
}

type Properties struct {
	Entries []Property `xml:",any"` // Capture all elements inside <properties>
}

type Property struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// Dependency represents a single Maven dependency.
type dependency struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
	Scope      string `xml:"scope,omitempty"`
}

// DependencyManagement includes a list of dependencies that are managed.
type dependencyManagement struct {
	Dependencies []dependency `xml:"dependencies>dependency"`
}

// Build represents the build configuration which can contain plugins.
type build struct {
	Plugins []plugin `xml:"plugins>plugin"`
}

// Plugin represents a build plugin.
type plugin struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
}

func toPom(filePath string) (*pom, error) {
	pom, err := unmarshalPomFromFilePath(filePath)
	if err != nil {
		return nil, err
	}
	readPropertiesToPropertyMap(&pom)
	updateVersionAccordingToPropertyMap(&pom)
	pom.path = filepath.Dir(filePath)
	return &pom, nil
}

func unmarshalPomFromFilePath(pomFilePath string) (pom, error) {
	bytes, err := os.ReadFile(pomFilePath)
	if err != nil {
		return pom{}, err
	}
	return unmarshalPomFromBytes(bytes)
}

func unmarshalPomFromString(pomString string) (pom, error) {
	return unmarshalPomFromBytes([]byte(pomString))
}

func unmarshalPomFromBytes(pomBytes []byte) (pom, error) {
	var unmarshalledPom pom
	if err := xml.Unmarshal(pomBytes, &unmarshalledPom); err != nil {
		return pom{}, fmt.Errorf("parsing xml: %w", err)
	}
	return unmarshalledPom, nil
}

func readPropertiesToPropertyMap(pom *pom) {
	if pom.propertyMap == nil {
		pom.propertyMap = make(map[string]string)
	}
	for _, entry := range pom.Properties.Entries {
		pom.propertyMap[entry.XMLName.Local] = entry.Value
	}
}

func updateVersionAccordingToPropertyMap(pom *pom) {
	for i, dep := range pom.DependencyManagement.Dependencies {
		if isVariable(dep.Version) {
			variableName := getVariableName(dep.Version)
			if variableValue, ok := pom.propertyMap[variableName]; ok {
				pom.DependencyManagement.Dependencies[i].Version = variableValue
			}
		}
	}
	for i, dep := range pom.Dependencies {
		if isVariable(dep.Version) {
			variableName := getVariableName(dep.Version)
			if variableValue, ok := pom.propertyMap[variableName]; ok {
				pom.Dependencies[i].Version = variableValue
			}
		}
	}
	for i, dep := range pom.Build.Plugins {
		if isVariable(dep.Version) {
			variableName := getVariableName(dep.Version)
			if variableValue, ok := pom.propertyMap[variableName]; ok {
				pom.Build.Plugins[i].Version = variableValue
			}
		}
	}
}

const variablePrefix = "${"
const variableSuffix = "}"

func isVariable(value string) bool {
	return strings.HasPrefix(value, variablePrefix) && strings.HasSuffix(value, variableSuffix)
}

func getVariableName(value string) string {
	return strings.TrimSuffix(strings.TrimPrefix(value, variablePrefix), variableSuffix)
}

func toEffectivePomByMvnCommand(pomPath string) (pom, error) {
	if !commandExistsInPath("java") {
		return pom{}, fmt.Errorf("can not get effective pom because java command not exist")
	}
	mvn, err := getMvnCommandFromPath(pomPath)
	if err != nil {
		return pom{}, err
	}
	cmd := exec.Command(mvn, "help:effective-pom", "-f", pomPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return pom{}, err
	}
	effectivePom, err := getEffectivePomFromConsoleOutput(string(output))
	if err != nil {
		return pom{}, err
	}
	var project pom
	if err := xml.Unmarshal([]byte(effectivePom), &project); err != nil {
		return pom{}, fmt.Errorf("parsing xml: %w", err)
	}
	return project, nil
}

func commandExistsInPath(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
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
