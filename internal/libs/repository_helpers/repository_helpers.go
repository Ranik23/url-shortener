package repositoryhelpers

import (
	"fmt"
)


func GetConnectionString(Type, Host, Port, User, Password, Name string) (connectionString string) {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
			Host, Port, User, Password, Name, "Europe/Moscow")
}