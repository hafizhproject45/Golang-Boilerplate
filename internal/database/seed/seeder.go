package seed

import (
	"fmt"

	mUser "github.com/hafizhproject45/Golang-Boilerplate.git/internal/modules/users/models"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// pw, err := secure.Hash("asdasdasd", nil)

		// if err != nil {
		// 	return err
		// }

		// ===== Users (user) =====
		user := mUser.User{
			Name: "Super Admin",
		}
		if err := tx.Where("email = ?", user.Id).FirstOrCreate(&user).Error; err != nil {
			return err
		}

		fmt.Println("✅ Seeder successfully")
		return nil
	})
}
