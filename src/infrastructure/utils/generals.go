package utils

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func GetTimeZoneThailand() time.Time {
	// Set the Bangkok, Thailand time zone
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return time.Now().Add(7 * time.Hour)
	}

	return time.Now().In(loc)
}

func CreateUUID(id string) string {
	uID := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(id))
	return uID.String()
}
