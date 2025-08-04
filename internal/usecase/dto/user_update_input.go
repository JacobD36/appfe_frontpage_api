package dto

type UpdateUserInput struct {
	ID             string
	Name           *string
	Password       *string
	Img            *string
	Role           *string
	Status         *bool
	EmailValidated *bool
}

func (u UpdateUserInput) GetID() string {
	return u.ID
}

func (u UpdateUserInput) FieldsToUpdate() map[string]any {
	fields := make(map[string]any)

	if u.Name != nil {
		fields["name"] = *u.Name
	}
	if u.Password != nil {
		fields["password"] = *u.Password
	}
	if u.Img != nil {
		fields["img"] = *u.Img
	}
	if u.Role != nil {
		fields["role"] = *u.Role
	}
	if u.Status != nil {
		fields["status"] = *u.Status
	}
	if u.EmailValidated != nil {
		fields["email_validated"] = *u.EmailValidated
	}
	return fields
}
