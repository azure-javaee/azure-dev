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
	XmlName                 xml.Name             `xml:"project"`
	Parent                  parent               `xml:"parent"`
	Modules                 []string             `xml:"modules>module"` // Capture the modules
	Properties              Properties           `xml:"properties"`
	Dependencies            []dependency         `xml:"dependencies>dependency"`
	DependencyManagement    dependencyManagement `xml:"dependencyManagement"`
	Build                   build                `xml:"build"`
	path                    string               // todo: add 'pom.xml' in the path.
	propertyMap             map[string]string
	dependencyManagementMap map[string]string
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

// Not strictly equal to effective pom. Just try best to make sure the Dependencies are accurate.
func createSimulatedEffectivePomFromFilePath(filePath string) (*pom, error) {
	pom, err := unmarshalPomFromFilePath(filePath)
	if err != nil {
		return nil, err
	}
	pom.path = filepath.Dir(filePath)
	convertToSimulatedEffectivePom(&pom)
	return &pom, nil
}

func convertToSimulatedEffectivePom(pom *pom) {
	updateVersionAccordingToPropertiesAndDependencyManagement(pom)
	absorbInformationFromParentAndImportedDependenciesInDependencyManagement(pom)
}

func updateVersionAccordingToPropertiesAndDependencyManagement(pom *pom) {
	readPropertiesToPropertyMap(pom)
	updateVersionAccordingToPropertyMap(pom)
	// not handle pluginManagement because now we only care about dependency's version.
	readDependencyManagementToDependencyManagementMap(pom)
	updateVersionAccordingToDependencyManagementMap(pom)
}

func absorbInformationFromParentAndImportedDependenciesInDependencyManagement(pom *pom) {
	absorbInformationFromParent(pom)
	absorbInformationFromImportedDependenciesInDependencyManagement(pom)
}

func absorbInformationFromParent(pom *pom) {
	// todo finish this
}

func absorbInformationFromImportedDependenciesInDependencyManagement(pom *pom) {
	for _, dep := range pom.DependencyManagement.Dependencies {
		if dep.Scope != "import" {
			continue
		}
		importedPom, err := getSimulatedEffectivePomByDependency(dep)
		if err != nil {
			continue // ignore error, because we want to get as more information as possible
		}
		for key, value := range importedPom.propertyMap {
			addToPropertyMapIfKeyIsNew(pom, key, value)
		}
		for _, dep := range importedPom.DependencyManagement.Dependencies {
			addToDependencyManagementMapIfDependencyIsNew(pom, dep)
		}
	}
	updateVersionAccordingToPropertyMap(pom)
	updateVersionAccordingToDependencyManagementMap(pom)
}

func getSimulatedEffectivePomByDependency(dependency dependency) (pom, error) {
	// todo finish this
	return pom{}, nil
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
	for _, entry := range pom.Properties.Entries {
		addToPropertyMapIfKeyIsNew(pom, entry.XMLName.Local, entry.Value)
	}
}

func addToPropertyMapIfKeyIsNew(pom *pom, key string, value string) {
	if pom.propertyMap == nil {
		pom.propertyMap = make(map[string]string)
	}
	if _, ok := pom.propertyMap[key]; ok {
		return
	}
	pom.propertyMap[key] = value
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

func toDependencyManagementMapKey(dependency dependency) string {
	return fmt.Sprintf("%s:%s", dependency.GroupId, dependency.ArtifactId)
}

func readDependencyManagementToDependencyManagementMap(pom *pom) {
	if pom.dependencyManagementMap == nil {
		pom.dependencyManagementMap = make(map[string]string)
	}
	for _, dep := range pom.DependencyManagement.Dependencies {
		addToDependencyManagementMapIfDependencyIsNew(pom, dep)
	}
}

func addToDependencyManagementMapIfDependencyIsNew(pom *pom, dependency dependency) {
	version := strings.TrimSpace(dependency.Version)
	if version == "" {
		return
	}
	key := toDependencyManagementMapKey(dependency)
	if _, ok := pom.dependencyManagementMap[key]; ok {
		return
	}
	pom.dependencyManagementMap[key] = version
}

func updateVersionAccordingToDependencyManagementMap(pom *pom) {
	if pom.dependencyManagementMap == nil {
		pom.dependencyManagementMap = make(map[string]string)
	}
	for i, dep := range pom.Dependencies {
		if strings.TrimSpace(dep.Version) != "" {
			continue
		}
		key := toDependencyManagementMapKey(dep)
		managedVersion := pom.dependencyManagementMap[key]
		if managedVersion != "" {
			pom.Dependencies[i].Version = managedVersion
		}
	}
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
