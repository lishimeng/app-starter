package persistencetest

// TestRecord 集成测试用表模型。
type TestRecord struct {
	Id     int    `gorm:"primaryKey;autoIncrement;column:id"`
	Name   string `gorm:"column:name"`
	Status int    `gorm:"column:status;default:0"`
}

func (TestRecord) TableName() string {
	return "test_record"
}
