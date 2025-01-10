package model


// DB Models

// GraphUser

// GraphUser is a user of the graph database
type GraphUser struct { // To track what user has what entities in his personal graph
	ID					string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" validate:"required,uuid"`

	UsersEntities		[]UsersEntity `gorm:"foreignKey:UserID;references:ID"`
	UsersConnectionTypes []UsersConnectionType `gorm:"foreignKey:UserID;references:ID"`
	UsersPropertyTypes	[]UsersPropertyType `gorm:"foreignKey:UserID;references:ID"`
}

func (gu *GraphUser) TableName() string {
	return "users"
}

// Entity
// Entity is a node in the graph (same for all users)
type Entity struct {
	ID					string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" validate:"required,uuid"`

	UsersEntities		[]UsersEntity `gorm:"foreignKey:EntityID;references:ID"`
}

func (e *Entity) TableName() string {
	return "entities"
}

// UsersEntity is a specific users version of an entity
type UsersEntity struct { // junction table for many-to-many relationship between User and Entity
	Name				string `gorm:"type:varchar(255);not null;index:idx_users_entity_name" validate:"required,max=255"`
	Definition			string `gorm:"type:varchar(4096);not null" validate:"required,max=4096"`
	UserID				string `gorm:"type:uuid;primaryKey" validate:"required,uuid"`
	EntityID			string `gorm:"type:uuid;primaryKey" validate:"required,uuid"`
}

func (ue *UsersEntity) TableName() string {
	return "users_entities"
}

// ConnectionType
// ConnectionType is a type of a edge in the graph (same for all users)
type ConnectionType struct { 
	ID					string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" validate:"required,uuid"`

	UsersConnectionTypes []UsersConnectionType `gorm:"foreignKey:ConnectionTypeID;references:ID"`
	// TODO Add InverseConnectionType
}

func (ct *ConnectionType) TableName() string {
	return "connection_types"
}

// UsersConnectionType is a specific users version of a connection type
type UsersConnectionType struct { // junction table for many-to-many relationship between User and ConnectionType
	Name       string `gorm:"type:varchar(255);not null;index:idx_users_connection_type" validate:"required,max=255"`
	Definition string `gorm:"type:varchar(4096);not null" validate:"required,max=4096"`

	UserID           string `gorm:"type:uuid;primaryKey"`
	ConnectionTypeID string `gorm:"type:uuid;primaryKey"`
}

func (uct *UsersConnectionType) TableName() string {
	return "users_connection_types"
}

// PropertyType
// PropertyType is a type of a property in the graph (same for all users)
type PropertyType struct { 
	ID string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" validate:"required,uuid"`
	// The value type of the property
	ValueType string `gorm:"type:varchar(10);check:value_type in ('string','int','float','boolean')" validate:"required,oneof=string int float boolean"`

	UsersPropertyTypes []UsersPropertyType `gorm:"foreignKey:PropertyTypeID;references:ID"`
}

func (pt *PropertyType) TableName() string {
	return "property_types"
}

// UsersPropertyType is a specific users version of a property type
type UsersPropertyType struct { // junction table for many-to-many relationship between User and PropertyType
	Name       string `gorm:"type:varchar(255);not null;index:idx_users_property_type" validate:"required,max=255"`
	Definition string `gorm:"type:varchar(4096);not null" validate:"required,max=4096"`

	UserID         string `gorm:"type:uuid;primaryKey" validate:"required,uuid"`
	PropertyTypeID string `gorm:"type:uuid;primaryKey" validate:"required,uuid"`
}

func (upt *UsersPropertyType) TableName() string {
	return "users_property_types"
}


// DTOs

// Common

type SearchRequest struct {
	Name string `json:"name" validate:"required,max=255"`
}

type UserDataResponse struct {
	ID string `json:"id" validate:"required,uuid"`

	Entities        []UsersEntity `json:"entities"`
	ConnectionTypes []UsersConnectionType `json:"connection_types"`
	PropertyTypes   []PropertyTypeResponse `json:"property_types"`
}

// GraphUser

type UserRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

// Entity

type EntityRequest struct {
	ID         string `json:"id" validate:"omitempty,uuid"`
	Name       string `json:"name" validate:"required,max=255"`
	Definition string `json:"definition" validate:"required,max=4096"`
}

// ConnectionType

type ConnectionTypeRequest struct {
	ID         string `json:"id" validate:"omitempty,uuid"`
	Name       string `json:"name" validate:"required,max=255"`
	Definition string `json:"definition" validate:"required,max=4096"`
}

// PropertyType

type PropertyTypeRequest struct {
	ID         string `json:"id" validate:"omitempty,uuid"`
	Name       string `json:"name" validate:"required,max=255"`
	Definition string `json:"definition" validate:"required,max=4096"`
	ValueType  string `json:"value_type" validate:"required,oneof=string int float boolean"`
}

type PropertyTypeResponse struct {
	UserID         string `json:"user_id" validate:"required,uuid"`
	PropertyTypeID string `json:"property_type_id" validate:"required,uuid"`
	Name       string `json:"name" validate:"required,max=255"`
	Definition string `json:"definition" validate:"required,max=4096"`
	ValueType  string `json:"value_type" validate:"required,oneof=string int float boolean"`
}
