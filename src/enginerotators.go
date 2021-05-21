package main

type engineRotator interface {
	Rotate(credentials databaseCredentials, newPassword string) error
	Test(credentials databaseCredentials) error
}

var postgresRotator = postgresEngineRotator{}
var sqlserverRotator = sqlserverEngineRotator{}
var engineRotators = map[string]engineRotator{
	"postgres":  &postgresRotator,
	"sqlserver": &sqlserverRotator,
}
