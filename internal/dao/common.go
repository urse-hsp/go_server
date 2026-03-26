// 通用dao
package dao

import "gorm.io/gorm"

// 通用分页 page从1开始
func Paginate(db *gorm.DB, out interface{}, page, pageSize int) (int64, error) {
	var total int64

	// ⚠️ 注意：Count 要在 Offset/Limit 之前
	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}

	err := db.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(out).Error

	return total, err
}
