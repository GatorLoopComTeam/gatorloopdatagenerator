package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gatorloopdatagenerator/database"
	_ "github.com/go-sql-driver/mysql"
)

func random(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func getVelAccandPos(currentVel, currentAcc, currentPos float64) (float64, float64, float64) {
	if currentPos > 1609 {
		currentAcc -= 0.5
		if currentVel <= 0 {
			return 0, 0, currentPos
		}
	} else if currentAcc < 20 {
		currentAcc += 0.1
	} else if currentAcc > 20 && currentAcc < 30 {
		currentAcc += 0.1
	}

	currentVel = currentVel + currentAcc/20
	if currentVel > 250 {
		currentVel = 250
		currentAcc = 0
	}
	currentPos = currentPos + currentVel/20

	return currentVel, currentAcc, currentPos
}

var acceleration, velocity, position, pressure, roll, pitch, yaw, temperature float64
var pVoltage, pSOC, pTemp, pAmpHour float64
var aVoltage, aSOC, aTemp, aAmpHour float64
var startTime time.Time

func main() {
	database.InitDB()
	if len(os.Args) == 2 && os.Args[1] == "cleanup" {
		database.DB.Exec("DELETE FROM gatorloop.Acceleration")
		database.DB.Exec("DELETE FROM gatorloop.Position")
		database.DB.Exec("DELETE FROM gatorloop.Rotation")
		database.DB.Exec("DELETE FROM gatorloop.Temperature")
		database.DB.Exec("DELETE FROM gatorloop.Velocity")
		database.DB.Exec("DELETE FROM gatorloop.PrimaryBattery")
		database.DB.Exec("DELETE FROM gatorloop.AuxiliaryBattery")
		log.Info("Deleted all entries in database")
		return
	}

	pSOC = 1.0
	aSOC = 1.0
	rand.Seed(time.Now().Unix())
	startTime = time.Now()
	resetAcceleration := false
	for {
		database.DB.Exec("INSERT INTO gatorloop.Acceleration VALUES(NULL, " + fmt.Sprintf("%f", acceleration) + ")")
		database.DB.Exec("INSERT INTO gatorloop.Position VALUES(NULL, " + fmt.Sprintf("%f", position) + ")")
		database.DB.Exec("INSERT INTO gatorloop.Pressure VALUES(NULL, " + fmt.Sprintf("%f", pressure) + ")")
		database.DB.Exec("INSERT INTO gatorloop.Rotation VALUES(NULL, " + fmt.Sprintf("%f,%f,%f", roll, pitch, yaw) + ")")
		database.DB.Exec("INSERT INTO gatorloop.Temperature VALUES(NULL, " + fmt.Sprintf("%f", temperature) + ")")
		database.DB.Exec("INSERT INTO gatorloop.Velocity VALUES(NULL, " + fmt.Sprintf("%f", velocity) + ")")
		database.DB.Exec("INSERT INTO gatorloop.PrimaryBattery VALUES(NULL," + fmt.Sprintf("%f,%f,%f,%f", pVoltage, pSOC, pTemp, pAmpHour) + ")")
		database.DB.Exec("INSERT INTO gatorloop.AuxiliaryBattery VALUES(NULL," + fmt.Sprintf("%f,%f,%f,%f", aVoltage, aSOC, aTemp, aAmpHour) + ")")

		velocity, acceleration, position = getVelAccandPos(velocity, acceleration, position)
		if position > 1609 && !resetAcceleration {
			acceleration = 0
			resetAcceleration = true
		}
		pressure = random(40, 60)
		roll = random(0, 1)
		pitch = random(0, 1)
		yaw = random(0, 1)
		temperature = random(75, 80)
		pVoltage = random(4.8, 5)
		pSOC -= random(.0001, .001)
		pTemp = random(75, 80)
		pAmpHour = random(180, 200)
		aVoltage = random(4.8, 5)
		if pSOC <= 0 {
			aSOC -= random(.0001, .001)
		}
		aTemp = random(75, 80)
		aAmpHour = random(180, 200)
		time.Sleep(time.Millisecond * 50)
	}
}
