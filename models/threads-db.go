package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

func (m *DBModel) Get(id int) (*Thread, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, title, content, author_id, author_name, upvotes, created_at, updated_at, is_solved
		      FROM threads WHERE id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var thread Thread

	err := row.Scan(
		&thread.ID,
		&thread.Title,
		&thread.Content,
		&thread.AuthorID,
		&thread.AuthorName,
		&thread.Upvotes,
		&thread.CreatedAt,
		&thread.UpdatedAt,
		&thread.IsSolved,
	)

	if err != nil {
		return nil, err
	}

	query = `select
				tc.id, tc.thread_id, tc.category_id, c.category_name
			from
				threads_categories tc
				left join categories c on (c.id = tc.category_id)
			where
				tc.thread_id = $1
	`

	rows, _ := m.DB.QueryContext(ctx, query, id)
	defer rows.Close()

	categories := make(map[int]string)
	for rows.Next() {
		var tc ThreadCategory
		err := rows.Scan(
			&tc.ID,
			&tc.ThreadID,
			&tc.CategoryID,
			&tc.Category.CategoryName,
		)

		if err != nil {
			return nil, err
		}
		categories[tc.CategoryID] = tc.Category.CategoryName
	}

	thread.ThreadCategory = categories

	return &thread, nil
}

func (m *DBModel) All(category ...int) ([]*Thread, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	where := ""
	if len(category) > 0 {
		where = fmt.Sprintf("where id in (select thread_id from threads_categories where category_id = %d)", category[0])
	}

	query := fmt.Sprintf(
		`SELECT id, title, content, author_id, author_name, upvotes, created_at, updated_at, is_solved
			FROM threads %s order by id desc`,
		where)

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []*Thread

	for rows.Next() {

		var thread Thread
		err := rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Content,
			&thread.AuthorID,
			&thread.AuthorName,
			&thread.Upvotes,
			&thread.CreatedAt,
			&thread.UpdatedAt,
			&thread.IsSolved,
		)

		if err != nil {
			return nil, err
		}

		categoryQuery := `select
							tc.id, tc.thread_id, tc.category_id, c.category_name
						from
							threads_categories tc
							left join categories c on (c.id = tc.category_id)
						where
							tc.thread_id = $1
						`

		categoryRows, _ := m.DB.QueryContext(ctx, categoryQuery, thread.ID)
		categories := make(map[int]string)
		for categoryRows.Next() {
			var tc ThreadCategory
			err := categoryRows.Scan(
				&tc.ID,
				&tc.ThreadID,
				&tc.CategoryID,
				&tc.Category.CategoryName,
			)
			if err != nil {
				return nil, err
			}

			categories[tc.CategoryID] = tc.Category.CategoryName

		}
		categoryRows.Close()

		thread.ThreadCategory = categories
		threads = append(threads, &thread)

	}
	return threads, nil
}

