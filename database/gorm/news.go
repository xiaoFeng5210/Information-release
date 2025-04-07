package database

import (
	model "infomation-release/database/model"
	"time"
)

func PostNews(uid int, title string, content string) (int, error) {
	now := time.Now()

	news := &model.News{
		UserId:     uid,
		Title:      title,
		Content:    content,
		PostTime:   &now,
		DeleteTime: nil,
	}

}
