package repository

import (
	"service3/dto"
	"time"
)

type serviceResult struct {
	Id              int
	ExecutionTime   time.Time
	Service1Success bool
	Service2Success bool
	Service1Delay   int64
	Service2Delay   int64
}

type ServiceResultRepository interface {
	Save(dto dto.ServiceResultDto) error
}

type serviceResultRepositoryImpl struct {
	baseRepository
}

func GetServiceResultRepository() ServiceResultRepository {
	return &serviceResultRepositoryImpl{baseRepository{dbConn: GetDBConn()}}
}

func (repo *serviceResultRepositoryImpl) Save(resultDto dto.ServiceResultDto) error {
	serviceResultRecord := serviceResult{
		ExecutionTime:   resultDto.ExecutionTime,
		Service1Success: resultDto.Service1Success,
		Service2Success: resultDto.Service2Success,
		Service1Delay:   resultDto.Service1Delay,
		Service2Delay:   resultDto.Service2Delay,
	}

	err := repo.dbConn.Create(&serviceResultRecord).Error
	return err
}
