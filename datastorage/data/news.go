package data

import (
	"errors"
	"time"
	"tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NewsSymbol struct {
	Symbol          string `gorm:"primaryKey"`
	NewsFingerprint string `gorm:"primaryKey;foreignKey"`
}
type News struct {
	Id                 int64
	Author             string
	CreatedAtTimestamp time.Time `gorm:"index"`
	UpdatedAtTimestamp time.Time
	Headline           string
	Summary            string
	Content            string
	URL                string
	Symbols            []NewsSymbol `gorm:"foreignKey:NewsFingerprint"`
	Fingerprint        string       `gorm:"primaryKey"`
	Source             string
	Sentiment          []Sentiment `gorm:"foreignKey:NewsFingerprint"`
}

type Sentiment struct {
	Timestamp                time.Time `gorm:"index"`
	Sentiment                string
	SentimentAnalysisProcess string
	NewsFingerprint          string
	// Sentiment fingerprint might not result always from hashing the same struct
	// TODO: decide what to do about this
	Fingerprint  string `gorm:"primaryKey"`
	LLMName      string `gorm:"foreignKey:LLMName"` // New field for LLM name
	Symbol       string
	SystemPrompt string
	Failed       bool
	RawSentiment string
}

type LLM struct {
	Name       string      `gorm:"primaryKey"`
	Sentiments []Sentiment `gorm:"foreignKey:LLMName;references:Name"`
}

func SentimentFromEntity(entity *entities.NewsSentiment) Sentiment {
	return Sentiment{
		Timestamp:                time.Unix(entity.Timestamp, 0),
		Sentiment:                entity.Sentiment,
		SentimentAnalysisProcess: entity.SentimentAnalysisProcess,
		Fingerprint:              entity.Fingerprint,
		LLMName:                  entity.LLM,
		Symbol:                   entity.Symbol,
		SystemPrompt:             entity.SystemPrompt,
		Failed:                   entity.Failed,
		RawSentiment:             entity.RawSentiment,
	}
}

func NewsFromEntity(entity *entities.News) News {
	symbols := make([]NewsSymbol, len(entity.Symbols))
	sentiment := make([]Sentiment, len(entity.Sentiments))
	for i, s := range entity.Symbols {
		symbols[i] = NewsSymbol{Symbol: s}
	}

	for i, s := range entity.Sentiments {
		sentiment[i] = SentimentFromEntity(s)
	}

	return News{
		Id:                 entity.Id,
		Author:             entity.Author,
		CreatedAtTimestamp: time.Unix(entity.CreatedAt, 0),
		UpdatedAtTimestamp: time.Unix(entity.UpdatedAt, 0),
		Headline:           entity.Headline,
		Summary:            entity.Summary,
		Content:            entity.Content,
		URL:                entity.URL,
		Symbols:            symbols,
		Fingerprint:        entity.Fingerprint,
		Source:             entity.Source,
		Sentiment:          sentiment,
	}
}

func SentimentToEntity(sentiment Sentiment) *entities.NewsSentiment {
	return &entities.NewsSentiment{
		Timestamp:                sentiment.Timestamp.Unix(),
		Sentiment:                sentiment.Sentiment,
		SentimentAnalysisProcess: sentiment.SentimentAnalysisProcess,
		Fingerprint:              sentiment.Fingerprint,
		LLM:                      sentiment.LLMName,
		Symbol:                   sentiment.Symbol,
		SystemPrompt:             sentiment.SystemPrompt,
		Failed:                   sentiment.Failed,
		RawSentiment:             sentiment.RawSentiment,
	}
}

func NewsToEntity(news News) *entities.News {
	symbols := make([]string, len(news.Symbols))
	sentiment := make([]*entities.NewsSentiment, len(news.Sentiment))
	for i, s := range news.Symbols {
		symbols[i] = s.Symbol
	}

	for i, s := range news.Sentiment {
		sentiment[i] = SentimentToEntity(s)
	}

	return &entities.News{
		Id:          news.Id,
		Author:      news.Author,
		CreatedAt:   news.CreatedAtTimestamp.Unix(),
		UpdatedAt:   news.UpdatedAtTimestamp.Unix(),
		Headline:    news.Headline,
		Summary:     news.Summary,
		Content:     news.Content,
		URL:         news.URL,
		Symbols:     symbols,
		Fingerprint: news.Fingerprint,
		Source:      news.Source,
		Sentiments:  sentiment,
	}
}

