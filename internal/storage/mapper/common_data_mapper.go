package mapper

import "github.com/rinefica/voice_null_files/internal/domain/model"

func MapFileToCommon(files []model.FileModel) []*model.CommonData {
	var data []*model.CommonData

	for _, file := range files {
		d := &model.CommonData{
			UUID:           file.UUID,
			AdditionalData: file.Filename,
			Type:           typeFile,
		}
		data = append(data, d)
	}
	return data
}
func MapInfoDataToCommon(infoDataModels []model.InfoDataModel) []*model.CommonData {
	var data []*model.CommonData

	for _, infoData := range infoDataModels {
		d := &model.CommonData{
			UUID:           infoData.UUID,
			AdditionalData: infoData.AdditionalData,
			Type:           infoData.Type,
		}
		data = append(data, d)
	}
	return data
}

const typeFile = "file"
