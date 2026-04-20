package demodto

import "go-server/internal/model"

// ================= DTO 转换 =================

// 他人可见
func ToPublicDTO(u *model.Demo) PublicDTO {
	return PublicDTO{}
}

// 自己可见
func ToPrivateDTO(u *model.Demo) PrivateDTO {
	return PrivateDTO{
		PublicDTO: ToPublicDTO(u),
	}
}

func ListToPublic(users []model.Demo) []PublicDTO {
	list := make([]PublicDTO, 0, len(users))
	for _, u := range users {
		list = append(list, ToPublicDTO(&u))
	}
	return list
}
