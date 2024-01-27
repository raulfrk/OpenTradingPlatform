package data

import (
	"tradingplatform/shared/logging"

	"gorm.io/gorm/clause"
)

type SubscribedTopic struct {
	Topic       string `gorm:"primaryKey"`
	AgentsCount int    `gorm:"default:0"`
}

func AddSubscribedTopic(topic string, agentsCount int) {
	LocalDBLock.Lock()

	tx := LocalDB.Clauses(clause.OnConflict{DoNothing: true}).Create(&SubscribedTopic{
		Topic:       topic,
		AgentsCount: agentsCount,
	})

	if tx.Error != nil {
		logging.Log().Error().
			Err(tx.Error).
			Msg("adding subscribed topic to local database")
	}

	LocalDBLock.Unlock()
}

func RemoveSubscribedTopic(topic string) {
	LocalDBLock.Lock()

	tx := LocalDB.Delete(&SubscribedTopic{}, "topic = ?", topic)

	if tx.Error != nil {
		logging.Log().Error().
			Err(tx.Error).
			Msg("removing subscribed topic from local database")
	}

	LocalDBLock.Unlock()
}

func GetSubscribedTopics() []SubscribedTopic {
	var subscribedTopics []SubscribedTopic

	LocalDBLock.Lock()

	tx := LocalDB.Find(&subscribedTopics)

	if tx.Error != nil {
		logging.Log().Error().
			Err(tx.Error).
			Msg("getting subscribed topics from local database")
	}

	LocalDBLock.Unlock()

	return subscribedTopics
}

func GetSubscribedTopic(topic string) SubscribedTopic {
	var subscribedTopic SubscribedTopic

	LocalDBLock.Lock()

	tx := LocalDB.First(&subscribedTopic, "topic = ?", topic)

	if tx.Error != nil {
		logging.Log().Error().
			Err(tx.Error).
			Msg("getting subscribed topic from local database")
	}

	LocalDBLock.Unlock()

	return subscribedTopic
}
