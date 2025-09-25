package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseRepository[T any] interface {
	GetAll(ctx context.Context, offset, limit int, modifier func(*gorm.DB) *gorm.DB) ([]T, int64, error)
	GetByID(ctx context.Context, id uint, modifier func(*gorm.DB) *gorm.DB) (*T, error)
	GetByIDs(ctx context.Context, ids []uint, modifier func(*gorm.DB) *gorm.DB) ([]T, error)

	CreateOne(ctx context.Context, entity *T, modifier func(*gorm.DB) *gorm.DB) error
	CreateMany(ctx context.Context, entities []*T, modifier func(*gorm.DB) *gorm.DB) error

	UpdateOne(ctx context.Context, id uint, entity *T, modifier func(*gorm.DB) *gorm.DB) error
	UpdateMany(ctx context.Context, entities []*T, modifier func(*gorm.DB) *gorm.DB) error
	PatchOne(ctx context.Context, id uint, updates map[string]any, modifier func(*gorm.DB) *gorm.DB) error

	DeleteOne(ctx context.Context, id uint) error
	DeleteMany(ctx context.Context, modifier func(*gorm.DB) *gorm.DB) error

	Upsert(ctx context.Context, entity *T, conflictColumns []clause.Column, modifier func(*gorm.DB) *gorm.DB) error

	WithTx(tx *gorm.DB) BaseRepository[T]
	DB() *gorm.DB
}

type BaseRepositoryImpl[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepositoryImpl[T] {
	return &BaseRepositoryImpl[T]{db: db}
}

func (r *BaseRepositoryImpl[T]) GetAll(
	ctx context.Context,
	offset, limit int,
	modifier func(*gorm.DB) *gorm.DB,
) ([]T, int64, error) {
	var entities []T
	var total int64

	q := r.db.WithContext(ctx).Model(new(T))
	if modifier != nil {
		q = modifier(q)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(offset).Limit(limit).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *BaseRepositoryImpl[T]) GetByID(
	ctx context.Context,
	id uint,
	modifier func(*gorm.DB) *gorm.DB,
) (*T, error) {
	entity := new(T)
	q := r.db.WithContext(ctx)
	if modifier != nil {
		q = modifier(q)
	}
	if err := q.First(entity, id).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *BaseRepositoryImpl[T]) GetByIDs(
	ctx context.Context,
	ids []uint,
	modifier func(*gorm.DB) *gorm.DB,
) ([]T, error) {
	var entities []T
	q := r.db.WithContext(ctx).Model(new(T))
	if modifier != nil {
		q = modifier(q)
	}
	if err := q.Where("id IN ?", ids).Find(&entities).Error; err != nil {
		return nil, err
	}

	if len(entities) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return entities, nil
}

// ---- CREATE ----
func (r *BaseRepositoryImpl[T]) CreateOne(
	ctx context.Context,
	entity *T,
	modifier func(*gorm.DB) *gorm.DB,
) error {
	q := r.db.WithContext(ctx)
	if modifier != nil {
		q = modifier(q)
	}
	return q.Create(entity).Error
}

func (r *BaseRepositoryImpl[T]) CreateMany(
	ctx context.Context,
	entities []*T,
	modifier func(*gorm.DB) *gorm.DB,
) error {
	q := r.db.WithContext(ctx)
	if modifier != nil {
		q = modifier(q)
	}
	return q.Create(&entities).Error
}

// ---- UPDATE ----
func (r *BaseRepositoryImpl[T]) UpdateOne(
	ctx context.Context,
	id uint,
	entity *T,
	modifier func(*gorm.DB) *gorm.DB,
) error {
	q := r.db.WithContext(ctx).Model(new(T)).Where("id = ?", id)
	if modifier != nil {
		q = modifier(q)
	}

	result := q.Updates(entity)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *BaseRepositoryImpl[T]) UpdateMany(
	ctx context.Context,
	entities []*T,
	modifier func(*gorm.DB) *gorm.DB,
) error {
	q := r.db.WithContext(ctx)
	if modifier != nil {
		q = modifier(q)
	}

	result := q.Save(&entities)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *BaseRepositoryImpl[T]) PatchOne(
	ctx context.Context,
	id uint,
	updates map[string]any,
	modifier func(*gorm.DB) *gorm.DB,
) error {
	q := r.db.WithContext(ctx).Model(new(T)).Where("id = ?", id)
	if modifier != nil {
		q = modifier(q)
	}

	result := q.Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// ---- DELETE ----
func (r *BaseRepositoryImpl[T]) DeleteOne(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(new(T), id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *BaseRepositoryImpl[T]) DeleteMany(ctx context.Context, modifier func(*gorm.DB) *gorm.DB) error {
	q := r.db.WithContext(ctx).Model(new(T))
	if modifier != nil {
		q = modifier(q)
	}

	result := q.Delete(new(T))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// ---- UPSERT ----
func (r *BaseRepositoryImpl[T]) Upsert(
	ctx context.Context,
	entity *T,
	conflictColumns []clause.Column,
	modifier func(*gorm.DB) *gorm.DB,
) error {
	q := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   conflictColumns,
		UpdateAll: true,
	})
	if modifier != nil {
		q = modifier(q)
	}
	return q.Create(entity).Error
}

func (r *BaseRepositoryImpl[T]) WithTx(tx *gorm.DB) BaseRepository[T] {
	return &BaseRepositoryImpl[T]{db: tx}
}

func (r *BaseRepositoryImpl[T]) DB() *gorm.DB {
	return r.db
}
