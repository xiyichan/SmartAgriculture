package model

type Admin struct {
	Uuid string `gorm:"unique;not null;primary_key"`
	NickName string `gorm:"DEFAULT:NULL"`
	Email string `gorm:"unique;DEFAULT:NULL"`
	Phone string `gorm:"unique;DEFAULT:NULL"`
	Password string `gorm:"DEFAULT:NULL"`
	Avater string `gorm:"DEFAULT:NULL"`
}
type AdminDto struct {
	Uuid string `json:"uuid"`
	NickName string `json:"nickName"`
	Phone string `json:"phone"`
	Email string `json:"Email"`
	Avater string `json:"avater"`

}
func ToAdminDto(admin Admin) AdminDto{
	return AdminDto{
		Uuid: admin.Uuid,
		NickName: admin.NickName,
		Phone: admin.Phone,
		Email:admin.Email,
		Avater: admin.Avater,
	}

}
