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
	"gorm.io/gorm"
)

func TestCategoryRepository_Create(t *testing.T) {
	gormDB, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := repository.NewCategoryRepositoryImpl(gormDB)

	cat := &entities.CategoryEntity{
		Name: "Technology",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `categories` (`name`,`created_at`,`updated_at`) VALUES (?,?,?)")).
		WithArgs(cat.Name, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	result, err := repo.Create(context.Background(), cat)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepository_GetAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		repo := repository.NewCategoryRepositoryImpl(gormDB)

		catRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(1, "Science", time.Now(), time.Now())
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` ORDER BY id asc")).
			WillReturnRows(catRows)

		records, err := repo.GetAll(context.Background())
		assert.NoError(t, err)
		assert.Len(t, records, 1)
		assert.Equal(t, "Science", records[0].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query error", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		repo := repository.NewCategoryRepositoryImpl(gormDB)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` ORDER BY id asc")).
			WillReturnError(errors.New("db error"))

		records, err := repo.GetAll(context.Background())
		assert.Error(t, err)
		assert.Nil(t, records)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCategoryRepository_GetById(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		repo := repository.NewCategoryRepositoryImpl(gormDB)

		row := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(1, "Tech", time.Now(), time.Now())
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`id` = ? ORDER BY `categories`.`id` LIMIT ?")).
			WithArgs(1, 1).
			WillReturnRows(row)

		cat, err := repo.GetById(context.Background(), 1)
		assert.NoError(t, err)
		assert.NotNil(t, cat)
		assert.Equal(t, uint(1), cat.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not found", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		repo := repository.NewCategoryRepositoryImpl(gormDB)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`id` = ? ORDER BY `categories`.`id` LIMIT ?")).
			WithArgs(99, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		cat, err := repo.GetById(context.Background(), 99)
		assert.Error(t, err)
		assert.Equal(t, "category not found", err.Error())
		assert.Nil(t, cat)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCategoryRepository_Update(t *testing.T) {
	gormDB, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := repository.NewCategoryRepositoryImpl(gormDB)

	cat := &entities.CategoryEntity{
		ID:   1,
		Name: "Updated Tech",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `categories` SET `name`=?,`created_at`=?,`updated_at`=? WHERE `id` = ?")).
		WithArgs(cat.Name, sqlmock.AnyArg(), sqlmock.AnyArg(), cat.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updated, err := repo.Update(context.Background(), cat)
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Tech", updated.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepository_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		repo := repository.NewCategoryRepositoryImpl(gormDB)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `categories` WHERE `categories`.`id` = ?")).
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

		repo := repository.NewCategoryRepositoryImpl(gormDB)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `categories` WHERE `categories`.`id` = ?")).
			WithArgs(99).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.Delete(context.Background(), 99)
		assert.Error(t, err)
		assert.Equal(t, "category not found", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
