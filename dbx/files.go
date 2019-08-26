package dbx

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	dbxFiles "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
)

type Folder struct {
	Name string
	Path string
}

type File struct {
	Name string
	Path string
}

func List(token, path string) ([]Folder, []File, error) {

	dropboxConf := dropbox.Config{
		Token:    token,
		LogLevel: dropbox.LogInfo,
	}
	client := files.New(dropboxConf)

	result, err := client.ListFolder(&dbxFiles.ListFolderArg{
		Path: path,
	})
	if err != nil {
		return nil, nil, err
	}
	var folders []Folder
	var files []File

	for _, entry := range result.Entries {
		switch meta := entry.(type) {
		case *dbxFiles.FolderMetadata:
			folders = append(folders, Folder{
				Name: meta.Name,
				Path: meta.PathLower,
			})
		case *dbxFiles.FileMetadata:
			files = append(files, File{
				Name: meta.Name,
				Path: meta.PathLower,
			})
		}
	}
	return folders, files, nil
}
