package engine

import (
	"bytes"
	"context"
	"encoding/base32"
	"encoding/gob"
	jsonLib "encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tangvis/erp/common"
	"github.com/tangvis/erp/pkg/crypto"
	"net/http"
	"strings"
	"time"

	ginSession "github.com/gin-contrib/sessions"
	rediGo "github.com/gomodule/redigo/redis"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	goRedis "github.com/redis/go-redis/v9"
	"github.com/tangvis/erp/agent/redis"
)

// Amount of time for cookies/redis keys to expire.
const (
	sessionExpire = 86400 * 30
	keyPrefix     = "session_"
)

type Store interface {
	ginSession.Store
	SessionHandler() gin.HandlerFunc
	OnlineUsers(ctx context.Context, id uint64) ([]common.UserInfo, error)
	ForeLogout(ctx context.Context, id uint64, sid string) error
}

// SessionSerializer provides an interface hook for alternative serializers
type SessionSerializer interface {
	Deserialize(d []byte, ss *sessions.Session) error
	Serialize(ss *sessions.Session) ([]byte, error)
}

// JSONSerializer encode the session map to JSON.
type JSONSerializer struct{}

// Serialize to JSON. Will err if there are unmarshalable key values
func (s JSONSerializer) Serialize(ss *sessions.Session) ([]byte, error) {
	m := make(map[string]interface{}, len(ss.Values))
	for k, v := range ss.Values {
		ks, ok := k.(string)
		if !ok {
			err := fmt.Errorf("non-string key value, cannot serialize session to JSON: %v", k)
			fmt.Printf("redistore.JSONSerializer.serialize() Error: %v", err)
			return nil, err
		}
		m[ks] = v
	}
	return jsonLib.Marshal(m)
}

// Deserialize back to map[string]interface{}
func (s JSONSerializer) Deserialize(d []byte, ss *sessions.Session) error {
	m := make(map[string]interface{})
	err := jsonLib.Unmarshal(d, &m)
	if err != nil {
		fmt.Printf("redistore.JSONSerializer.deserialize() Error: %v", err)
		return err
	}
	for k, v := range m {
		ss.Values[k] = v
	}
	return nil
}

// GobSerializer uses gob package to encode the session map
type GobSerializer struct{}

// Serialize using gob
func (s GobSerializer) Serialize(ss *sessions.Session) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(ss)
	if err == nil {
		return buf.Bytes(), nil
	}
	return nil, err
}

// Deserialize back to map[interface{}]interface{}
func (s GobSerializer) Deserialize(d []byte, ss *sessions.Session) error {
	dec := gob.NewDecoder(bytes.NewBuffer(d))
	return dec.Decode(&ss)
}

// SessionStore stores sessions in a redis backend.
type SessionStore struct {
	cli           redis.Cache
	Codecs        []securecookie.Codec
	Opts          *sessions.Options // default configuration
	DefaultMaxAge int               // default Redis TTL for a MaxAge == 0 session
	maxLength     int
	keyPrefix     string
	serializer    SessionSerializer

	sessionHandler gin.HandlerFunc
}

func (s *SessionStore) ForeLogout(ctx context.Context, id uint64, sid string) error {
	if len(sid) > 0 {
		return s.cli.Del(ctx, sid)
	}
	onlineUsers, err := s.OnlineUsers(ctx, id)
	if err != nil {
		return err
	}
	sessionIDs := make([]string, 0, len(onlineUsers))
	for _, onlineUser := range onlineUsers {
		sessionIDs = append(sessionIDs, onlineUser.SessionID)
	}
	return s.cli.Del(ctx, sessionIDs...)
}

func (s *SessionStore) Options(options ginSession.Options) {
	s.Opts = options.ToGorillaOptions()
}

// SetMaxLength sets SessionStore.maxLength if the `l` argument is greater or equal 0
// maxLength restricts the maximum length of new sessions to l.
// If l is 0 there is no limit to the size of a session, use with caution.
// The default for a new SessionStore is 4096. Redis allows for max.
// value sizes of up to 512MB (http://redis.io/topics/data-types)
// Default: 4096,
func (s *SessionStore) SetMaxLength(l int) {
	if l >= 0 {
		s.maxLength = l
	}
}

// SetKeyPrefix set the prefix
func (s *SessionStore) SetKeyPrefix(p string) {
	s.keyPrefix = p
}

// SetSerializer sets the serializer
func (s *SessionStore) SetSerializer(ss SessionSerializer) {
	s.serializer = ss
}

// SetMaxAge restricts the maximum age, in seconds, of the session record
// both in database and a browser. This is to change session storage configuration.
// If you want just to remove session use your session `s` object and change it's
// `Opts.MaxAge` to -1, as specified in
//
//	http://godoc.org/github.com/gorilla/sessions#Options
//
// Default is the one provided by this package value - `sessionExpire`.
// Set it to 0 for no restriction.
// Because we use `MaxAge` also in SecureCookie crypting algorithm you should
// use this function to change `MaxAge` value.
func (s *SessionStore) SetMaxAge(v int) {
	var c *securecookie.SecureCookie
	var ok bool
	s.Opts.MaxAge = v
	for i := range s.Codecs {
		if c, ok = s.Codecs[i].(*securecookie.SecureCookie); ok {
			c.MaxAge(v)
		} else {
			fmt.Printf("Can't change MaxAge on codec %v\n", s.Codecs[i])
		}
	}
}

