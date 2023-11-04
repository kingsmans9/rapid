package handlers

import (
	"fmt"
	"net/http"
)

func (s *site) GetData(w http.ResponseWriter, r *http.Request) {
	sites, err := s.store.ListSites(r.Context())
	if err != nil {
		JSONWithError(w, http.StatusInternalServerError, "error listing sites")
		return
	}

	if len(sites) == 0 {
		JSON(w, http.StatusOK, []string{})
		return
	}

	dbData, err := s.store.DumpData(r.Context())
	if err != nil {
		JSONWithError(w, http.StatusInternalServerError, "error dumping data")
		return
	}

	tmpDir, err := s.fs.TempDir("", "data")
	if err != nil {
		JSONWithError(w, http.StatusInternalServerError, fmt.Errorf("error creating temp dir: %v", err).Error())
		return
	}
	defer s.fs.RemoveAll(tmpDir)

	tarFile := tmpDir + "/data.tar.gz"
	err = s.rapidDataCopier.CopyDataAndArchive(r.Context(), sites, dbData, tarFile)
	if err != nil {
		JSONWithError(w, http.StatusInternalServerError, fmt.Errorf("error copying data: %v", err).Error())
		return
	}

	http.ServeFile(w, r, tarFile)
}
