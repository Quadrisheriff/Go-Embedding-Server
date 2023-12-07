package internal

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"github.com/pkg/errors"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// store embeddings in database
func (r *Repository) StoreEmbeddingsInDB(ctx context.Context, embedding Embedding) error {
	stmnt := "insert into embeddings (id, text, created_at, embedding) values ($1, $2, $3, $4)"

	_, err := r.db.ExecContext(ctx, stmnt, uuid.NewString(), embedding.Text, time.Now(), pgvector.NewVector(embedding.Embedding))
	if err != nil {
		return errors.Wrap(err, "cannot store embeddings in db currently")
	}

	return nil
}

// retrieve top 5 most similar embedding from database
func (r *Repository) RetrieveFiveSimilarEmbedding(ctx context.Context, embedding []float32) ([]Embedding, error) {
	stmnt := "select id, text, created_at, embedding from content_embeddings ORDER BY embedding <-> $1 LIMIT 5"
	rows, err := r.db.QueryContext(ctx, stmnt, pgvector.NewVector(embedding))
	if err != nil {
		return []Embedding{}, errors.Wrap(err, "cannot retrieve embeddings from db at the moment")
	}
	defer rows.Close()

	var embeds []Embedding

	for rows.Next() {
		var embed Embedding

		err = rows.Scan(&embed.ID, &embed.Text, &embed.CreatedAt, &embed.Embedding)
		if err != nil {
			return []Embedding{}, errors.Wrap(err, "cannot retrieve embeddings from db at the moment")
		}

		embeds = append(embeds, embed)
	}

	return embeds, nil
}
