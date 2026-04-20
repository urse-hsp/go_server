package userdto

import "go-server/internal/model"

// ================= DTO 转换 =================

// 他人可见
func ToPublicDTO(u *model.User) PublicDTO {
	return PublicDTO{
		ID:       u.ID,
		Username: u.Username,
		Avatar:   u.Avatar,
	}
}

// 自己可见
func ToPrivateDTO(u *model.User) PrivateDTO {
	return PrivateDTO{
		PublicDTO: ToPublicDTO(u),
	}
}

func ListToPublic(users []model.User) []PublicDTO {
	list := make([]PublicDTO, 0, len(users))
	for _, u := range users {
		list = append(list, ToPublicDTO(&u))
	}
	return list
}
