package appdetect

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestToEffectivePom(t *testing.T) {
	tests := []struct {
		name       string
		pomContent string
		expected   []dependency
	}{
		{
			name: "Test with two dependencies",
			pomContent: `
				<project>
					<modelVersion>4.0.0</modelVersion>
					<groupId>com.example</groupId>
					<artifactId>example-project</artifactId>
					<version>1.0.0</version>
					<dependencies>
						<dependency>
							<groupId>org.springframework</groupId>
							<artifactId>spring-core</artifactId>
							<version>5.3.8</version>
							<scope>compile</scope>
						</dependency>
						<dependency>
							<groupId>junit</groupId>
							<artifactId>junit</artifactId>
							<version>4.13.2</version>
							<scope>test</scope>
						</dependency>
					</dependencies>
				</project>
				`,
			expected: []dependency{
				{
					GroupId:    "org.springframework",
					ArtifactId: "spring-core",
					Version:    "5.3.8",
					Scope:      "compile",
				},
				{
					GroupId:    "junit",
					ArtifactId: "junit",
					Version:    "4.13.2",
					Scope:      "test",
				},
			},
		},
		{
			name: "Test with no dependencies",
			pomContent: `
				<project>
					<modelVersion>4.0.0</modelVersion>
					<groupId>com.example</groupId>
					<artifactId>example-project</artifactId>
					<version>1.0.0</version>
					<dependencies>
					</dependencies>
				</project>
				`,
			expected: []dependency{},
		},
		{
			name: "Test with one dependency which version is decided by dependencyManagement",
			pomContent: `
				<project>
					<modelVersion>4.0.0</modelVersion>
					<groupId>com.example</groupId>
					<artifactId>example-project</artifactId>
					<version>1.0.0</version>
					<dependencies>
						<dependency>
							<groupId>org.slf4j</groupId>
							<artifactId>slf4j-api</artifactId>
						</dependency>
					</dependencies>
					<dependencyManagement>
						<dependencies>
							<dependency>
								<groupId>org.springframework.boot</groupId>
								<artifactId>spring-boot-dependencies</artifactId>
								<version>3.0.0</version>
								<type>pom</type>
								<scope>import</scope>
							</dependency>
						</dependencies>
					</dependencyManagement>
				</project>
				`,
			expected: []dependency{
				{
					GroupId:    "org.slf4j",
					ArtifactId: "slf4j-api",
					Version:    "2.0.4",
					Scope:      "compile",
				},
			},
		},
		{
			name: "Test with one dependency which version is decided by parent",
			pomContent: `
				<project>
					<parent>
						<groupId>org.springframework.boot</groupId>
						<artifactId>spring-boot-starter-parent</artifactId>
						<version>3.0.0</version>
						<relativePath/> <!-- lookup parent from repository -->
					</parent>
					<modelVersion>4.0.0</modelVersion>
					<groupId>com.example</groupId>
					<artifactId>example-project</artifactId>
					<version>1.0.0</version>
					<dependencies>
						<dependency>
							<groupId>org.slf4j</groupId>
							<artifactId>slf4j-api</artifactId>
						</dependency>
					</dependencies>
				</project>
				`,
			expected: []dependency{
				{
					GroupId:    "org.slf4j",
					ArtifactId: "slf4j-api",
					Version:    "2.0.4",
					Scope:      "compile",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "TestToEffectivePom")
			if err != nil {
				t.Fatalf("Failed to create temp directory: %v", err)
			}
			defer func(path string) {
				err := os.RemoveAll(path)
				if err != nil {
					t.Fatalf("Failed to remove all in directory: %v", err)
				}
			}(tempDir)

			pomPath := filepath.Join(tempDir, "pom.xml")
			err = os.WriteFile(pomPath, []byte(tt.pomContent), 0600)
			if err != nil {
				t.Fatalf("Failed to write temp POM file: %v", err)
			}

			effectivePom, err := toEffectivePomByMvnCommand(pomPath)
			if err != nil {
				t.Fatalf("toEffectivePomByMvnCommand failed: %v", err)
			}

			if len(effectivePom.Dependencies) != len(tt.expected) {
				t.Fatalf("Expected %d dependencies, got %d", len(tt.expected), len(effectivePom.Dependencies))
			}

			for i, dep := range effectivePom.Dependencies {
				if dep != tt.expected[i] {
					t.Errorf("Expected dependency %v, got %v", tt.expected[i], dep)
				}
			}
		})
	}
}

func TestReadPropertiesToPropertyMap(t *testing.T) {
	tests := []struct {
		name       string
		pomContent string
		expected   map[string]string
	}{
		{
			name: "Test with two dependencies",
			pomContent: `
				<project>
					<modelVersion>4.0.0</modelVersion>
					<groupId>com.example</groupId>
					<artifactId>example-project</artifactId>
					<version>1.0.0</version>
					<properties>
						<version.spring.boot>3.3.5</version.spring.boot>
						<version.spring.cloud>2023.0.3</version.spring.cloud>
						<version.spring.cloud.azure>5.18.0</version.spring.cloud.azure>
					</properties>
				</project>
				`,
			expected: map[string]string{
				"version.spring.boot":        "3.3.5",
				"version.spring.cloud":       "2023.0.3",
				"version.spring.cloud.azure": "5.18.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pom, err := unmarshalPomFromString(tt.pomContent)
			if err != nil {
				t.Fatalf("Failed to unmarshal content: %v", err)
			}
			readPropertiesToPropertyMap(&pom)
			if !reflect.DeepEqual(pom.propertyMap, tt.expected) {
				t.Fatalf("Expected %s dependencies, got %s", tt.expected, pom.propertyMap)
			}
		})
	}
}

