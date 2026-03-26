// model 集中注册
package model

func GetModels() []interface{} {
	return []interface{}{
		&User{},
		// &Order{},
		// &Match{},
	}
}
