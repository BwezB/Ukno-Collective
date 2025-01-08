package model

// DB Models
type User struct { // To track what user has what entities in his personal graph
	ID string `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	Entities        []Entity         `gorm:"many2many:user_entities;"`
	ConnectionTypes []ConnectionType `gorm:"many2many:user_connection_types;"`
	PropertyTypes   []PropertyType   `gorm:"many2many:user_property_types;"`
}

type Entity struct { // A node in the graph, the same for all users
	ID string `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	Users []User `gorm:"many2many:user_entities;"`
}

type UserEntity struct { // junction table for many-to-many relationship between User and Entity
	Name       string `gorm:"not null"`
	Definition string `gorm:"not null"`

	UserID   string `gorm:"type:uuid;primaryKey"`
	EntityID string `gorm:"type:uuid;primaryKey"`
}

type ConnectionType struct { // A type of a edge in the graph, the same for all users
	ID string `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	Users []User `gorm:"many2many:user_connection_types;"`
}

type UserConnectionType struct { // junction table for many-to-many relationship between User and ConnectionType
	Name       string `gorm:"not null"`
	Definition string `gorm:"not null"`

	UserID           string `gorm:"type:uuid;primaryKey"`
	ConnectionTypeID string `gorm:"type:uuid;primaryKey"`
}

type PropertyType struct { // A type of a property in the graph, the same for all users
	ID string `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	// The value type of the property
	ValueType string `gorm:"not null;check:value_type in ('string','int','float','boolean')"`

	Users []User `gorm:"many2many:user_property_types;"`
}

type UserPropertyType struct { // junction table for many-to-many relationship between User and PropertyType
	Name       string `gorm:"not null"`
	Definition string `gorm:"not null"`

	UserID         string `gorm:"type:uuid;primaryKey"`
	PropertyTypeID string `gorm:"type:uuid;primaryKey"`
}

// DTOs

type UserRequest struct {
	ID string `json:"id" validate:"required"`
}

type EntityRequest struct {
	ID         string `json:"id" validate:"-"`
	Name       string `json:"name" validate:"required"`
	Definition string `json:"definition" validate:"required"`
}

type SearchRequest struct {
	Name string `json:"name" validate:"required"`
}

type ConnectionTypeRequest struct {
	ID         string `json:"id" validate:"-"`
	Name       string `json:"name" validate:"required"`
	Definition string `json:"definition" validate:"required"`
}

type PropertyTypeRequest struct {
	ID         string `json:"id" validate:"-"`
	Name       string `json:"name" validate:"required"`
	Definition string `json:"definition" validate:"required"`
	ValueType  string `json:"value_type" validate:"required,oneof=string int float boolean"`
}
