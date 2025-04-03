package database

import (
	"errors"
	"mononoke-go/config"
	"mononoke-go/model"
	"mononoke-go/utils"

	"gorm.io/gorm"
)

func (d *GormDatabase) GetUserByNameAndPW(name, password string, conf *config.Configuration) (*model.Accounts, bool) {
	account := new(model.Accounts)
	err := d.DB.Where("account_name = ?", name).Find(account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false
	}

	if account.AccountName != name {
		return nil, false
	}

	if account.Password == "" {
		// In case the bcrypt password is wrong, try migration from an md5 password?
		if !conf.Database.RunPasswordMigration {
			return nil, false
		}
		if !utils.VerifyPasswordMD5(password, account.PasswordMD5) {
			return nil, false
		}
		account.Password, err = utils.HashPassword(password)
		if err != nil {
			return nil, false
		}

		account.PasswordMD5 = ""
		d.DB.Save(&account)
	}

	if !utils.VerifyPassword(password, account.Password) {
		return nil, false
	}

	return account, true
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
