package utils

import (
	"strconv"
	"strings"
	"time"
	"user_role_permissions/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func ParseTurnoverYears(c *gin.Context, index int) []model.TurnoverYear {

	var years []model.TurnoverYear

	yIndex := 0
	for {
		yearKey := "data[" + strconv.Itoa(index) + "][years][" + strconv.Itoa(yIndex) + "][year]"
		amountKey := "data[" + strconv.Itoa(index) + "][years][" + strconv.Itoa(yIndex) + "][amount]"

		yearVal := c.PostForm(yearKey)
		amountVal := c.PostForm(amountKey)

		if yearVal == "" {
			break
		}

		amount, _ := strconv.ParseFloat(amountVal, 64)

		years = append(years, model.TurnoverYear{
			FinancialYear:   yearVal,
			Amount: amount,
		})

		yIndex++
	}

	return years
}

func ParseDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}
	layout := "02/01/2006"
	parsed, err := time.Parse(layout, dateStr)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func MaskMobile(mobile string) string {
    if len(mobile) <= 4 {
        return mobile
    }
    return strings.Repeat("*", len(mobile)-4) + mobile[len(mobile)-4:]
}

func MaskAadhar(aadhar string) string {
    if len(aadhar) <= 4 {
        return aadhar
    }
    return strings.Repeat("*", len(aadhar)-4) + aadhar[len(aadhar)-4:]
}

func MaskEmail(email string) string {
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return email
    }

    name := parts[0]
    domain := parts[1]

    if len(name) <= 2 {
        return strings.Repeat("*", len(name)) + "@" + domain
    }

    return name[:2] + strings.Repeat("*", len(name)-2) + "@" + domain
}