package database

import (
	"errors"
	"fmt"
	model "infomation-release/database/model"
	"log/slog"
	"time"
)

// 创建新闻
func PostNews(uid int, title string, content string) (int, error) {
	now := time.Now()

	news := &model.News{
		UserId:     uid,
		Title:      title,
		Content:    content,
		PostTime:   &now,
		DeleteTime: nil,
	}
	err := globalDB.Create(news).Error
	if err != nil {
		slog.Error("发布新闻失败", "error", err)
		return 0, errors.New("发布新闻失败")
	}
	return news.Id, nil
}

func DeleteNews(id int) error {
	tx := globalDB.Model(&model.News{}).Where("id = ?", id).Update("delete_time", time.Now())
	if tx.Error != nil {
		slog.Error("更新新闻失败", "error", tx.Error)
		return errors.New("更新新闻失败")
	} else {
		if tx.RowsAffected <= 0 {
			return fmt.Errorf("新闻id[%d]不存在", id)
		} else {
			return nil
		}
	}
}

func UpdateNews(id int, title, content string) error {
	tx := globalDB.Model(&model.News{}).Where("id = ?", id).Updates(map[string]any{
		"title":   title,
		"content": content,
	})
	if tx.Error != nil {
		slog.Error("更新新闻失败", "error", tx.Error)
		return errors.New("更新新闻失败")
	} else {
		if tx.RowsAffected <= 0 {
			return fmt.Errorf("新闻id[%d]不存在", id)
		} else {
			return nil
		}
	}
}

func GetNewsByID(id int) *model.News {
	news := &model.News{Id: id}
	err := globalDB.Select("*").Where("delete_time is null and id = ?", id).Find(news).Error
	if err != nil {
		slog.Error("获取新闻失败", "error", err)
		return nil
	}
	news.ViewPostTime = news.PostTime.Format("2006-01-02 15:04:05")
	return news
}

func GetNewsByUid(uid int) *[]model.News {
	var news *[]model.News
	err := globalDB.Select("*").Where("delete_time is null and user_id = ?", uid).Find(news).Error
	if err != nil {
		slog.Error("获取新闻失败", "error", err)
		return nil
	}
	return news
}

func GetNewsByPage(page int, pageSize int) (int, []*model.News) {
	var total int64
	var news []*model.News
	err := globalDB.Model(&model.News{}).Where("delete_time is null").Count(&total).Error
	if err != nil {
		slog.Error("获取新闻总数失败", "error", err)
		return 0, nil
	}

	err = globalDB.Select("*").Where("delete_time is null").Order("create_time desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&news).Error
	if err != nil {
		slog.Error("GetNewsByPage failed", "pageNo", page, "pageSize", pageSize, "error", err)
		return 0, nil
	}

	if len(news) > 0 {
		for _, v := range news {
			v.ViewPostTime = v.PostTime.Format("2006-01-02 15:04:05")
		}
	}
	return int(total), news
}
