package model

type BaseModel struct {
	ID uint `gorm:"primarykey;notNull;autoIncrement" json:"id,omitempty"`
}
