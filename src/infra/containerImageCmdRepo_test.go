package infra

import (
	"testing"

	testHelpers "github.com/goinfinite/ez/src/devUtils"
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/useCase"
)

func TestContainerImageCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	containerImageCmdRepo := NewContainerImageCmdRepo(persistentDbSvc)
	containerImageQueryRepo := NewContainerImageQueryRepo(persistentDbSvc)

	t.Run("CreateSnapshot", func(t *testing.T) {
		containerQueryRepo := NewContainerQueryRepo(persistentDbSvc)

		readContainersRequestDto := dto.ReadContainersRequest{
			Pagination: useCase.ContainersDefaultPagination,
		}

		readContainersResponseDto, err := containerQueryRepo.Read(readContainersRequestDto)
		if err != nil || len(readContainersResponseDto.Containers) == 0 {
			t.Fatal(err)
		}
		containerEntity := readContainersResponseDto.Containers[0]

		createDto := dto.CreateContainerSnapshotImage{
			ContainerId: containerEntity.Id,
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
