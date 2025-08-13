package version

import (
	"runtime"

	"gorm.io/gorm"

	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1752583284513Test)
}

func _1752583284513Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		// TODO: 这里开始写入要变更的内容

		// TODO: 例如 修改表字段 使用过程中请删除此段代码
		//err := tx.Migrator().RenameColumn(&models.SysConfig{}, "config_id", "id")
		//if err != nil {
		// 	return err
		//}

		// TODO: 例如 新增表结构 使用过程中请删除此段代码
		//err = tx.Debug().Migrator().AutoMigrate(
		//		new(models.CasbinRule),
		// 		)
		//if err != nil {
		// 	return err
		//}

		//这次升级文件还没有用，下次修改后用参数"args": ["migrate"]运行一次即可

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
