package appdetect

import (
	"context"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/azure/azure-dev/cli/azd/internal/tracing"
	"github.com/azure/azure-dev/cli/azd/internal/tracing/fields"
)

type javaDetector struct {
	rootProjects      []pom
	mavenWrapperPaths []mavenWrapper
}

type mavenWrapper struct {
	posixPath string
	winPath   string
}

// JavaProjectOptionMavenParentPath The parent module path of the maven multi-module project
const JavaProjectOptionMavenParentPath = "parentPath"

// JavaProjectOptionPosixMavenWrapperPath The path to the maven wrapper script for POSIX systems
const JavaProjectOptionPosixMavenWrapperPath = "posixMavenWrapperPath"

// JavaProjectOptionWinMavenWrapperPath The path to the maven wrapper script for Windows systems
const JavaProjectOptionWinMavenWrapperPath = "winMavenWrapperPath"

func (jd *javaDetector) Language() Language {
	return Java
}

func (jd *javaDetector) DetectProject(ctx context.Context, path string, entries []fs.DirEntry) (*Project, error) {
	for _, entry := range entries {
		if strings.ToLower(entry.Name()) == "pom.xml" {
			tracing.SetUsageAttributes(fields.AppInitJavaDetect.String("start"))
			pomFile := filepath.Join(path, entry.Name())
			project, err := toPom(pomFile)
			if err != nil {
				log.Printf("Please edit azure.yaml manually to satisfy your requirement. azd can not help you "+
					"to that by detect your java project because error happened when reading pom.xml: %s. ", err)
				return nil, nil
			}

			if len(project.Modules) > 0 {
				// This is a multi-module project, we will capture the analysis, but return nil
				// to continue recursing
				jd.rootProjects = append(jd.rootProjects, *project)
				jd.mavenWrapperPaths = append(jd.mavenWrapperPaths, mavenWrapper{
					posixPath: detectMavenWrapper(path, "mvnw"),
					winPath:   detectMavenWrapper(path, "mvnw.cmd"),
				})
				return nil, nil
			}

			var currentRoot *pom
			var currentWrapper mavenWrapper
			for i, rootProject := range jd.rootProjects {
				// we can say that the project is in the root project if the path is under the project
				if inRoot := strings.HasPrefix(pomFile, rootProject.path); inRoot {
					currentRoot = &rootProject
					currentWrapper = jd.mavenWrapperPaths[i]
				}
			}

			result, err := detectDependencies(currentRoot, project, &Project{
				Language:      Java,
				Path:          path,
				DetectionRule: "Inferred by presence of: pom.xml",
			})
			if currentRoot != nil {
				result.Options = map[string]interface{}{
					JavaProjectOptionMavenParentPath:       currentRoot.path,
					JavaProjectOptionPosixMavenWrapperPath: currentWrapper.posixPath,
					JavaProjectOptionWinMavenWrapperPath:   currentWrapper.winPath,
				}
			}
			if err != nil {
				log.Printf("Please edit azure.yaml manually to satisfy your requirement. azd can not help you "+
					"to that by detect your java project because error happened when detecting dependencies: %s", err)
				return nil, nil
			}

			tracing.SetUsageAttributes(fields.AppInitJavaDetect.String("finish"))
			return result, nil
		}
	}

	return nil, nil
}

func detectDependencies(rootPom *pom, pom *pom, project *Project) (*Project, error) {
	detectAzureDependenciesByAnalyzingSpringBootProject(rootPom, pom, project)
	return project, nil
}

func detectMavenWrapper(path string, executable string) string {
	wrapperPath := filepath.Join(path, executable)
	if fileExists(wrapperPath) {
		return wrapperPath
	}
	return ""
}
