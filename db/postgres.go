package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pgPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	err = pgPool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return pgPool, nil
}

type PgStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(pool *pgxpool.Pool) DatabaseService {
	return &PgStorage{pool: pool}
}

func (p *PgStorage) InsertMessage(ctx context.Context, msg string) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := "INSERT INTO msg_varchar (msg) VALUES ($1)"
	_, err = tx.Exec(ctx, query, msg)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (p *PgStorage) InsertMessages(ctx context.Context, msgs []string) error {
	batch := &pgx.Batch{}
	for _, msg := range msgs {
		batch.Queue("INSERT INTO msg_varchar (msg) VALUES ($1)", msg)
	}

	br := p.pool.SendBatch(ctx, batch)
	defer br.Close()

	_, err := br.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (p *PgStorage) InsertMessageJson(ctx context.Context, msg string) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := "INSERT INTO msg_jsonb (msg) VALUES ($1::jsonb)"
	_, err = tx.Exec(ctx, query, msg)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (p *PgStorage) InsertMessagesJson(ctx context.Context, msgs []string) error {
	batch := &pgx.Batch{}
	for _, msg := range msgs {
		batch.Queue("INSERT INTO msg_jsonb (msg) VALUES ($1::jsonb)", msg)
	}
	br := p.pool.SendBatch(ctx, batch)

	defer br.Close()

	//_, err := br.Exec()

	//if err != nil {
	//	return err
	//}
	return nil
}
