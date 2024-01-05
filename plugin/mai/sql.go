package mai

import (
	"strconv"
	"sync"
	"time"

	sql "github.com/FloatTech/sqlite"
)

type UserIDToQQ struct {
	QQ     int64  `db:"user_qq"` // qq nums
	Userid string `db:"user_id"` // user_id
}

type DataHostSQL struct {
	QQ         int64  `db:"user_qq"` // qq nums
	Plate      string `db:"plate"`   // plate
	Background string `db:"bg"`      // bg
}

type UserSwitcherService struct {
	QQ     int64 `db:"user_qq"`
	IsUsed bool  `db:"isused"` // true == lxns service \ false == Diving Fish.
}

var (
	maiDatabase = &sql.Sqlite{}
	maiLocker   = sync.Mutex{}
)

func init() {
	maiDatabase.DBPath = engine.DataFolder() + "maisql.db"
	err := maiDatabase.Open(time.Hour * 24)
	if err != nil {
		panic(err)
	}
	InitDataBase()
}

func FormatUserDataBase(qq int64, plate string, bg string) *DataHostSQL {
	return &DataHostSQL{QQ: qq, Plate: plate, Background: bg}
}

func FormatUserSwitcher(qq int64, isSwitcher bool) *UserSwitcherService {
	return &UserSwitcherService{QQ: qq, IsUsed: isSwitcher}
}

func FormatUserIDDatabase(qq int64, userid string) *UserIDToQQ {
	return &UserIDToQQ{QQ: qq, Userid: userid}
}

func InitDataBase() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	maiDatabase.Create("userinfo", &DataHostSQL{})
	maiDatabase.Create("useridinfo", &UserIDToQQ{})
	maiDatabase.Create("userswitcherinfo", &UserSwitcherService{})
	return nil
}

func GetUserSwitcherInfoFromDatabase(userid int64) bool {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var info UserSwitcherService
	userIDStr := strconv.FormatInt(userid, 10)
	err := maiDatabase.Find("userswitcherinfo", &info, "where user_qq is "+userIDStr)
	if err != nil {
		return false
	}
	return info.IsUsed
}

func (info *UserSwitcherService) ChangeUserSwitchInfoFromDataBase() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	return maiDatabase.Insert("userswitcherinfo", info)
}

// GetUserIDFromDatabase Params: user qq id ==> user maimai id.
func GetUserIDFromDatabase(userID int64) UserIDToQQ {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var infosql UserIDToQQ
	userIDStr := strconv.FormatInt(userID, 10)
	_ = maiDatabase.Find("useridinfo", &infosql, "where user_qq is "+userIDStr)
	return infosql
}

func RemoveUserIdFromDatabase(qq int64) {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	userIDStr := strconv.FormatInt(qq, 10)
	_ = maiDatabase.Del("useridinfo", "where user_qq is "+userIDStr)
}

func (info *UserIDToQQ) BindUserIDDataBase() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	return maiDatabase.Insert("useridinfo", info)
}

// maimai render b50

func GetUserInfoFromDatabase(userID int64) string {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var infosql DataHostSQL
	userIDStr := strconv.FormatInt(userID, 10)
	_ = maiDatabase.Find("userinfo", &infosql, "where user_qq is "+userIDStr)
	return infosql.Plate
}

func GetUserDefaultinfoFromDatabase(userID int64) string {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var infosql DataHostSQL
	userIDStr := strconv.FormatInt(userID, 10)
	_ = maiDatabase.Find("userinfo", &infosql, "where user_qq is "+userIDStr)
	return infosql.Background
}

func (info *DataHostSQL) BindUserDataBase() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	return maiDatabase.Insert("userinfo", info)
}