func (m *DBModel) YourThreads(authorID int) ([]*Thread, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    query := `SELECT id, title, content, author_id, created_at, updated_at, is_solved
              FROM threads
              WHERE author_id = $1
              ORDER BY created_at DESC`

    rows, err := m.DB.QueryContext(ctx, query, authorID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var threads []*Thread
    for rows.Next() {
        var thread Thread
        err := rows.Scan(
            &thread.ID,
            &thread.Title,
            &thread.Content,
            &thread.AuthorID,
            &thread.CreatedAt,
            &thread.UpdatedAt,
            &thread.IsSolved,
        )
        if err != nil {
            return nil, err
        }

        categoryQuery := `SELECT
                            tc.id, tc.thread_id, tc.category_id, c.category_name
                          FROM
                            threads_categories tc
                            LEFT JOIN categories c ON (c.id = tc.category_id)
                          WHERE
                            tc.thread_id = $1`

        categoryRows, _ := m.DB.QueryContext(ctx, categoryQuery, thread.ID)
        categories := make(map[int]string)
        for categoryRows.Next() {
            var tc ThreadCategory
            err := categoryRows.Scan(
                &tc.ID,
                &tc.ThreadID,
                &tc.CategoryID,
                &tc.Category.CategoryName,
            )
            if err != nil {
                return nil, err
            }

            categories[tc.CategoryID] = tc.Category.CategoryName
        }
        categoryRows.Close()

        thread.ThreadCategory = categories
        threads = append(threads, &thread)
    }

    return threads, nil
}


func (m *DBModel) CategoriesAll() ([]*Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, category_name, created_at, updated_at FROM categories order by id asc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category

	for rows.Next() {
		var category Category
		err := rows.Scan(
			&category.ID,
			&category.CategoryName,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}

	return categories, nil
}

func (m *DBModel) InsertThread(thread Thread) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var categoryExists bool
	err := m.DB.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)`, thread.CategoryID).Scan(&categoryExists)
	if err != nil {
		log.Println(err)
		return err
	}
	if !categoryExists {
		return errors.New("category_id does not exist")
	}

	stmt := `INSERT INTO threads (title, content, author_id, author_name, upvotes, created_at, updated_at, is_solved)
             VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	err = m.DB.QueryRowContext(ctx, stmt,
		thread.Title,
		thread.Content,
		thread.AuthorID,
		thread.AuthorName,
		0,
		thread.CreatedAt,
		thread.UpdatedAt,
		false,
	).Scan(&thread.ID)

	if err != nil {
		log.Println(err)
		return err
	}

	stmt = `INSERT INTO threads_categories (thread_id, category_id, created_at, updated_at)
            VALUES ($1, $2, $3, $4)`
	_, err = m.DB.ExecContext(ctx, stmt, thread.ID, thread.CategoryID, thread.CreatedAt, thread.UpdatedAt)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (m *DBModel) UpdateThread(thread Thread) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `UPDATE threads SET title = $1, content = $2, updated_at = $3 WHERE id = $4 RETURNING id`
	err := m.DB.QueryRowContext(ctx, stmt,
		thread.Title,
		thread.Content,
		thread.UpdatedAt,
		thread.ID,
	).Scan(&thread.ID)

	if err != nil {
		log.Println(err)
		return err
	}

	stmt = `UPDATE threads_categories SET category_id = $1, updated_at = $2 WHERE thread_id = $3`
	_, err = m.DB.ExecContext(ctx, stmt, thread.CategoryID, thread.UpdatedAt, thread.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *DBModel) DeleteThread(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmtCategory := `delete from threads_categories where thread_id = $1`
	_, err := m.DB.ExecContext(ctx, stmtCategory, id)
	if err != nil {
		log.Println(err)
		return err
	}

	stmt := `delete from threads where id = $1`

	_, err = m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *DBModel) ToggleSolved(id int) (*Thread, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE threads SET is_solved = NOT is_solved WHERE id = $1 RETURNING id, title, content, author_id, author_name, upvotes, created_at, updated_at, is_solved`
	row := m.DB.QueryRowContext(ctx, query, id)

	var thread Thread
	err := row.Scan(
		&thread.ID,
		&thread.Title,
		&thread.Content,
		&thread.AuthorID,
		&thread.AuthorName,
		&thread.Upvotes,
		&thread.CreatedAt,
		&thread.UpdatedAt,
		&thread.IsSolved,
	)

	if err != nil {
		return nil, err
	}

	return &thread, nil
}

func (m *DBModel) GetUserByUsername(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT user_id, username, password FROM users WHERE username = $1`
	row := m.DB.QueryRowContext(ctx, query, username)

	var user User
	err := row.Scan(&user.UserID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *DBModel) InsertUser(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING user_id`
	err := m.DB.QueryRowContext(ctx, query, user.Username, user.Password).Scan(&user.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) GetReplies(threadID int) ([]*Reply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, thread_id, content, author_id, author_name, created_at, is_answer
              FROM replies
              WHERE thread_id = $1
              ORDER BY is_answer DESC, created_at ASC`

	rows, err := m.DB.QueryContext(ctx, query, threadID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []*Reply
	for rows.Next() {
		var reply Reply
		err := rows.Scan(&reply.ID, &reply.ThreadID, &reply.Content, &reply.AuthorID, &reply.AuthorName, &reply.CreatedAt, &reply.IsAnswer)
		if err != nil {
			return nil, err
		}
		replies = append(replies, &reply)

	}
	return replies, nil
}

func (m *DBModel) InsertReply(reply Reply) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO replies (thread_id, content, author_id, author_name, created_at, is_answer) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := m.DB.ExecContext(ctx, stmt,
		reply.ThreadID,
		reply.Content,
		reply.AuthorID,
		reply.AuthorName,
		reply.CreatedAt,
		false)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func (m *DBModel) DeleteReply(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `delete from replies where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *DBModel) ToggleAnswer(id int) (*Reply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE replies SET is_answer = NOT is_answer WHERE id = $1 RETURNING id, thread_id, content, author_id, author_name, created_at, is_answer`
	row := m.DB.QueryRowContext(ctx, query, id)

	var reply Reply
	err := row.Scan(&reply.ID, &reply.ThreadID, &reply.Content, &reply.AuthorID, &reply.AuthorName, &reply.CreatedAt, &reply.IsAnswer)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}


func (m *DBModel) StarThread(userID, threadID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO starred_threads (user_id, thread_id) VALUES ($1, $2)`
	_, err := m.DB.ExecContext(ctx, stmt, userID, threadID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (m *DBModel) UnstarThread(userID, threadID int) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    query := `DELETE FROM starred_threads WHERE user_id = $1 AND thread_id = $2`
    _, err := m.DB.ExecContext(ctx, query, userID, threadID)
    return err
}

func (m *DBModel) GetStarredThreads(userID int) ([]*Thread, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

    query := `SELECT t.id, t.title, t.content, t.author_id, t.created_at, t.updated_at, t.is_solved
              FROM threads t
              JOIN starred_threads s ON t.id = s.thread_id
              WHERE s.user_id = $1
              ORDER BY s.created_at DESC`

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []*Thread
	for rows.Next() {
		var thread Thread
		err := rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Content,
			&thread.AuthorID,
			&thread.CreatedAt,
			&thread.UpdatedAt,
			&thread.IsSolved,
		)
		if err != nil {
			return nil, err
		}

		categoryQuery := `SELECT
							tc.id, tc.thread_id, tc.category_id, c.category_name
						  FROM
							threads_categories tc
							LEFT JOIN categories c ON (c.id = tc.category_id)
						  WHERE
							tc.thread_id = $1`

		categoryRows, _ := m.DB.QueryContext(ctx, categoryQuery, thread.ID)
		categories := make(map[int]string)
		for categoryRows.Next() {
			var tc ThreadCategory
			err := categoryRows.Scan(
				&tc.ID,
				&tc.ThreadID,
				&tc.CategoryID,
				&tc.Category.CategoryName,
			)
			if err != nil {
				return nil, err
			}

			categories[tc.CategoryID] = tc.Category.CategoryName
		}
		categoryRows.Close()

		thread.ThreadCategory = categories
		threads = append(threads, &thread)
	}

	return threads, nil

}