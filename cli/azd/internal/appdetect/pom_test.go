package appdetect

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestToEffectivePom(t *testing.T) {
	tests := []struct {
		name      string
		pomString string
		expected  []dependency
	}{
		{
			name: "Test with two dependencies",
			pomString: `
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
			pomString: `
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
			pomString: `
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
			pomString: `
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
			err = os.WriteFile(pomPath, []byte(tt.pomString), 0600)
			if err != nil {
				t.Fatalf("Failed to write temp POM file: %v", err)
			}

			effectivePom, err := toEffectivePomByMvnCommand(pomPath)
			if err != nil {
				t.Fatalf("toEffectivePomByMvnCommand failed: %v", err)
			}

			if len(effectivePom.Dependencies) != len(tt.expected) {
				t.Fatalf("Expected: %d\nActual: %d", len(tt.expected), len(effectivePom.Dependencies))
			}

			for i, dep := range effectivePom.Dependencies {
				if dep != tt.expected[i] {
					t.Errorf("\nExpected: %s\nActual:   %s", tt.expected[i], dep)
				}
			}
		})
	}
}

func TestCreatePropertyMap(t *testing.T) {
	tests := []struct {
		name      string
		pomString string
		expected  map[string]string
	}{
		{
			name: "Test createPropertyMap",
			pomString: `
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
			pom, err := unmarshalPomFromString(tt.pomString)
			if err != nil {
				t.Fatalf("Failed to unmarshal string: %v", err)
			}
			createPropertyMap(&pom)
			if !reflect.DeepEqual(pom.propertyMap, tt.expected) {
				t.Fatalf("\nExpected: %s\nActual:   %s", tt.expected, pom.propertyMap)
			}
		})
	}
}

func TestReplacePropertyPlaceHolder(t *testing.T) {
	var tests = []struct {
		name     string
		inputPom pom
		expected pom
	}{
		{
			name: "Test replacePropertyPlaceHolder",
			inputPom: pom{
				GroupId:    "sampleGroupId",
				ArtifactId: "sampleArtifactId",
				Version:    "1.0.0",
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
					{
						GroupId:    "${project.groupId}",
						ArtifactId: "artifactIdTwo",
						Version:    "${project.version}",
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
				propertyMap: map[string]string{
					"version.spring.boot":        "3.3.5",
					"version.spring.cloud":       "2023.0.3",
					"version.spring.cloud.azure": "5.18.0",
					"another.property":           "${version.spring.cloud.azure}",
				},
				dependencyManagementMap: map[string]string{
					"groupIdOne:artifactIdOne:compile": "${version.spring.boot}",
				},
			},
			expected: pom{
				GroupId:    "sampleGroupId",
				ArtifactId: "sampleArtifactId",
				Version:    "1.0.0",
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
					{
						GroupId:    "sampleGroupId",
						ArtifactId: "artifactIdTwo",
						Version:    "1.0.0",
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
				propertyMap: map[string]string{
					"version.spring.boot":        "3.3.5",
					"version.spring.cloud":       "2023.0.3",
					"version.spring.cloud.azure": "5.18.0",
					"another.property":           "5.18.0",
					"project.groupId":            "sampleGroupId",
					"project.version":            "1.0.0",
				},
				dependencyManagementMap: map[string]string{
					"groupIdOne:artifactIdOne:compile": "3.3.5",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addCommonPropertiesLikeProjectGroupIdAndProjectVersionToPropertyMap(&tt.inputPom)
			replacePropertyPlaceHolderInPropertyMap(&tt.inputPom)
			replacePropertyPlaceHolderInGroupId(&tt.inputPom)
			createDependencyManagementMap(&tt.inputPom)
			replacePropertyPlaceHolderInVersion(&tt.inputPom)
			if !reflect.DeepEqual(tt.inputPom, tt.expected) {
				t.Fatalf("\nExpected: %s\nActual:   %s", tt.expected, tt.inputPom)
			}
		})
	}
}

func TestCreateDependencyManagementMap(t *testing.T) {
	var tests = []struct {
		name     string
		inputPom pom
		expected pom
	}{
		{
			name: "Test createDependencyManagementMap",
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
			createDependencyManagementMap(&tt.inputPom)
			if !reflect.DeepEqual(tt.inputPom, tt.expected) {
				t.Fatalf("\nExpected: %s\nActual:   %s", tt.expected, tt.inputPom)
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
			name: "Test updateDependencyVersionAccordingToDependencyManagement",
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
			updateDependencyVersionAccordingToDependencyManagement(&tt.inputPom)
			if !reflect.DeepEqual(tt.inputPom, tt.expected) {
				t.Fatalf("\nExpected: %s\nActual:   %s", tt.expected, tt.inputPom)
			}
		})
	}
}

func TestUpdateVersionAccordingToPropertiesAndDependencyManagement(t *testing.T) {
	var tests = []struct {
		name      string
		pomString string
		expected  []dependency
	}{
		{
			name: "Test updateVersionAccordingToPropertiesAndDependencyManagement",
			pomString: `
				<project>
					<modelVersion>4.0.0</modelVersion>
					<groupId>com.example</groupId>
					<artifactId>example-project</artifactId>
					<version>1.0.0</version>
					<properties>
						<version.slf4j>1.0.0</version.slf4j>
						<version.junit>2.0.0</version.junit>
					</properties>
					<dependencyManagement>
						<dependencies>
							<dependency>
								<groupId>org.slf4j</groupId>
								<artifactId>slf4j-api</artifactId>
								<version>${version.slf4j}</version>
							</dependency>
						</dependencies>
					</dependencyManagement>
					<dependencies>
						<dependency>
							<groupId>org.slf4j</groupId>
							<artifactId>slf4j-api</artifactId>
						</dependency>
						<dependency>
							<groupId>junit</groupId>
							<artifactId>junit</artifactId>
							<version>${version.junit}</version>
							<scope>test</scope>
						</dependency>
					</dependencies>
				</project>
				`,
			expected: []dependency{
				{
					GroupId:    "org.slf4j",
					ArtifactId: "slf4j-api",
					Version:    "1.0.0",
				},
				{
					GroupId:    "junit",
					ArtifactId: "junit",
					Version:    "2.0.0",
					Scope:      "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pom, err := unmarshalPomFromString(tt.pomString)
			if err != nil {
				t.Fatalf("Failed to unmarshal POM string: %v", err)
			}

			updateVersionAccordingToPropertiesAndDependencyManagement(&pom)
			if !reflect.DeepEqual(pom.Dependencies, tt.expected) {
				t.Fatalf("\nExpected: %s\nActual:   %s", tt.expected, pom.Dependencies)
			}
		})
	}
}

func TestGetMavenRepositoryUrl(t *testing.T) {
	var tests = []struct {
		name       string
		groupId    string
		artifactId string
		version    string
		expected   string
	}{
		{
			name:       "spring-boot-starter-parent",
			groupId:    "org.springframework.boot",
			artifactId: "spring-boot-starter-parent",
			version:    "3.4.0",
			expected: "https://repo.maven.apache.org/maven2/org/springframework/boot/spring-boot-starter-parent/3.4.0/" +
				"spring-boot-starter-parent-3.4.0.pom",
		},
		{
			name:       "spring-boot-dependencies",
			groupId:    "org.springframework.boot",
			artifactId: "spring-boot-dependencies",
			version:    "3.4.0",
			expected: "https://repo.maven.apache.org/maven2/org/springframework/boot/spring-boot-dependencies/3.4.0/" +
				"spring-boot-dependencies-3.4.0.pom",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := getRemoteMavenRepositoryUrl(tt.groupId, tt.artifactId, tt.version)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Fatalf("\nExpected: %s\nActual:   %s", tt.expected, actual)
			}
		})
	}
}

func TestGetSimulatedEffectivePom(t *testing.T) {
	var tests = []struct {
		name       string
		groupId    string
		artifactId string
		version    string
		expected   int
	}{
		{
			name:       "spring-boot-starter-parent",
			groupId:    "org.springframework.boot",
			artifactId: "spring-boot-starter-parent",
			version:    "3.4.0",
			expected:   1496,
		},
		{
			name:       "spring-boot-dependencies",
			groupId:    "org.springframework.boot",
			artifactId: "spring-boot-dependencies",
			version:    "3.4.0",
			expected:   1496,
		},
		{
			name:       "kotlin-bom",
			groupId:    "org.jetbrains.kotlin",
			artifactId: "kotlin-bom",
			version:    "1.9.25",
			expected:   23,
		},
		{
			name:       "infinispan-bom",
			groupId:    "org.infinispan",
			artifactId: "infinispan-bom",
			version:    "15.0.11.Final",
			expected:   65,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pom, err := getSimulatedEffectivePomFromRemoteMavenRepository(tt.groupId, tt.artifactId, tt.version)
			if err != nil {
				t.Fatalf("Failed to create temp directory: %v", err)
			}
			for _, value := range pom.dependencyManagementMap {
				if isVariable(value) {
					t.Fatalf("Unresolved property: value = %s", value)
				}
			}
			actual := len(pom.dependencyManagementMap)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Fatalf("\nExpected: %d\nActual:   %d", tt.expected, actual)
			}
		})
	}
}
