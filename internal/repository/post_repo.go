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
	err := r.db.Select(&posts, "SELECT id, title, body, user_id FROM posts WHERE user_id = $1 LIMIT $2 OFFSET $3", userID, limit, offset)
	return posts, err
}

func (r *PostRepo) GetAll(limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Select(&posts, "SELECT id, title, body, user_id FROM posts LIMIT $1 OFFSET $2", limit, offset)
	return posts, err
}

func (r *PostRepo) Update(id int, title, body string) (*models.Post, error) {
	var post models.Post
	query := "UPDATE posts SET title=$1, body=$2 WHERE id=$3 RETURNING id, title, body, user_id"
	err := r.db.QueryRowx(query, title, body, id).Scan(&post.ID, &post.Title, &post.Body, &post.UserID)
	return &post, err
}

func (r *PostRepo) Count() (int, error) {
	var count int
	err := r.db.QueryRowx("SELECT COUNT(*) FROM posts").Scan(&count)
	return count, err
}

// Search ищет посты по подстроке в title или body (ILIKE — регистр не важен)
func (r *PostRepo) Search(query string, limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	like := "%" + query + "%"
	err := r.db.Select(&posts,
		`SELECT id, title, body, user_id FROM posts
		 WHERE title ILIKE $1 OR body ILIKE $1
		 LIMIT $2 OFFSET $3`,
		like, limit, offset)
	return posts, err
}

func (r *PostRepo) SearchCount(query string) (int, error) {
	var count int
	like := "%" + query + "%"
	err := r.db.QueryRowx(
		`SELECT COUNT(*) FROM posts WHERE title ILIKE $1 OR body ILIKE $1`, like,
	).Scan(&count)
	return count, err
}

func (r *PostRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM posts WHERE id = $1", id)
	return err
}
