package mysql

import (
	"backbend/models"
	"go.uber.org/zap"
)



func CreateReply(reply *models.Reply) (err error) {
	sqlStr := `insert into reply(
	comment_id, content, author_id, author_name, parent_id, parent_name, create_time)
	values(?,?,?,?,?,?,?)`
	_, err = db.Exec(sqlStr, reply.CommentID, reply.Content, reply.AuthorID,
		reply.AuthorName, reply.ParentID, reply.AuthorName, reply.CreateTime)
	if err != nil {
		zap.L().Error("insert reply failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}


func GetReplyList(id int64) (replyList []*models.Reply, err error) {
	sqlStr := `select content, author_id, author_name, parent_id, parent_name, create_time
	from reply where comment_id = ?`
	
	err = db.Select(&replyList, sqlStr, id)
	if(err != nil) {
		return 
	}
	return
}
