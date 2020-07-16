package model


type User struct {
	Uuid string `gorm:"unique;not null;primary_key"`
	NickName string `gorm:"DEFAULT:NULL"`
	Phone string `gorm:"unique;DEFAULT:NULL"`
	Email string `gorm:"unique;DEFAULT:NULL"`
	Password string `gorm:"DEFAULT:NULL"`
	Avater string `gorm:"DEFAULT:NULL"`
}

type UserDto struct {
	Uuid string `json:"uuid"`
	NickName string `json:"nickName"`
	Phone string `json:"phone"`
	Email string `json:"Email"`
	Avater string `json:"avater"`

}
func ToUserDto(user User) UserDto{
	return UserDto{
		Uuid: user.Uuid,
		NickName: user.NickName,
		Phone: user.Phone,
		Email:user.Email,
		Avater: user.Avater,
	}

}

