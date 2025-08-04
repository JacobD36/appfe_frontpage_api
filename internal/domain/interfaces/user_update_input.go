package interfaces

type UpdateUserInput interface {
	GetID() string
	FieldsToUpdate() map[string]any
}