func TestFulfillPropertyValue(t *testing.T) {
	var tests = []struct {
		name     string
		inputPom pom
		expected pom
	}{
		{
			name: "Test updateVersionAccordingToPropertyMap",
			inputPom: pom{
				propertyMap: map[string]string{
					"version.spring.boot":        "3.3.5",
					"version.spring.cloud":       "2023.0.3",
					"version.spring.cloud.azure": "5.18.0",
				},
				DependencyManagement: dependencyManagement{
					Dependencies: []dependency{
						{
							GroupId:    "groupIdOne",
							ArtifactId: "artifactIdOne",
							Version:    "${version.spring.boot}",
						},
					},
				},
				Dependencies: []dependency{
					{
						GroupId:    "groupIdTwo",
						ArtifactId: "artifactIdTwo",
						Version:    "${version.spring.cloud}",
					},
				},
				Build: build{
					Plugins: []plugin{
						{
							GroupId:    "groupIdThree",
							ArtifactId: "artifactIdThree",
							Version:    "${version.spring.cloud.azure}",
						},
					},
				},
			},
			expected: pom{
				propertyMap: map[string]string{
					"version.spring.boot":        "3.3.5",
					"version.spring.cloud":       "2023.0.3",
					"version.spring.cloud.azure": "5.18.0",
				},
				DependencyManagement: dependencyManagement{
					Dependencies: []dependency{
						{
							GroupId:    "groupIdOne",
							ArtifactId: "artifactIdOne",
							Version:    "3.3.5",
						},
					},
				},
				Dependencies: []dependency{
					{
						GroupId:    "groupIdTwo",
						ArtifactId: "artifactIdTwo",
						Version:    "2023.0.3",
					},
				},
				Build: build{
					Plugins: []plugin{
						{
							GroupId:    "groupIdThree",
							ArtifactId: "artifactIdThree",
							Version:    "5.18.0",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateVersionAccordingToPropertyMap(&tt.inputPom)
			if !reflect.DeepEqual(tt.inputPom, tt.expected) {
				t.Fatalf("Expected %s dependencies, got %s", tt.expected, tt.inputPom)
			}
		})
	}
}

func TestReadDependencyManagementToDependencyManagementMap(t *testing.T) {
	var tests = []struct {
		name     string
		inputPom pom
		expected pom
	}{
		{
			name: "Test updateVersionAccordingToPropertyMap",
			inputPom: pom{
				DependencyManagement: dependencyManagement{
					Dependencies: []dependency{
						{
							GroupId:    "groupIdOne",
							ArtifactId: "artifactIdOne",
							Version:    "1.0.0",
						},
					},
				},
				Dependencies: []dependency{
					{
						GroupId:    "groupIdOne",
						ArtifactId: "artifactIdOne",
					},
				},
			},
			expected: pom{
				DependencyManagement: dependencyManagement{
					Dependencies: []dependency{
						{
							GroupId:    "groupIdOne",
							ArtifactId: "artifactIdOne",
							Version:    "1.0.0",
						},
					},
				},
				Dependencies: []dependency{
					{
						GroupId:    "groupIdOne",
						ArtifactId: "artifactIdOne",
					},
				},
				dependencyManagementMap: map[string]string{
					"groupIdOne:artifactIdOne": "1.0.0",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readDependencyManagementToDependencyManagementMap(&tt.inputPom)
			if !reflect.DeepEqual(tt.inputPom, tt.expected) {
				t.Fatalf("Expected %s dependencies, got %s", tt.expected, tt.inputPom)
			}
		})
	}
}

func TestUpdateVersionAccordingToDependencyManagementMap(t *testing.T) {
	var tests = []struct {
		name     string
		inputPom pom
		expected pom
	}{
		{
			name: "Test updateVersionAccordingToPropertyMap",
			inputPom: pom{
				Dependencies: []dependency{
					{
						GroupId:    "groupIdOne",
						ArtifactId: "artifactIdOne",
					},
				},
				dependencyManagementMap: map[string]string{
					"groupIdOne:artifactIdOne": "1.0.0",
				},
			},
			expected: pom{
				dependencyManagementMap: map[string]string{
					"groupIdOne:artifactIdOne": "1.0.0",
				},
				Dependencies: []dependency{
					{
						GroupId:    "groupIdOne",
						ArtifactId: "artifactIdOne",
						Version:    "1.0.0",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateVersionAccordingToDependencyManagementMap(&tt.inputPom)
			if !reflect.DeepEqual(tt.inputPom, tt.expected) {
				t.Fatalf("Expected %s dependencies, got %s", tt.expected, tt.inputPom)
			}
		})
	}
}
