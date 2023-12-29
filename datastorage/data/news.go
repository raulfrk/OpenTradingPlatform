package data

import (
	"time"
	"tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"

	"gorm.io/gorm/clause"
)

type NewsSymbol struct {
	Symbol          string `gorm:"primaryKey"`
	NewsFingerprint string `gorm:"primaryKey"`
}
type News struct {
	Id          int64
	Author      string
	CreatedAt   time.Time `gorm:"index"`
	UpdatedAt   time.Time
	Headline    string
	Summary     string
	Content     string
	URL         string
	Symbols     []NewsSymbol `gorm:"foreignKey:NewsFingerprint"`
	Fingerprint string       `gorm:"primaryKey"`
	Source      string
	Sentiment   []Sentiment `gorm:"foreignKey:NewsFingerprint"`
}

type Sentiment struct {
	Timestamp                time.Time `gorm:"index"`
	Sentiment                string
	SentimentAnalysisProcess string
	NewsFingerprint          string
	LLM                      []LLM  `gorm:"many2many:sentiment_llm"`
	Fingerprint              string `gorm:"primaryKey"`
}

type LLM struct {
	Sentiment []Sentiment `gorm:"many2many:sentiment_llm"`
	Name      string      `gorm:"primaryKey;"`
}

func SentimentFromEntity(entity *entities.NewsSentiment) Sentiment {
	llm := make([]LLM, len(entity.LLM))

	for i, l := range entity.LLM {
		llm[i] = LLM{
			Name: l,
		}
	}
	return Sentiment{
		Timestamp:                time.Unix(entity.Timestamp, 0),
		Sentiment:                entity.Sentiment,
		SentimentAnalysisProcess: entity.SentimentAnalysisProcess,
		Fingerprint:              entity.Fingerprint,
		LLM:                      llm,
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
		Id:          entity.Id,
		Author:      entity.Author,
		CreatedAt:   time.Unix(entity.CreatedAt, 0),
		UpdatedAt:   time.Unix(entity.UpdatedAt, 0),
		Headline:    entity.Headline,
		Summary:     entity.Summary,
		Content:     entity.Content,
		URL:         entity.URL,
		Symbols:     symbols,
		Fingerprint: entity.Fingerprint,
		Source:      entity.Source,
		Sentiment:   sentiment,
	}
}

func SentimentToEntity(sentiment Sentiment) *entities.NewsSentiment {
	llm := make([]string, len(sentiment.LLM))

	for i, l := range sentiment.LLM {
		llm[i] = l.Name
	}
	return &entities.NewsSentiment{
		Timestamp:                sentiment.Timestamp.Unix(),
		Sentiment:                sentiment.Sentiment,
		SentimentAnalysisProcess: sentiment.SentimentAnalysisProcess,
		Fingerprint:              sentiment.Fingerprint,
		LLM:                      llm,
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
		CreatedAt:   news.CreatedAt.Unix(),
		UpdatedAt:   news.UpdatedAt.Unix(),
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

func InsertNews(news *entities.News) {
	dbNews := NewsFromEntity(news)
	tx := DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&dbNews)

	if tx.Error != nil {
		logging.Log().Error().
			Err(tx.Error).
			RawJSON("news", entities.GenerateJson(news)).
			Msg("inserting news")
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

func GetNewsFromRequest(symbol string, req requests.DataRequest) ([]*entities.News, error) {
	return GetNews(string(req.GetSource()),
		symbol,
		req.GetStartTime(),
		req.GetEndTime()), nil
}

func GetNews(source string, symbol string, startTime int64, endTime int64) []*entities.News {
	var news []News

	tx := DB.Preload("Symbols").
		Preload("Sentiment").
		Preload("LLM").
		Where("source = ? AND updated_at >= ? AND updated_at <= ?",
			source,
			time.Unix(startTime, 0),
			time.Unix(endTime, 0)).Order("updated_at DESC").Find(&news)

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
