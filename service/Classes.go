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
		if err != nil {
			return
		}
	}
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
			mentorNotification := models.Notification{
				Receiver: i.UserID,
				Type:     "class comleted",
				Data:     i.ClassDataName,
			}
			mentiNotification := models.Notification{
				Receiver: i.MentiId,
				Type:     "class comleted",
				Data:     i.ClassDataName,
			}
			data1 := c.repo.SaveNotification(mentiNotification)
			data2 := c.repo.SaveNotification(mentorNotification)
			jsonNotification, _ := json.Marshal(data1)
			sendToServer(string(jsonNotification), i.UserID)
			jsonNotification, _ = json.Marshal(data2)
			sendToServer(string(jsonNotification), i.MentiId)
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

	http.Post("http://152.70.189.77/backend/notifications/", "application/json",
		bytes.NewBuffer(json_data))
}
