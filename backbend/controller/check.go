package controller

import (
	"backbend/models"
	"backbend/logic"
	"go.uber.org/zap"
)

func CheckPost(post *models.Post) (ok bool) {
	ok = false
	hits, title := logic.Check(post.Title)
	if len(hits) > 0 {
		zap.L().Debug("帖子标题命中违禁词", zap.Any("hits:", hits))
		ok = true
	}

	hits, content := logic.Check(post.Content)
	if len(hits) > 0 {
		zap.L().Debug("帖子内容命中违禁词", zap.Any("hits:", hits))
		ok = true
	}

	post.Title = title
	post.Content = content
	return
}


func CheckComment(comment *models.Comment) (ok bool) {
	ok = false
	hits, content := logic.Check(comment.Content)
	if len(hits) > 0 {
		zap.L().Debug("帖子内容命中违禁词", zap.Any("hits:", hits))
		ok = true
	}
	comment.Content = content
	return
}