package model

type Accounts struct {
	AccountID          uint32 `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	AccountName        string `gorm:"type:varchar(61);unique_index"`
	Password           string `gorm:"type:varchar(60)"`
	Email              string `gorm:"type:varchar(32);unique_index"`
	Blocked            bool
	Age                uint8
	LastLoginServerIdx uint32
	Permission         uint32
}
