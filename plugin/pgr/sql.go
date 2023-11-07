package pgr

import (
	sql "github.com/FloatTech/sqlite"
	"strconv"
	"sync"
	"time"
)

type PhigrosSQL struct {
	QQ         int64  `db:"user_qq"` // qq nums
	PhiSession string `db:"session"` // pgr session
	Time       int64  `db:"time"`    // time.
}

var (
	pgrDatabase = &sql.Sqlite{}
	pgrLocker   = sync.Mutex{}
)

func init() {
	pgrDatabase.DBPath = engine.DataFolder() + "pgrsql.db"
	err := pgrDatabase.Open(time.Hour * 24)
	if err != nil {
		panic(err)
	}
	_ = InitDataBase()
}

func FormatUserDataBase(qq int64, session string, Time int64) *PhigrosSQL {
	return &PhigrosSQL{QQ: qq, PhiSession: session, Time: Time}
}

func InitDataBase() error {
	pgrLocker.Lock()
	defer pgrLocker.Unlock()
	return pgrDatabase.Create("userinfo", &PhigrosSQL{})
}

func GetUserInfoFromDatabase(userID int64) *PhigrosSQL {
	pgrLocker.Lock()
	defer pgrLocker.Unlock()
	var infosql PhigrosSQL
	userIDStr := strconv.FormatInt(userID, 10)
	_ = pgrDatabase.Find("userinfo", &infosql, "where user_qq is "+userIDStr)
	return &infosql
}

func GetUserInfoTimeFromDatabase(userID int64) int64 {
	pgrLocker.Lock()
	defer pgrLocker.Unlock()
	var infosql PhigrosSQL
	userIDStr := strconv.FormatInt(userID, 10)
	_ = pgrDatabase.Find("userinfo", &infosql, "where user_qq is "+userIDStr)
	return infosql.Time
}

func (info *PhigrosSQL) BindUserDataBase() error {
	pgrLocker.Lock()
	defer pgrLocker.Unlock()
	return pgrDatabase.Insert("userinfo", info)
}
