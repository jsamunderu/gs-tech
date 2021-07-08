package questions

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"k8s.io/klog/v2"
)

// Profile is the class for profiles
type Profile struct {
	dbHost        string
	profileClient *redis.Client
	sessionClient *redis.Client
	usersClient   *redis.Client
	profileSChema int
	sessionSchema int
	usersSchema   int
}

const sessionTimeout = time.Minute * 10

// NewProfile defines a profile instance
func NewProfile(dbHost, profileSchema, sessionSchema, usersSchema string) (*Profile, error) {
	profile := &Profile{dbHost: dbHost}

	var err error
	profile.profileSChema, err = strconv.Atoi(profileSchema)
	if err != nil {
		return nil, err
	}
	profile.sessionSchema, err = strconv.Atoi(sessionSchema)
	if err != nil {
		return nil, err
	}
	profile.usersSchema, err = strconv.Atoi(usersSchema)
	if err != nil {
		return nil, err
	}

	profile.profileClient = redis.NewClient(&redis.Options{
		Addr:     profile.dbHost,
		Password: "",                    // no password set
		DB:       profile.profileSChema, // use default DB
	})

	profile.sessionClient = redis.NewClient(&redis.Options{
		Addr:     profile.dbHost,
		Password: "",                    // no password set
		DB:       profile.sessionSchema, // use default DB
	})

	profile.usersClient = redis.NewClient(&redis.Options{
		Addr:     profile.dbHost,
		Password: "",                  // no password set
		DB:       profile.usersSchema, // use default DB
	})
	return profile, nil
}

func (p *Profile) setProfile(key, value string) bool {
	err := p.profileClient.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		klog.ErrorS(err, "Could not set profile value", "key", key, "value", value)
		return false
	}
	return true
}

func (p *Profile) getProfile(key string) (string, bool) {
	value, err := p.profileClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		klog.ErrorS(err, "Profile: Failed to read REDIS", "key", key)
		return "", false
	}
	if err == redis.Nil {
		klog.ErrorS(err, "profile Key does not exist", "key", key)
		return "", false
	}
	return value, true
}

func (p *Profile) deleteProfile(key string) bool {
	err := p.profileClient.Del(context.Background(), key).Err()
	if err != nil && err != redis.Nil {
		klog.ErrorS(err, "Profile: Failed to read REDIS", "key", key)
		return false
	}
	if err == redis.Nil {
		klog.ErrorS(err, "profile Key does not exist", "key", key)
		return false
	}
	return true
}

func (p *Profile) setSession(key, value string) bool {
	err := p.sessionClient.Set(context.Background(), key, value, sessionTimeout).Err()
	if err != nil {
		klog.ErrorS(err, "Could not set session value", "key", key, "value", value)
		return false
	}
	return true
}

func (p *Profile) getSession(key string) (string, bool) {
	value, err := p.sessionClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		klog.ErrorS(err, "Session: Failed to read REDIS", "key", key)
		return "", false
	}
	if err == redis.Nil {
		klog.ErrorS(err, "session Key does not exist", "key", key)
		return "", false
	}
	return value, true
}

func (p *Profile) setUser(key string, value []byte) {
	err := p.usersClient.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		klog.ErrorS(err, "Could not set user value", "key", key, "value", value)
		return
	}
}

func (p *Profile) getUser(key string) ([]byte, bool) {
	value, err := p.usersClient.Get(context.Background(), key).Bytes()
	if err != nil && err != redis.Nil {
		klog.ErrorS(err, "User: Failed to read REDIS", "key", key)
		return nil, false
	}
	if err == redis.Nil {
		klog.ErrorS(err, "user Key does not exist", "key", key)
		return nil, false
	}
	return value, true
}
