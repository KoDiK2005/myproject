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
	AuthorName string `db:"author_name"`
	Total      int    `db:"total"`
}

func toPost(r postRow) models.Post {
	return models.Post{
		ID: r.ID, Title: r.Title, Body: r.Body, UserID: r.UserID,
		Visibility: r.Visibility, AuthorName: r.AuthorName,
	}
}

func (r *PostRepo) Create(post *models.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	if post.Visibility == "" {
		post.Visibility = "public"
	}
	query := `WITH inserted AS (
	            INSERT INTO posts (title, body, user_id, visibility)
	            VALUES ($1, $2, $3, $4)
	            RETURNING id, title, body, user_id, visibility
	          )
	          SELECT inserted.id, users.name AS author_name
	          FROM inserted JOIN users ON users.id = inserted.user_id`
	return r.db.QueryRowxContext(ctx, query, post.Title, post.Body, post.UserID, post.Visibility).
		Scan(&post.ID, &post.AuthorName)
}

func (r *PostRepo) GetByID(id int) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var post models.Post
	err := r.db.GetContext(ctx, &post,
		`SELECT posts.id, posts.title, posts.body, posts.user_id, posts.visibility, users.name AS author_name
		 FROM posts JOIN users ON users.id = posts.user_id
		 WHERE posts.id = $1`, id)
	return &post, err
}

// GetByUserIDForViewer — посты юзера с учётом видимости.
// viewerID=0 → гость, видит только public.
// viewerID=ownerID → видит все свои.
// иначе → public + friends если дружат.
func (r *PostRepo) GetByUserIDForViewer(ownerID, viewerID, limit, offset int) ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var posts []models.Post
	err := r.db.SelectContext(ctx, &posts,
		`SELECT posts.id, posts.title, posts.body, posts.user_id, posts.visibility, users.name AS author_name
		 FROM posts JOIN users ON users.id = posts.user_id
		 WHERE posts.user_id = $1
		   AND (
		     posts.visibility = 'public'
		     OR $2 = $1
		     OR (posts.visibility = 'friends' AND $2 != 0 AND EXISTS (
		         SELECT 1 FROM friendships
		         WHERE (
		             (requester_id = $1 AND addressee_id = $2)
		             OR (requester_id = $2 AND addressee_id = $1)
		         ) AND status = 'accepted'
		     ))
		   )
		 ORDER BY posts.id DESC
		 LIMIT $3 OFFSET $4`,
		ownerID, viewerID, limit, offset)
	return posts, err
}

// GetAllWithCount — публичная лента (для неавторизованных)
func (r *PostRepo) GetAllWithCount(limit, offset int) ([]models.Post, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var rows []postRow
	err := r.db.SelectContext(ctx, &rows,
		`SELECT posts.id, posts.title, posts.body, posts.user_id, posts.visibility, users.name AS author_name, COUNT(*) OVER() AS total
		 FROM posts JOIN users ON users.id = posts.user_id
		 WHERE posts.visibility = 'public'
		 ORDER BY posts.id DESC
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
		`SELECT posts.id, posts.title, posts.body, posts.user_id, posts.visibility, users.name AS author_name, COUNT(*) OVER() AS total
		 FROM posts JOIN users ON users.id = posts.user_id
		 WHERE
		   -- свои посты видишь всегда
		   posts.user_id = $1
		   OR
		   -- публичные посты всех
		   posts.visibility = 'public'
		   OR
		   -- приватные посты друзей
		   (posts.visibility = 'friends' AND posts.user_id IN (
		       SELECT CASE
		           WHEN requester_id = $1 THEN addressee_id
		           ELSE requester_id
		       END
		       FROM friendships
		       WHERE (requester_id = $1 OR addressee_id = $1) AND status = 'accepted'
		   ))
		 ORDER BY posts.id DESC
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
		`SELECT posts.id, posts.title, posts.body, posts.user_id, posts.visibility, users.name AS author_name, COUNT(*) OVER() AS total
		 FROM posts JOIN users ON users.id = posts.user_id
		 WHERE posts.visibility = 'public' AND (posts.title ILIKE $1 OR posts.body ILIKE $1)
		 ORDER BY posts.id DESC
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
	query := `WITH updated AS (
	            UPDATE posts SET title=$1, body=$2, visibility=$3
	            WHERE id=$4 RETURNING id, title, body, user_id, visibility
	          )
	          SELECT updated.id, updated.title, updated.body, updated.user_id, updated.visibility, users.name AS author_name
	          FROM updated JOIN users ON users.id = updated.user_id`
	err := r.db.QueryRowxContext(ctx, query, title, body, visibility, id).
		Scan(&post.ID, &post.Title, &post.Body, &post.UserID, &post.Visibility, &post.AuthorName)
	return &post, err
}

func (r *PostRepo) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err := r.db.ExecContext(ctx, "DELETE FROM posts WHERE id = $1", id)
	return err
}

// IsFriend — проверка дружбы для проверки видимости friends-only постов
func (r *PostRepo) IsFriend(userA, userB int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var exists bool
	err := r.db.QueryRowxContext(ctx,
		`SELECT EXISTS(
		   SELECT 1 FROM friendships
		   WHERE ((requester_id = $1 AND addressee_id = $2)
		       OR (requester_id = $2 AND addressee_id = $1))
		   AND status = 'accepted'
		 )`,
		userA, userB).Scan(&exists)
	return exists, err
}

// IsBlocked — есть ли блокировка между двумя юзерами в любую сторону
func (r *PostRepo) IsBlocked(userA, userB int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var exists bool
	err := r.db.QueryRowxContext(ctx,
		`SELECT EXISTS(
		   SELECT 1 FROM blocks
		   WHERE (blocker_id = $1 AND blocked_id = $2)
		      OR (blocker_id = $2 AND blocked_id = $1)
		 )`,
		userA, userB).Scan(&exists)
	return exists, err
}
