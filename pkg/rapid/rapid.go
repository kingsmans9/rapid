package rapid

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spectrocloud/rapid-agent/pkg/store/types"
	"github.com/spectrocloud/rapid-agent/pkg/util/filesystem"
	"github.com/spectrocloud/rapid-agent/pkg/util/logger"
)

type RapidDataCopier interface {
	CopyDataAndArchive(ctx context.Context, sites []*types.Site, dbData map[string]interface{}, tarFile string) error
}

type rapidDataCopier struct {
	rapidDataDir string
	fs           filesystem.Interface
}

func NewRapidDataCopier(rapidDataDir string, fs filesystem.Interface) RapidDataCopier {
	return &rapidDataCopier{
		rapidDataDir: rapidDataDir,
		fs:           fs,
	}
}

func (r *rapidDataCopier) CopyDataAndArchive(ctx context.Context, sites []*types.Site, dbData map[string]interface{}, tarFile string) error {
	err := r.zipRapidData(tarFile, sites, dbData)
	if err != nil {
		return fmt.Errorf("error creating zip file: %v", err)
	}

	return nil
}

func (r *rapidDataCopier) zipRapidData(destinationPath string, sites []*types.Site, dbData map[string]interface{}) error {
	pathToZip := r.rapidDataDir
	matchDirPath := make(map[string]bool)
	for _, site := range sites {
		matchDirPath[site.SiteName] = true
	}

	destinationFile, err := r.fs.Create(destinationPath)
	if err != nil {
		return err
	}

	zipWriter := zip.NewWriter(destinationFile)
	err = filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if _, ok := matchDirPath[info.Name()]; !ok && filePath != pathToZip {
				logger.Debugf("skipping dir %s", filePath)
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Dir(filePath) == pathToZip {
			logger.Debugf("skipping file %s", filePath)
			return nil
		}

		relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))
		zipFile, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, site := range sites {
		siteDBData := map[string]interface{}{}
		for key, value := range dbData {
			if strings.HasPrefix(key, site.SiteName+"_") {
				siteDBData[key] = value
				delete(dbData, key)
			}
		}

		siteDBDataBytes, err := json.Marshal(siteDBData)
		if err != nil {
			return err
		}

		filePath := pathToZip + "/" + site.SiteName + "/postgres_data.json"
		relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))
		zipFile, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		_, err = zipFile.Write(siteDBDataBytes)
		if err != nil {
			return err
		}

	}

	dbDataBytes, err := json.Marshal(dbData)
	if err != nil {
		return err
	}

	filePath := pathToZip + "/" + "postgres_data.json"
	relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))
	zipFile, err := zipWriter.Create(relPath)
	if err != nil {
		return err
	}

	_, err = zipFile.Write(dbDataBytes)
	if err != nil {
		return err
	}

	err = zipWriter.Close()
	if err != nil {
		return err
	}
	return nil
}
