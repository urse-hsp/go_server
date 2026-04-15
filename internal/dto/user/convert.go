package userdto

import "go-server/internal/model"

// ================= DTO 转换 =================

// 👉 他人可见
func ToPublicDTO(u *model.User) PrivateDTO {
	return PrivateDTO{
		ID:       u.ID,
		Username: u.Username,
		Avatar:   u.Avatar,
	}
}

// 👉 自己可见
func ToPrivateDTO(u *model.User) PrivateDTO {
	return PrivateDTO{
		ID:       u.ID,
		Username: u.Username,
		Avatar:   u.Avatar,
	}
}

func ListToPublic(users []model.User) []PrivateDTO {
	list := make([]PrivateDTO, 0, len(users))
	for _, u := range users {
		list = append(list, ToPublicDTO(&u))
	}
	return list
}
