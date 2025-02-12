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

	t.Run("CreateArchive", func(t *testing.T) {
		imagesList, err := containerImageQueryRepo.Read()
		if err != nil {
			t.Fatal(err)
		}
		if len(imagesList) == 0 {
			t.Fatal("NoImagesFound")
		}

		createDto := dto.CreateContainerImageArchive{
			AccountId: imagesList[0].AccountId,
			ImageId:   imagesList[0].Id,
		}
		_, err = containerImageCmdRepo.CreateArchive(createDto)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("DeleteArchive", func(t *testing.T) {
		responseDto, err := containerImageQueryRepo.ReadArchives(
			dto.ReadContainerImageArchivesRequest{},
		)
		if err != nil {
			t.Fatal(err)
		}
		if len(responseDto.Archives) == 0 {
			t.Fatal("NoArchiveFilesFound")
		}

		err = containerImageCmdRepo.DeleteArchive(responseDto.Archives[0])
		if err != nil {
			t.Fatal(err)
		}
	})
}
