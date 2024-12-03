package main

import (
	"context"
	"fmt"

	"github.com/qdrant/go-client/qdrant"
)

func InitQdrantClient() (*qdrant.Client, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})
	if err != nil {
		fmt.Errorf("search error: %w", err)
	}
	return client, err
}

func CreateCollection(client *qdrant.Client, collectionName string, vectorSize int) error {
	client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(vectorSize),
			Distance: qdrant.Distance_Cosine,
		}),
	})

	return nil
}

func UpserVector(client *qdrant.Client, collectionName string, points []*qdrant.PointStruct) error {
	_, err := client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points:         points,
	})
	return err
}

func SearchSimilarVectors(client *qdrant.Client, collectionName string, vector []float32, topK int) ([]*qdrant.ScoredPoint, error) {
	searchRequest := &qdrant.SearchPoints{
		CollectionName: collectionName,
		Vector:         vector,
		Limit:          uint64(topK),
	}

	response, err := client.GetPointsClient().Search(context.Background(), searchRequest)

	if err != nil {
		return nil, fmt.Errorf("search error: %w", err)
	}

	return response.Result, nil
}

func GetChatHistory(qdrantClient *qdrant.Client, collectionName string, chatID int64, limit uint32) ([]string, error) {
	// Фильтр по chat_id
	filter := &qdrant.Filter{
		Must: []*qdrant.Condition{
			{
				ConditionOneOf: &qdrant.Condition_Field{
					Field: &qdrant.FieldCondition{
						Key: "chat_id",
						Match: &qdrant.Match{
							MatchValue: &qdrant.Match_Integer{
								Integer: chatID,
							},
						},
					},
				},
			},
		},
	}

	// Point selection request
	response, err := qdrantClient.Scroll(context.Background(), &qdrant.ScrollPoints{
		CollectionName: collectionName,
		Filter:         filter,
		Limit:          &limit,
		WithPayload:    &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
	})
	if err != nil {
		return nil, err
	}

	var messages []string
	for _, point := range response {
		message := point.Payload["message"].GetStringValue()
		messages = append(messages, message)
	}

	return messages, nil
}

func DeleteCollection(client *qdrant.Client, collectionName string) error {
	err := client.DeleteCollection(context.Background(), collectionName)
	return err
}
