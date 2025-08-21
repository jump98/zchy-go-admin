package version

import (
	"fmt"
	"runtime"

	"gorm.io/gorm"

	"go-admin/cmd/migrate/migration"
	"go-admin/cmd/migrate/migration/models"
	common "go-admin/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1752490642079Test)
}

func _1752490642079Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		fmt.Println("创建数据库的表=================================")
		err := tx.Debug().Migrator().AutoMigrate(
			new(models.SysRadar),
			new(models.RadarPoint),
		)
		if err != nil {
			fmt.Println("创建数据库的表出错:", err)
			return err
		}

		return tx.Create(&common.Migration{

			Version: version,
		}).Error
	})
}
