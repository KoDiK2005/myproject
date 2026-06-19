package repository

import (
	"context"
	"myproject/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

const dbTimeout = 5 * time.Second

type PostRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) *PostRepo {
	return &PostRepo{db: db}
}

type postRow struct {
	ID         int    `db:"id"`
	Title      string `db:"title"`
	Body       string `db:"body"`
	UserID     int    `db:"user_id"`
	Visibility string `db:"visibility"`
	Total      int    `db:"total"`
}

func toPost(r postRow) models.Post {
	return models.Post{ID: r.ID, Title: r.Title, Body: r.Body, UserID: r.UserID, Visibility: r.Visibility}
}

func (r *PostRepo) Create(post *models.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	if post.Visibility == "" {
		post.Visibility = "public"
	}
	query := `INSERT INTO posts (title, body, user_id, visibility)
	          VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRowxContext(ctx, query, post.Title, post.Body, post.UserID, post.Visibility).Scan(&post.ID)
}

func (r *PostRepo) GetByID(id int) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var post models.Post
	err := r.db.GetContext(ctx, &post,
		"SELECT id, title, body, user_id, visibility FROM posts WHERE id = $1", id)
	return &post, err
}

func (r *PostRepo) GetByUserID(userID, limit, offset int) ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var posts []models.Post
	err := r.db.SelectContext(ctx, &posts,
		`SELECT id, title, body, user_id, visibility
		 FROM posts WHERE user_id = $1 ORDER BY id DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset)
	return posts, err
}

// GetAllWithCount — публичная лента (для неавторизованных)
func (r *PostRepo) GetAllWithCount(limit, offset int) ([]models.Post, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var rows []postRow
	err := r.db.SelectContext(ctx, &rows,
		`SELECT id, title, body, user_id, visibility, COUNT(*) OVER() AS total
		 FROM posts
		 WHERE visibility = 'public'
		 ORDER BY id DESC
		 LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		var total int
		_ = r.db.QueryRowxContext(ctx, "SELECT COUNT(*) FROM posts WHERE visibility = 'public'").Scan(&total)
		return nil, total, nil
	}
	posts := make([]models.Post, len(rows))
	for i, row := range rows {
		posts[i] = toPost(row)
	}
	return posts, rows[0].Total, nil
}

// GetFeedWithCount — персональная лента авторизованного юзера:
// свои посты + публичные посты всех + приватные посты друзей
func (r *PostRepo) GetFeedWithCount(userID, limit, offset int) ([]models.Post, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var rows []postRow
	err := r.db.SelectContext(ctx, &rows,
		`SELECT id, title, body, user_id, visibility, COUNT(*) OVER() AS total
		 FROM posts
		 WHERE
		   -- свои посты видишь всегда
		   user_id = $1
		   OR
		   -- публичные посты всех
		   visibility = 'public'
		   OR
		   -- приватные посты друзей
		   (visibility = 'friends' AND user_id IN (
		       SELECT CASE
		           WHEN requester_id = $1 THEN addressee_id
		           ELSE requester_id
		       END
		       FROM friendships
		       WHERE (requester_id = $1 OR addressee_id = $1) AND status = 'accepted'
		   ))
		 ORDER BY id DESC
		 LIMIT $2 OFFSET $3`,
		userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return nil, 0, nil
	}
	posts := make([]models.Post, len(rows))
	for i, row := range rows {
		posts[i] = toPost(row)
	}
	return posts, rows[0].Total, nil
}

// SearchWithCount — поиск только по публичным постам
func (r *PostRepo) SearchWithCount(query string, limit, offset int) ([]models.Post, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var rows []postRow
	like := "%" + query + "%"
	err := r.db.SelectContext(ctx, &rows,
		`SELECT id, title, body, user_id, visibility, COUNT(*) OVER() AS total
		 FROM posts
		 WHERE visibility = 'public' AND (title ILIKE $1 OR body ILIKE $1)
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
		posts[i] = toPost(row)
	}
	return posts, rows[0].Total, nil
}

func (r *PostRepo) Update(id int, title, body, visibility string) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var post models.Post
	query := `UPDATE posts SET title=$1, body=$2, visibility=$3
	          WHERE id=$4 RETURNING id, title, body, user_id, visibility`
	err := r.db.QueryRowxContext(ctx, query, title, body, visibility, id).
		Scan(&post.ID, &post.Title, &post.Body, &post.UserID, &post.Visibility)
	return &post, err
}

func (r *PostRepo) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err := r.db.ExecContext(ctx, "DELETE FROM posts WHERE id = $1", id)
	return err
}
