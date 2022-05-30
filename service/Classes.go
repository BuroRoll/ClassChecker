package service

import (
	"ClassChecker/models"
	"ClassChecker/repository"
	"bytes"
	"encoding/json"
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
		classData := c.repo.GetClassData(j.BookingClassID)
		if err != nil {
			return
		}
		mentorData := CreateDataToSend(classData.MentiFirstName, classData.MentiSecondName, classData.ClassDataName, classData.ID, classData.UserID, classData.MentiId, j.Time, 0)
		mentorNotification := models.ClassNotification{
			Receiver: classData.UserID,
			Type:     "lesson complete",
			Data:     mentorData,
		}
		mentiData := CreateDataToSend(classData.MentorFirstName, classData.MentorSecondName, classData.ClassDataName, classData.ID, classData.UserID, classData.MentiId, j.Time, 0)
		mentiNotification := models.ClassNotification{
			Receiver: classData.MentiId,
			Type:     "lesson complete",
			Data:     mentiData,
		}
		data1 := c.repo.SaveNotification(mentiNotification)
		data2 := c.repo.SaveNotification(mentorNotification)
		jsonNotification, _ := json.Marshal(data1)
		sendToServer(string(jsonNotification), classData.MentiId)
		jsonNotification, _ = json.Marshal(data2)
		sendToServer(string(jsonNotification), classData.UserID)
	}
}

type booking_data struct {
	ClassDataName    string `json:"class_data_name"`
	FirstName        string `json:"first_name"`
	SecondName       string `json:"second_name"`
	BookingId        uint   `json:"booking_id"`
	MentorId         uint   `json:"mentor_id"`
	CommentRecipient uint   `json:"comment_recipient"`
	Time             string `json:"time"`
	LessonCount      uint   `json:"lesson_count"`
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
			mentorData := CreateDataToSend(i.MentiFirstName, i.MentiSecondName, i.ClassDataName, i.ID, i.UserID, i.MentiId, "", i.LessonCount)
			mentorNotification := models.ClassNotification{
				Receiver: i.UserID,
				Type:     "class comleted",
				Data:     mentorData,
			}
			mentiData := CreateDataToSend(i.MentorFirstName, i.MentorSecondName, i.ClassDataName, i.ID, i.UserID, i.UserID, "", i.LessonCount)
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

func CreateDataToSend(firstName string, secondName string, className string, bookingId uint, mentorId uint, commentRecipient uint, time string, lessonCount uint) string {
	var data booking_data
	if time == "" {
		data = booking_data{
			ClassDataName:    className,
			FirstName:        firstName,
			SecondName:       secondName,
			BookingId:        bookingId,
			MentorId:         mentorId,
			CommentRecipient: commentRecipient,
			LessonCount:      lessonCount,
		}
	} else {
		data = booking_data{
			ClassDataName:    className,
			FirstName:        firstName,
			SecondName:       secondName,
			BookingId:        bookingId,
			MentorId:         mentorId,
			CommentRecipient: commentRecipient,
			Time:             time,
		}
	}

	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}
