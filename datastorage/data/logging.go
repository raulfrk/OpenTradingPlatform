package data

type Log struct {
	Message string `gorm:"type:json;default:'[]';not null"`
}

// TODO: Handle errors on insertions
func InsertLog(message []byte) {
	data := Log{
		Message: string(message),
	}
	DB.Create(&data).Commit()
}
