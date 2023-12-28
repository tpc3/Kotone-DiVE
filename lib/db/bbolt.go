package db

import (
	"Kotone-DiVE/lib/config"
	"encoding/json"

	"go.etcd.io/bbolt"
)

var db *bbolt.DB

const (
	bucketGuild = "guild"
	bucketUser  = "user"
	Bbolt       = "bbolt"
)

func LoadBbolt() error {
	var err error
	db, err = bbolt.Open(config.CurrentConfig.Db.Path, 0600, nil)
	if err != nil {
		return err
	}
	return nil
}

func CloseBbolt() error {
	err := db.Close()
	if err != nil {
		return err
	}
	return nil
}

func LoadGuildBbolt(id string) (*config.Guild, error) {
	guild := config.Guild{}
	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketGuild))
		if bucket != nil {
			value := bucket.Get([]byte(id))
			if value != nil {
				err := json.Unmarshal(value, &guild)
				if err != nil {
					return err
				}
			} else {
				guild = config.CurrentConfig.Guild
			}
		} else {
			guild = config.CurrentConfig.Guild
		}
		return nil
	})
	if err != nil {
		return &guild, err
	}
	return &guild, nil
}

func SaveGuildBbolt(id string, guild *config.Guild) error {
	jsonGuild, err := json.Marshal(guild)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketGuild))
		if err != nil {
			return err
		}
		bucket.Put([]byte(id), jsonGuild)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func LoadUserBbolt(id string) (*config.User, error) {
	exists := true
	user := config.User{}
	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketUser))
		if bucket != nil {
			value := bucket.Get([]byte(id))
			if value != nil {
				err := json.Unmarshal(value, &user)
				if err != nil {
					return err
				}
			} else {
				exists = false
			}
		} else {
			exists = false
		}
		return nil
	})
	if err != nil {
		return &user, err
	}
	if exists {
		return &user, nil
	}
	return nil, nil
}

func SaveUserBbolt(id string, user *config.User) error {
	jsonUser, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketUser))
		if err != nil {
			return err
		}
		bucket.Put([]byte(id), jsonUser)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserBbolt(id string) error {
	err := db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketUser))
		if err != nil {
			return err
		}
		return bucket.Delete([]byte(id))
	})
	if err != nil {
		return err
	}
	return nil
}
