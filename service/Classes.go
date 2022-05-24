package service

import (
	"ClassChecker/models"
	"ClassChecker/repository"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ClassesService struct {
	repo repository.Classes
}

func NewClassesService(repo repository.Classes) *ClassesService {
	return &ClassesService{repo: repo}
}

func (c ClassesService) CheckClasses() {
	classesTimes := c.repo.GetClassesTimes()
	for _, j := range classesTimes {
		j.IsEnd = true
		err := c.repo.SaveClassesTimes(j)
		if err != nil {
			return
		}
	}
}

type booking_data struct {
	ClassName        string `json:"class_name"`
	FirstName        string `json:"first_name"`
	SecondName       string `json:"second_name"`
	BookingId        uint   `json:"booking_id"`
	MentorId         uint   `json:"mentor_id"`
	CommentRecipient uint   `json:"comment_recipient"`
}

func (c ClassesService) CheckClassEnd() {
	classesTimes := c.repo.GetClassesWithTime()
	for _, i := range classesTimes {
		isEnd := true
		for _, j := range i.Time {
			if !j.IsEnd {
				isEnd = false
				break
			}
		}
		if isEnd {
			i.Status = "completed"
			c.repo.SaveClass(i)
			mentorData := CreateDataToSend(i.MentiFirstName, i.MentiSecondName, i.ClassDataName, i.ID, i.UserID, i.MentiId)
			mentorNotification := models.ClassNotification{
				Receiver: i.UserID,
				Type:     "class comleted",
				Data:     mentorData,
			}
			mentiData := CreateDataToSend(i.MentorFirstName, i.MentorSecondName, i.ClassDataName, i.ID, i.UserID, i.UserID)
			mentiNotification := models.ClassNotification{
				Receiver: i.MentiId,
				Type:     "class comleted",
				Data:     mentiData,
			}
			data1 := c.repo.SaveNotification(mentiNotification)
			data2 := c.repo.SaveNotification(mentorNotification)
			jsonNotification, _ := json.Marshal(data1)
			sendToServer(string(jsonNotification), i.MentiId)
			jsonNotification, _ = json.Marshal(data2)
			sendToServer(string(jsonNotification), i.UserID)
		}
	}
}

func sendToServer(data string, userId uint) {
	type DataToServer struct {
		Data   string
		UserId uint
	}
	d := DataToServer{
		Data:   data,
		UserId: userId,
	}
	json_data, _ := json.Marshal(d)

	//http.Post("http://localhost:8000/notifications/class/", "application/json",
	http.Post("http://152.70.189.77/backend/notifications/class", "application/json",
		bytes.NewBuffer(json_data))
}

func CreateDataToSend(firstName string, secondName string, className string, bookingId uint, mentorId uint, commentRecipient uint) string {
	data := booking_data{
		ClassName:        className,
		FirstName:        firstName,
		SecondName:       secondName,
		BookingId:        bookingId,
		MentorId:         mentorId,
		CommentRecipient: commentRecipient,
	}
	jsonData, _ := json.Marshal(data)
	fmt.Println(string(jsonData))
	return string(jsonData)
}
