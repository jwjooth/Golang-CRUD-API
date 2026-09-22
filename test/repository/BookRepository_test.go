package repository_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"golang-restful-api/entities"
	"golang-restful-api/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestBookRepository_Create(t *testing.T) {
	gormDB, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := repository.NewBookRepositoryImpl(gormDB)

	catID := uint(2)
	book := &entities.BookEntity{
		Title:      "Clean Code",
		CategoryId: &catID,
		Author:     "Robert C. Martin",
		Stock:      10,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `books` (`category_id`,`title`,`author`,`stock`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?)")).
		WithArgs(book.CategoryId, book.Title, book.Author, book.Stock, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	result, err := repo.Create(context.Background(), book)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepository_GetAll(t *testing.T) {
	t.Run("Success with records and count", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		repo := repository.NewBookRepositoryImpl(gormDB)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `books`")).
			WillReturnRows(countRows)

		catID := uint(1)
		bookRows := sqlmock.NewRows([]string{"id", "category_id", "title", "author", "stock", "created_at", "updated_at"}).
			AddRow(1, catID, "Book 1", "Author 1", 5, time.Now(), time.Now())
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `books` ORDER BY id asc LIMIT ?")).
			WithArgs(10).
			WillReturnRows(bookRows)

		records, total, err := repo.GetAll(context.Background(), 10, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, records, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Count query fails", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		repo := repository.NewBookRepositoryImpl(gormDB)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `books`")).
			WillReturnError(errors.New("db error"))

		records, total, err := repo.GetAll(context.Background(), 10, 0)
		assert.Error(t, err)
		assert.Equal(t, int64(0), total)
		assert.Nil(t, records)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBookRepository_GetById(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		repo := repository.NewBookRepositoryImpl(gormDB)

		catID := uint(1)
		row := sqlmock.NewRows([]string{"id", "category_id", "title", "author", "stock", "created_at", "updated_at"}).
			AddRow(1, catID, "Book 1", "Author 1", 5, time.Now(), time.Now())
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `books` WHERE `books`.`id` = ? ORDER BY `books`.`id` LIMIT ?")).
			WithArgs(1, 1).
			WillReturnRows(row)

		book, err := repo.GetById(context.Background(), 1)
		assert.NoError(t, err)
		assert.NotNil(t, book)
		assert.Equal(t, uint(1), book.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		repo := repository.NewBookRepositoryImpl(gormDB)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `books` WHERE `books`.`id` = ? ORDER BY `books`.`id` LIMIT ?")).
			WithArgs(99, 1).
			WillReturnError(errors.New("not found error"))

		book, err := repo.GetById(context.Background(), 99)
		assert.Error(t, err)
		assert.Nil(t, book)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBookRepository_Update(t *testing.T) {
	gormDB, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := repository.NewBookRepositoryImpl(gormDB)

	catID := uint(1)
	book := &entities.BookEntity{
		ID:         1,
		Title:      "Clean Code 2nd",
		CategoryId: &catID,
		Author:     "Robert C. Martin",
		Stock:      12,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `books` SET `category_id`=?,`title`=?,`author`=?,`stock`=?,`created_at`=?,`updated_at`=? WHERE `id` = ?")).
		WithArgs(book.CategoryId, book.Title, book.Author, book.Stock, sqlmock.AnyArg(), sqlmock.AnyArg(), book.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updated, err := repo.Update(context.Background(), book)
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Clean Code 2nd", updated.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepository_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		repo := repository.NewBookRepositoryImpl(gormDB)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `books` WHERE `books`.`id` = ?")).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Delete(context.Background(), 1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not found (RowsAffected = 0)", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		repo := repository.NewBookRepositoryImpl(gormDB)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `books` WHERE `books`.`id` = ?")).
			WithArgs(99).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.Delete(context.Background(), 99)
		assert.Error(t, err)
		assert.ErrorIs(t, err, repository.ErrBookNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
