package internal

import (
	"context"
	embed_grpc "go-server/grpc"

	"github.com/goccy/go-json"
	"google.golang.org/grpc"
)

type Service struct {
	repository *Repository
	embedding_grpc  embed_grpc.InferenceAPIsServiceClient
	logger Logger
}

func NewService(repository *Repository, conn *grpc.ClientConn, 	logger Logger) *Service {
	embedding_grpc := embed_grpc.NewInferenceAPIsServiceClient(conn)
	return &Service{repository: repository, embedding_grpc: embedding_grpc, logger: logger}
}

// generate and store embeddings
// @todo - check if text exceeds token limit before generating embedding
func (s *Service) GenerateAndStoreTextEmbeddings(ctx context.Context, text EmbeddingRequest) error {
	var text_embedding Embedding
	// generate embeddings
	s.logger.LogInfo("generating text embeddings...")
	results, err := s.PerformTextEmbedding(ctx, text.Text)
	if err != nil {
		s.logger.LogError("cannot generate text embeddings", err.Error())
		return err
	}

	embeds := results.GetPrediction()

	var embeddings [][]float32

	json.Unmarshal(embeds, &embeddings)
	text_embedding.Text = text.Text
	text_embedding.Embedding = embeddings[0]

	// store embeddings in db
	s.logger.LogInfo("storing text embeddings...")
	return s.StoreEmbeddings(ctx, text_embedding)
}
// perform text embedding
func (s *Service) PerformTextEmbedding(ctx context.Context, text string) (*embed_grpc.PredictionResponse, error) {
	x := map[string][]byte{"input": []byte(text)}
	input := &embed_grpc.PredictionsRequest{
		ModelName: "my_model",
		Input:     x,
	}

	res, err := s.embedding_grpc.Predictions(ctx, input)
	if err != nil {
		s.logger.LogError(err.Error())
		return &embed_grpc.PredictionResponse{}, err
	}

	return res, nil
}


// store embeddings in db
func (s *Service) StoreEmbeddings(ctx context.Context,embeddings Embedding) error {
	return s.repository.StoreEmbeddingsInDB(ctx, embeddings)
}

// retrive five similar embeddings from db
func (s *Service) RetrieveFiveSimilarEmbeddingService(ctx context.Context, text string) ([]Embedding, error) {
	results, err := s.PerformTextEmbedding(ctx, text)
	if err != nil {
		s.logger.LogError(err.Error())
		return []Embedding{}, err
	}

	embeds := results.GetPrediction()

	var embeddings [][]float32

	json.Unmarshal(embeds, &embeddings)

	return s.repository.RetrieveFiveSimilarEmbedding(ctx, embeddings[0])
}
