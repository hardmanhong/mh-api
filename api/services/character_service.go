package services

import (
	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/models"
	"github.com/hardmanhong/api/utils"
)

type CharacterService interface {
	GetList() utils.ApiResponse
	GetItem(id uint32) utils.ApiResponse
	Create(character *models.Character) utils.ApiResponse
	Exists(id uint32) (bool, error)
	Update(id uint32, character *models.Character) utils.ApiResponse
	Delete(id uint32) utils.ApiResponse
}

type characterService struct {
	dao          *dao.CharacterDAO
	equipmentDao *dao.EquipmentDAO
	petDao       *dao.PetDAO
}

func NewCharacterService(dao *dao.CharacterDAO, equipment *dao.EquipmentDAO, petDao *dao.PetDAO) *characterService {
	return &characterService{dao, equipment, petDao}
}

func (s *characterService) GetList() utils.ApiResponse {
	res, err := s.dao.GetList()
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(res)
}

func (s *characterService) GetItem(id uint32) utils.ApiResponse {
	tx := s.dao.GetDB().Begin()
	item, err := s.dao.GetItem(id, tx)
	if err != nil {
		return utils.ApiErrorResponse(-1, "Failed to get item")
	}
	// 返回结果
	if item == nil {
		return utils.ApiErrorResponse(-1, "Item not found")
	}
	return utils.ApiSuccessResponse(&item)
}

func (s *characterService) Create(character *models.Character) utils.ApiResponse {
	db := s.dao.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	exists, err := s.dao.ExistsByName(character.Name)
	if err != nil {
		tx.Rollback()
		return utils.ApiErrorResponse(-1, err.Error())
	}
	if exists {
		tx.Rollback()
		return utils.ApiErrorResponse(-1, "角色名已存在")
	}
	result, err := s.dao.Create(character, tx)
	if err != nil || result.ID == 0 {
		tx.Rollback()
		return utils.ApiErrorResponse(-1, err.Error())
	}
	character.Equipment.CharacterID = result.ID
	// NOTE: 由于在结构体中定了 Equipment 和 Pet, gorm会自动执行创建这两个表的sql，导致下面报错重复创建
	// _, err = s.equipmentDao.Create(&character.Equipment, tx)
	// if err != nil {
	// 	tx.Rollback()
	// 	return utils.ApiErrorResponse(-1, err.Error())
	// }
	character.Pet.CharacterID = result.ID
	// _, err = s.petDao.Create(&character.Pet, tx)
	// if err != nil {
	// 	tx.Rollback()
	// 	return utils.ApiErrorResponse(-1, err.Error())
	// }
	// Commit the transaction
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(&result)
}

func (s *characterService) Exists(id uint32) (bool, error) {
	return s.dao.Exists(id)
}

func (s *characterService) Update(id uint32, character *models.Character) utils.ApiResponse {
	exists, err := s.Exists(id)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	if !exists {
		return utils.ApiErrorResponse(404, "记录不存在")
	}

	db := s.dao.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 使用 Preload 方法，预加载 Equipment 和 Pet 的关联数据
	find, err := s.dao.GetItem(id, tx)
	if err != nil {
		tx.Rollback()
		return utils.ApiErrorResponse(-1, err.Error())
	}
	err = s.dao.Update(id, character)
	if err != nil {
		tx.Rollback()
		return utils.ApiErrorResponse(-1, err.Error())
	}
	if find.Equipment.ID > 0 {
		err = s.equipmentDao.Update(find.Equipment.ID, &character.Equipment)
	} else {
		character.Equipment.CharacterID = find.ID
		_, err = s.equipmentDao.Create(&character.Equipment, tx)
	}
	if err != nil {
		tx.Rollback()
		return utils.ApiErrorResponse(-1, err.Error())
	}
	if find.Pet.ID > 0 {
		err = s.petDao.Update(find.Pet.ID, &character.Pet)
	} else {
		character.Pet.CharacterID = find.ID
		_, err = s.petDao.Create(&character.Pet, tx)
	}
	if err != nil {
		tx.Rollback()
		return utils.ApiErrorResponse(-1, err.Error())
	}
	// Commit the transaction
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(nil)
}

func (s *characterService) Delete(id uint32) utils.ApiResponse {
	err := s.dao.Delete(id)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(nil)
}
