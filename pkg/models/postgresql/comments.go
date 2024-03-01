package mysql

import (
	"alexedwards.net/snippetbox/pkg/models"
	"database/sql"
	"fmt"
)

type CommentsModel struct {
	DB *sql.DB
}

func (m *CommentsModel) Insert(userId, newsId int, text string) error {
	stmt := `INSERT INTO comments (user_id, news_id, text) VALUES ($1,$2,$3) Returning id`
	var id int
	err := m.DB.QueryRow(stmt, userId, newsId, text).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}
func (m *NewsModel) GetComment(newsId int) ([]*models.Comments, error) {
	stmt := "Select id, user_id, news_id, text from comments where news_id = $1"
	rows, err := m.DB.Query(stmt, newsId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	commentList := []*models.Comments{}
	for rows.Next() {
		k := &models.Comments{}
		err := rows.Scan(&k.ID, &k.UserId, &k.NewsId, &k.Text)
		if err != nil {
			return nil, err
		}
		commentList = append(commentList, k)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	fmt.Println(commentList)
	return commentList, nil
}

func (m *CommentsModel) Delete(commentID int) error {
	stmt := "DELETE FROM comments WHERE id = $1"
	_, err := m.DB.Exec(stmt, commentID)
	if err != nil {
		return err
	}
	return nil
}

func (m *CommentsModel) GetNewsId(commentId int) (int, error) {
	stmt := `SELECT news_id FROM comments WHERE id = $1`
	var newsId int
	err := m.DB.QueryRow(stmt, commentId).Scan(&newsId)
	if err != nil {
		return 0, err
	}
	return newsId, nil
}

func (m *CommentsModel) GetAuthorId(commentId int) (int, error) {
	stmt := `SELECT user_id FROM comments WHERE id = $1`
	var userId int
	err := m.DB.QueryRow(stmt, commentId).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
