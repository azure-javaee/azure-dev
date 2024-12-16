package appdetect

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMavenProjectInEffectivePom(t *testing.T) {
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
			tempDir, err := os.MkdirTemp("", "test")
			if err != nil {
				t.Fatalf("Failed to create temp directory: %v", err)
			}
			defer os.RemoveAll(tempDir)

			pomPath := filepath.Join(tempDir, "pom.xml")
			err = os.WriteFile(pomPath, []byte(tt.pomContent), 0644)
			if err != nil {
				t.Fatalf("Failed to write temp POM file: %v", err)
			}

			project, err := getMavenProjectOfEffectivePom(pomPath)
			if err != nil {
				t.Fatalf("getMavenProjectOfEffectivePom failed: %v", err)
			}

			if len(project.Dependencies) != len(tt.expected) {
				t.Fatalf("Expected %d dependencies, got %d", len(tt.expected), len(project.Dependencies))
			}

			for i, dep := range project.Dependencies {
				if dep != tt.expected[i] {
					t.Errorf("Expected dependency %v, got %v", tt.expected[i], dep)
				}
			}
		})
	}
}

func TestGetMvnwCommandInProject(t *testing.T) {
	cases := []struct {
		pomPath     string
		expected    string
		description string
	}{
		{"project1/pom.xml", "project1/mvnw", "Wrapper in same directory"},
		{"project2/sub-dir/pom.xml", "project2/mvnw", "Wrapper in parent directory"},
		{"project3/sub-dir/sub-sub-dir/pom.xml", "project3/mvnw", "Wrapper in grandparent directory"},
		{"project4/pom.xml", "", "No wrapper found"},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "testdata")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tempDir)

			pomPath := filepath.Join(tempDir, c.pomPath)
			err = os.MkdirAll(filepath.Dir(pomPath), os.ModePerm)
			if err != nil {
				t.Errorf("failed to mkdir")
			}
			err = os.WriteFile(pomPath, []byte("<project></project>"), os.ModePerm)
			if err != nil {
				t.Errorf("failed to write file")
			}
			if c.expected != "" {
				expectedPath := filepath.Join(tempDir, c.expected)
				err = os.WriteFile(expectedPath, []byte("#!/bin/sh"), os.ModePerm)
				if err != nil {
					t.Errorf("failed to write file")
				}
			}

			result, _ := getMvnwCommandInProject(pomPath)
			expectedResult := ""
			if c.expected != "" {
				expectedResult = filepath.Join(tempDir, c.expected)
			}
			if result != expectedResult {
				t.Errorf("getMvnw(%q) == %q, expected %q", pomPath, result, expectedResult)
			}
		})
	}
}

func TestGetDownloadedMvnCommand(t *testing.T) {
	maven, err := getDownloadedMvnCommand()
	if err != nil {
		t.Errorf("getDownloadedMvnCommand failed, %v", err)
	}
	if maven == "" {
		t.Errorf("getDownloadedMvnCommand failed")
	}
}
