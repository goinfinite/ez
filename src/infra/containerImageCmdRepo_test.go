package infra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/dto"
)

func TestContainerImageCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	containerImageCmdRepo := NewContainerImageCmdRepo(persistentDbSvc)
	containerImageQueryRepo := NewContainerImageQueryRepo(persistentDbSvc)

	t.Run("CreateSnapshot", func(t *testing.T) {
		containerQueryRepo := NewContainerQueryRepo(persistentDbSvc)
		containersList, err := containerQueryRepo.Read()
		if err != nil {
			t.Fatal(err)
		}
		if len(containersList) == 0 {
			t.Fatal("NoContainersFound")
		}

		createDto := dto.CreateContainerSnapshotImage{
			AccountId:   containersList[0].AccountId,
			ContainerId: containersList[0].Id,
		}
		_, err = containerImageCmdRepo.CreateSnapshot(createDto)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("DeleteImage", func(t *testing.T) {
		imagesList, err := containerImageQueryRepo.Read()
		if err != nil {
			t.Fatal(err)
		}

		deleteDto := dto.DeleteContainerImage{
			AccountId: imagesList[0].AccountId,
			ImageId:   imagesList[0].Id,
		}
		err = containerImageCmdRepo.Delete(deleteDto)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("CreateArchiveFile", func(t *testing.T) {
		imagesList, err := containerImageQueryRepo.Read()
		if err != nil {
			t.Fatal(err)
		}
		if len(imagesList) == 0 {
			t.Fatal("NoImagesFound")
		}

		createDto := dto.CreateContainerImageArchiveFile{
			AccountId: imagesList[0].AccountId,
			ImageId:   imagesList[0].Id,
		}
		_, err = containerImageCmdRepo.CreateArchiveFile(createDto)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("DeleteArchiveFile", func(t *testing.T) {
		archiveFilesList, err := containerImageQueryRepo.ReadArchiveFiles()
		if err != nil {
			t.Fatal(err)
		}
		if len(archiveFilesList) == 0 {
			t.Fatal("NoArchiveFilesFound")
		}

		err = containerImageCmdRepo.DeleteArchiveFile(archiveFilesList[0])
		if err != nil {
			t.Fatal(err)
		}
	})
}