func NewRedisStore(cache redis.Cache) Store {
	store := &SessionStore{
		cli:    cache,
		Codecs: securecookie.CodecsFromPairs([]byte("secret")),
		Opts: &sessions.Options{
			Path:     "/",
			MaxAge:   sessionExpire,
			HttpOnly: true,
			Secure:   true,
		},
		DefaultMaxAge: 60 * 20, // 20 minutes seems like a reasonable default
		maxLength:     4096,
		keyPrefix:     keyPrefix,
		serializer:    GobSerializer{},
	}
	store.sessionHandler = ginSession.Sessions("session_id", store)

	return store
}

// Get returns a session for the given name after adding it to the registry.
//
// See gorilla/sessions FilesystemStore.Get().
func (s *SessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

// New returns a session for the given name without adding it to the registry.
//
// See gorilla/sessions FilesystemStore.New().
func (s *SessionStore) New(r *http.Request, name string) (*sessions.Session, error) {
	var (
		err error
		ok  bool
	)
	session := sessions.NewSession(s, name)
	// make a copy
	options := *s.Opts
	session.Options = &options
	session.IsNew = true
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, s.Codecs...)
		if err == nil {
			ok, err = s.load(r.Context(), session)
			session.IsNew = !(err == nil && ok) // not new if no error and data available
		}
	}
	return session, err
}

// Save adds a single session to the response.
func (s *SessionStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// Marked for deletion.
	if session.Options.MaxAge <= 0 {
		if err := s.delete(r.Context(), session); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
	} else {
		// Build an alphanumeric key for the redis store.
		// !!!here user.SessionID is not reliable
		user := userInfo(session.Values[common.UserInfoKey], session.ID)
		if user == nil {
			return fmt.Errorf("no user info found in session: %v", session.Values)
		}
		if session.ID != "" {
			sp := strings.Split(session.ID, "_")
			if len(sp) > 0 && sp[0] != GenerateSessionID(user.ID) {
				session.ID = newSessionID(user.ID)
			}
		} else {
			session.ID = newSessionID(user.ID)
		}
		if err := s.save(r.Context(), session); err != nil {
			return err
		}
		encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, s.Codecs...)
		if err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	}
	return nil
}

// save stores the session in redis.
func (s *SessionStore) save(ctx context.Context, session *sessions.Session) error {
	b, err := s.serializer.Serialize(session)
	if err != nil {
		return err
	}
	if s.maxLength != 0 && len(b) > s.maxLength {
		return errors.New("SessionStore: the value to store is too big")
	}
	age := session.Options.MaxAge
	if age == 0 {
		age = s.DefaultMaxAge
	}
	return s.cli.SetEx(ctx, s.keyPrefix+session.ID, b, time.Duration(age)*time.Second)
}

// load reads the session from redis.
// returns true if there is a sessoin data in DB
func (s *SessionStore) load(ctx context.Context, session *sessions.Session) (bool, error) {
	data, err := s.cli.GetBytes(ctx, s.keyPrefix+session.ID)
	if err != nil {
		if errors.Is(err, goRedis.Nil) {
			return false, nil
		}
		return false, err
	}
	if data == nil {
		return false, nil // no data was associated with this key
	}
	b, err := rediGo.Bytes(data, err)
	if err != nil {
		return false, err
	}
	return true, s.serializer.Deserialize(b, session)
}

// delete removes keys from redis if MaxAge<0
func (s *SessionStore) delete(ctx context.Context, session *sessions.Session) error {
	return s.cli.Del(ctx, s.keyPrefix+session.ID)
}

func (s *SessionStore) OnlineUsers(ctx context.Context, id uint64) ([]common.UserInfo, error) {
	prefix := s.keyPrefix
	if id > 0 {
		prefix = prefix + GenerateSessionID(id)
	}
	keys, err := s.cli.Keys(ctx, prefix+"*")
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, nil
	}
	rawSessions, err := s.cli.MGet(ctx, keys...)
	if err != nil {
		return nil, err
	}
	result := make([]common.UserInfo, 0, len(rawSessions))
	for _, raw := range rawSessions {
		b, err := rediGo.Bytes(raw, nil)
		if err != nil {
			continue
		}
		var session sessions.Session
		if err = s.serializer.Deserialize(b, &session); err != nil {
			continue
		}
		user := userInfo(session.Values[common.UserInfoKey], s.keyPrefix+session.ID)
		if user == nil {
			continue
		}
		result = append(result, *user)
	}
	return result, nil
}
func (s *SessionStore) SessionHandler() gin.HandlerFunc {
	return s.sessionHandler
}

func GenerateSessionID(userID uint64) string {
	return crypto.GetMD5Hash(fmt.Sprintf("%d", userID))
}

func newSessionID(id uint64) string {
	return fmt.Sprintf("%s_%s", GenerateSessionID(id), strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "="))
}
