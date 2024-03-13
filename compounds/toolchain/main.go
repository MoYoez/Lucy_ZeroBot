package toolchain

import (
	"github.com/FloatTech/floatbox/file"
	sql "github.com/FloatTech/sqlite"
	"strconv"
	"strings"
	"sync"
	"time"
)

type AdvancedUserList struct {
	UserID int64 `json:"user_id"`
}

var (
	dataBaseLocater = &sql.Sqlite{}
	dataBaseLocker  = sync.Mutex{}
)

func init() {
	dataBaseLocater.DBPath = file.BOTPATH + "hosted/advanced.db"
	err := dataBaseLocater.Open(time.Hour * 24)
	if err != nil {
		panic(err)
	}

}

// SplitCommandTo Split Command and Adjust To.
func SplitCommandTo(raw string, setCommandStopper int) (splitCommandLen int, splitInfo []string) {
	rawSplit := strings.SplitN(raw, " ", setCommandStopper)
	return len(rawSplit), rawSplit
}

// IsAdvancedActionUser AdvancedUserPatherHere, to ensure user can locate
func IsAdvancedActionUser(uid int64) bool {
	// check user list here.
	dataBaseLocker.Lock()
	defer dataBaseLocker.Unlock()
	var dataFinder AdvancedUserList
	dataBaseLocater.Find("userlist", &dataFinder, "Where user_id is "+strconv.FormatInt(uid, 10))
	if dataFinder.UserID == 0 {
		return false
	}
	return true
}

func RemoveUserList(uid int64) {
	dataBaseLocker.Lock()
	defer dataBaseLocker.Unlock()
	dataBaseLocater.Del("userlist", "Where user_id is "+strconv.FormatInt(uid, 10))
}

func AddtionList(userID int64) {
	dataBaseLocker.Lock()
	defer dataBaseLocker.Unlock()
	dataBaseLocater.Insert("userlist", &AdvancedUserList{UserID: userID})
}
