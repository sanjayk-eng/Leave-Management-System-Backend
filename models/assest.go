package models

// EQUIPMENT

//Catagory

type EquipmentCategoryInput struct {
	Name        string  `json:"name" validate:"required,min=2,max=50"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
}
