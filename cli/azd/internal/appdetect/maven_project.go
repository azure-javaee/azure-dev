package appdetect

type mavenProject struct {
	pom pom
}

func toMavenProject(pomFilePath string) (*mavenProject, error) {
	pom, err := createSimulatedEffectivePom(pomFilePath)
	if err != nil {
		return nil, err
	}
	return &mavenProject{
		pom: *pom,
	}, nil
}
