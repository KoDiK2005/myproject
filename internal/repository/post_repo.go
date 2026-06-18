package repository

import (
	"myproject/internal/models"

	"github.com/jmoiron/sqlx"
)

type PostRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) *PostRepo {
	return &PostRepo{db: db}
}

// postRow — внутренняя структура для сканирования window function
type postRow struct {
	ID     int    `db:"id"`
	Title  string `db:"title"`
	Body   string `db:"body"`
	UserID int    `db:"user_id"`
	Total  int    `db:"total"`
}

func (r *PostRepo) Create(post *models.Post) error {
	query := "INSERT INTO posts (title, body, user_id) VALUES ($1, $2, $3) RETURNING id"
	return r.db.QueryRowx(query, post.Title, post.Body, post.UserID).Scan(&post.ID)
}

func (r *PostRepo) GetByID(id int) (*models.Post, error) {
	var post models.Post
	err := r.db.Get(&post, "SELECT id, title, body, user_id FROM posts WHERE id = $1", id)
	return &post, err
}

func (r *PostRepo) GetByUserID(userID, limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Select(&posts,
		"SELECT id, title, body, user_id FROM posts WHERE user_id = $1 ORDER BY id DESC LIMIT $2 OFFSET $3",
		userID, limit, offset)
	return posts, err
}

// GetAllWithCount возвращает посты и total за один запрос через window function
func (r *PostRepo) GetAllWithCount(limit, offset int) ([]models.Post, int, error) {
	var rows []postRow
	err := r.db.Select(&rows,
		`SELECT id, title, body, user_id, COUNT(*) OVER() AS total
		 FROM posts
		 ORDER BY id DESC
		 LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		var total int
		_ = r.db.QueryRowx("SELECT COUNT(*) FROM posts").Scan(&total)
		return nil, total, nil
	}
	posts := make([]models.Post, len(rows))
	for i, row := range rows {
		posts[i] = models.Post{ID: row.ID, Title: row.Title, Body: row.Body, UserID: row.UserID}
	}
	return posts, rows[0].Total, nil
}

// SearchWithCount ищет посты и возвращает total за один запрос
func (r *PostRepo) SearchWithCount(query string, limit, offset int) ([]models.Post, int, error) {
	var rows []postRow
	like := "%" + query + "%"
	err := r.db.Select(&rows,
		`SELECT id, title, body, user_id, COUNT(*) OVER() AS total
		 FROM posts
		 WHERE title ILIKE $1 OR body ILIKE $1
		 ORDER BY id DESC
		 LIMIT $2 OFFSET $3`,
		like, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return nil, 0, nil
	}
	posts := make([]models.Post, len(rows))
	for i, row := range rows {
		posts[i] = models.Post{ID: row.ID, Title: row.Title, Body: row.Body, UserID: row.UserID}
	}
	return posts, rows[0].Total, nil
}

func (r *PostRepo) Update(id int, title, body string) (*models.Post, error) {
	var post models.Post
	query := "UPDATE posts SET title=$1, body=$2 WHERE id=$3 RETURNING id, title, body, user_id"
	err := r.db.QueryRowx(query, title, body, id).Scan(&post.ID, &post.Title, &post.Body, &post.UserID)
	return &post, err
}

func (r *PostRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM posts WHERE id = $1", id)
	return err
}
