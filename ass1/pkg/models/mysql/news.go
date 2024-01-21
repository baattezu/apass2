// pkg/models/mysql/news.go

package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// NewsModel представляет модель для работы с новостями в MySQL.
type NewsModel struct {
	DB *sql.DB
}

// News представляет структуру новости.
type News struct {
	ID       int
	Title    string
	Content  string
	Tag      string
	ImageURL string
}

// InitDB открывает соединение с базой данных MySQL.
func (m *NewsModel) InitDB(dataSourceName string) error {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}
	m.DB = db
	return nil
}

// CreateNews создает новость в базе данных.
func (m *NewsModel) CreateNews(title, content, tag, imageURL string) (int, error) {
	result, err := m.DB.Exec("INSERT INTO news (title, content, tag, image_url) VALUES (?, ?, ?, ?)", title, content, tag, imageURL)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetNews возвращает список всех новостей из базы данных.
func (m *NewsModel) GetNews() ([]News, error) {
	rows, err := m.DB.Query("SELECT id, title, content, tag, image_url FROM news")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newsList []News
	for rows.Next() {
		var news News
		err := rows.Scan(&news.ID, &news.Title, &news.Content, &news.Tag, &news.ImageURL)
		if err != nil {
			return nil, err
		}
		newsList = append(newsList, news)
	}

	return newsList, nil
}
func (m *NewsModel) GetNewsByCategory(s string) ([]News, error) {
	rows, err := m.DB.Query("SELECT id, title, content, tag, image_url FROM news where tag = ?", s)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newsList []News
	for rows.Next() {
		var news News
		err := rows.Scan(&news.ID, &news.Title, &news.Content, &news.Tag, &news.ImageURL)
		if err != nil {
			return nil, err
		}
		newsList = append(newsList, news)
	}

	return newsList, nil
}

// DeleteNews удаляет новость по ее идентификатору из базы данных.
func (m *NewsModel) DeleteNews(title string) error {
	_, err := m.DB.Exec("DELETE FROM news WHERE title = ?", title)
	return err
}

func (m *NewsModel) UpdateNews(oldtitle, title, content, tag, imageURL string) error {
	_, err := m.DB.Exec("UPDATE news SET title = ?, content = ?, tag = ? where title = ?", title, content, tag, oldtitle)
	return err
}