func NewsToEntities(news []News) []*entities.News {
	entities := make([]*entities.News, len(news))
	for i, n := range news {
		entities[i] = NewsToEntity(n)
	}
	return entities
}

func NewsFromEntities(entities []*entities.News) []News {
	news := make([]News, len(entities))
	for i, entity := range entities {
		news[i] = NewsFromEntity(entity)
	}
	return news
}

func newsContainsSentiment(news News, sentiment Sentiment) bool {
	for _, s := range news.Sentiment {
		if s.Fingerprint == sentiment.Fingerprint {
			return true
		}
	}
	return false
}

func InsertNewsWithSentiment(news *entities.News) {
	dbNews := NewsFromEntity(news)

	var existingNews News
	if err := DB.Where("fingerprint = ?", dbNews.Fingerprint).Preload("Sentiment").First(&existingNews).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			InsertNews(news)
		} else {
			logging.Log().Error().
				Err(err).
				RawJSON("news", entities.GenerateJson(news)).
				Msg("finding news")
			return
		}
	}

	DB.Transaction(func(tx *gorm.DB) error {
		for _, sentiment := range dbNews.Sentiment {
			sentiment.NewsFingerprint = existingNews.Fingerprint
			if err := tx.Clauses(clause.OnConflict{
				DoNothing: true,
			}).Create(&LLM{Name: sentiment.LLMName}).Error; err != nil {
				logging.Log().Error().
					Err(err).
					RawJSON("sentiment", entities.GenerateJson(SentimentToEntity(sentiment))).
					Msg("inserting LLM")
			}

			if err := tx.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&sentiment).Error; err != nil {
				return err
			}

		}
		return nil
	})
}

func InsertNews(news *entities.News) {
	dbNews := NewsFromEntity(news)
	if err := DB.Create(&dbNews).Error; err != nil {
		logging.Log().Error().
			Err(err).
			RawJSON("news", entities.GenerateJson(news)).
			Msg("inserting news")
	}
}

func InsertBatchNewsWithSentiment(news []News) {
	for _, n := range news {
		InsertNewsWithSentiment(NewsToEntity(n))
	}
}

func InsertBatchNews(news []News) {
	logging.Log().Debug().Int("count", len(news)).Type("entity", news[0]).Msg("started inserting batch of news to db")
	tx := DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).CreateInBatches(news, 1000)

	if tx.Error != nil {
		logging.Log().Error().Err(tx.Error).Msg("failed to insert batch of entities to db")
		return
	}
	logging.Log().Debug().Int("count", len(news)).Type("entity", news[0]).Msg("finished inserting batch of news to db")
}

func GetNewsFromDataRequest(symbol string, req requests.DataRequest) ([]*entities.News, error) {
	fingerprint := req.GetFingerprint()
	if fingerprint != "" {
		return GetNewsFingerprint(fingerprint), nil
	}
	return GetNews(string(req.GetSource()),
		symbol,
		req.GetStartTime(),
		req.GetEndTime()), nil
}

func GetNewsFingerprint(fingerprint string) []*entities.News {
	var news []News

	tx := DB.Preload("Symbols").
		Preload("Sentiment").
		Preload("LLM").
		Where("fingerprint = ?", fingerprint).Find(&news)

	if tx.Error != nil {
		logging.Log().Error().
			Err(tx.Error).
			Str("fingerprint", fingerprint).
			Msg("getting news from database")
	}

	return NewsToEntities(news)
}

func GetNews(source string, symbol string, startTime int64, endTime int64) []*entities.News {
	var news []News
	tx := DB.Preload("Symbols").
		Joins("JOIN news_symbols ON news_symbols.news_fingerprint = news.fingerprint").
		Preload("Sentiment").
		Where("source = ? AND news_symbols.symbol = ? AND updated_at_timestamp >= ? AND updated_at_timestamp <= ?",
			source,
			symbol,
			time.Unix(startTime, 0),
			time.Unix(endTime, 0)).Order("updated_at_timestamp DESC").Find(&news)
	logging.Log().Debug().Int("count", len(news)).Msg("finished getting news from db")

	if tx.Error != nil {
		logging.Log().Error().
			Err(tx.Error).
			Str("source", source).
			Str("symbol", symbol).
			Int64("startTime", startTime).
			Int64("endTime", endTime).
			Msg("getting news from database")
	}

	return NewsToEntities(news)
}
