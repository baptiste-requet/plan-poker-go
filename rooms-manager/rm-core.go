package roomsmanager

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const ROOM_CODE_SIZE = 4 // TODO: in conf

type User struct {
	Uuid     string
	Nickname string
	Estimate string
}

type Room struct {
	Users              map[*User]bool
	Name               string
	Code               string // must be unique, can act as an ID
	IsEstimateRevealed bool
}

var rooms = make(map[string]*Room)

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func CreateRoom(roomName string) *Room {
	if roomName == "" {
		return nil
	}
	roomCode := generateUniqueRoomCode()
	log.Printf("Creating room with name '%s', generated code : '%s'", roomName, roomCode)
	createdRoom := &Room{
		Users:              make(map[*User]bool),
		Name:               roomName,
		Code:               roomCode,
		IsEstimateRevealed: false,
	}
	rooms[roomCode] = createdRoom
	return createdRoom
	// TODO: timeout 5min no one connected -> delete room
}

func FindRoomByRoomCode(roomCode string) *Room {
	room, exist := rooms[roomCode]
	if !exist {
		log.Printf("Room with code '%s' not found", roomCode)
		return nil
	}
	return room
}

func GetAllUserFromRoomByRoomCode(roomCode string) (users []*User) {
	room := FindRoomByRoomCode(roomCode)
	if room == nil {
		return users
	}

	for user, ok := range room.Users {
		if ok {
			users = append(users, user)
		}
	}
	return users
}

func ConnectNewUserToRoom(nickname, roomCode string) *User { // TODO: return errors saying why user cant connect (ex: room doesn't exist)
	room := FindRoomByRoomCode(roomCode)
	if room == nil {
		return nil
	}

	newUser := createUser(nickname)
	room.Users[newUser] = true
	log.Printf("User '%s' joined the room : '%s'", nickname, roomCode)

	return newUser
}

func DisconnectUserFromRoom(user *User, roomCode string) {
	room := FindRoomByRoomCode(roomCode)
	if room == nil {
		return
	}

	log.Printf("User '%s' disconnected from room '%s'", user.Nickname, roomCode)
	delete(room.Users, user)
	// TODO: timeout if no user left -> delete room
}

func SubmitEstimate(user *User, roomCode string, estimate string) (err error) {
	room := FindRoomByRoomCode(roomCode)
	if room == nil {
		errMsg := fmt.Sprintf("Room [%s] not found, cannot submit user %s estimate (%s)", roomCode, user.Uuid, estimate)
		return errors.New(errMsg)
	}

	log.Printf("User %s submitted estimate {%s} in room [%s]", user.Uuid, estimate, roomCode)
	user.Estimate = estimate

	return nil
}

func ToggleShouldRevealEstimateForRoom(roomCode string) (bool, error) {
	room := FindRoomByRoomCode(roomCode)
	if room == nil {
		errMsg := fmt.Sprintf("Room [%s] not found, cannot toggle estimate reveal", roomCode)
		return false, errors.New(errMsg)
	}

	room.IsEstimateRevealed = !room.IsEstimateRevealed
	log.Printf("Toggled estimate reveal to '%t' for room [%s]", room.IsEstimateRevealed, roomCode)

	return room.IsEstimateRevealed, nil
}

func ResetPlanningForRoom(roomCode string)  error {
	room := FindRoomByRoomCode(roomCode)
	if room == nil {
		errMsg := fmt.Sprintf("Room [%s] not found, cannot reset planning", roomCode)
		return errors.New(errMsg)
	}

	room.IsEstimateRevealed = false
	for user := range room.Users {
		user.Estimate = ""
	}
	
	log.Printf("Reseted planning for room [%s]", roomCode)

	return nil
}

func createUser(nickname string) *User {
	// Get the string representation of the UUID
	uuidStr := uuid.New().String()
	log.Printf("Creating user with nickname '%s', uuid : '%s'", nickname, uuidStr)
	return &User{
		Nickname: nickname,
		Uuid:     uuidStr,
	}
}

func generateUniqueRoomCode() string {
	var generatedRoomCode string
	for {
		generatedRoomCode = generateRoomCode()
		if _, exist := rooms[generatedRoomCode]; !exist {
			break // Exit the loop if the generated code is unique
		}
	}
	return generatedRoomCode
}

func generateRoomCode() string {
	b := make([]rune, ROOM_CODE_SIZE)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
