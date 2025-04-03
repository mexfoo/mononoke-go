package database

import (
	"errors"
	"mononoke-go/model"
	"mononoke-go/utils"

	"gorm.io/gorm"
)

func (d *GormDatabase) GetUserByNameAndPW(name, password string) (*model.Accounts, bool) {
	account := new(model.Accounts)
	err := d.DB.Where("account_name = ?", name).Find(account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false
	}

	if !utils.VerifyPassword(password, account.Password) {
		return nil, false
	}

	if account.AccountName == name {
		return account, true
	}
	return nil, false
}

func (d *GormDatabase) UpdateLastLoginServerIdx(accountID, lastLoginServerIdx uint32) (bool, error) {
	err := d.DB.Model(model.Accounts{}).
		Where("account_id", accountID).
		Update("last_login_server_idx", lastLoginServerIdx).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
