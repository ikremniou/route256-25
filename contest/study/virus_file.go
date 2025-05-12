package study

import "strings"

type JsonFileSystem struct {
	Dir     string           `json:"dir"`
	Files   []string         `json:"files" optional:"true"`
	Folders []JsonFileSystem `json:"folders" optional:"true"`
}

func isVirusFile(fileName string) bool {
	return strings.HasSuffix(fileName, ".hack")
}

func countVirusesInFs(fs *JsonFileSystem, numOfViruses *int, parentHasVirus bool) bool {
	var hasVirus = parentHasVirus
	for i := 0; i < len(fs.Files) && !hasVirus; i++ {
		if isVirusFile(fs.Files[i]) {
			hasVirus = true
		}
	}

	for i := 0; i < len(fs.Folders); i++ {
		countVirusesInFs(&fs.Folders[i], numOfViruses, hasVirus)
	}

	if hasVirus {
		*numOfViruses = *numOfViruses + len(fs.Files)
	}

	return hasVirus
}

func FindNumberOfVirusFiles(fs *JsonFileSystem) int {
	var countViruses = 0
	countVirusesInFs(fs, &countViruses, false)

	return countViruses
}
