package userdto

import "go-server/internal/model"

// ================= DTO 转换 =================

// 👉 他人可见
func ToUserPublicDTO(u *model.User) UserPublicDTO {
	return UserPublicDTO{
		ID:       u.ID,
		Username: u.Username,
		Avatar:   u.Avatar,
	}
}

// 👉 自己可见
func ToUserPrivateDTO(u *model.User) UserPrivateDTO {
	return UserPrivateDTO{
		ID:       u.ID,
		Username: u.Username,
		Avatar:   u.Avatar,
	}
}

func UserListToPublic(users []model.User) []UserPublicDTO {
	list := make([]UserPublicDTO, 0, len(users))
	for _, u := range users {
		list = append(list, ToUserPublicDTO(&u))
	}
	return list
}
