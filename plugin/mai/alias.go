package mai

import (
	"encoding/json"
	"github.com/FloatTech/floatbox/web"
	"os"
)

type AliasesReturnValue struct {
	Aliases []struct {
		DvId     int      `json:"dv_id"`
		SongName string   `json:"song_name"`
		SongId   []int    `json:"song_id"`
		Aliases  []string `json:"aliases"`
	} `json:"aliases"`
}

// only support LXNS because => if DivingFish then need Token.

// QueryReferSong use LocalStorageData.
func QueryReferSong(Alias string, isLxnet bool) (status bool, id []int, needAcc bool) {
	// unpackedData
	getData, err := os.ReadFile(engine.DataFolder() + "alias.json")
	if err != nil {
		panic(err)
	}
	var DataHandler AliasesReturnValue
	json.Unmarshal(getData, &DataHandler)
	var onloadList [][]int
	for _, dataSearcher := range DataHandler.Aliases {
		for _, aliasSearcher := range dataSearcher.Aliases {
			if aliasSearcher == Alias {
				onloadList = append(onloadList, dataSearcher.SongId) // write in memory
			}
		}
	}
	// if list is 2,query them is from the same song? | if above 2(3 or more ,means this song need acc.)
	switch {
	case len(onloadList) == 1: // only one query.
		if isLxnet {
			for _, listhere := range onloadList[0] {
				if listhere < 10000 {
					return true, []int{listhere}, false
				}
			}
		} else {
			return true, onloadList[0], false
		}
	// query length is 2,it means this maybe same name but diff id ==> (E.G: Oshama Scramble!)
	case len(onloadList) == 2:
		for _, listHere := range onloadList[0] {
			for _, listAliasTwo := range onloadList[1] {
				if listHere == listAliasTwo {
					// same list here.
					var returnIntList []int
					returnIntList = append(returnIntList, onloadList[0]...)
					returnIntList = append(returnIntList, onloadList[1]...)
					returnIntList = removeIntDuplicates(returnIntList)
					if isLxnet {
						for _, listhere := range returnIntList {
							if listhere < 10000 {
								return true, []int{listhere}, false
							}
						}
					} else {
						return true, returnIntList, false
					}
				}
			}
		}
		// if query is none, means it need moreacc
		return true, nil, true
	case len(onloadList) >= 3:
		return true, nil, true
	}
	// no found.
	return false, nil, false
}

// UpdateAliasPackage Use simple action to update alias.
func UpdateAliasPackage() {
	getData, err := web.GetData("https://maihook.lemonkoi.one/api/alias")
	if err != nil {
		panic(err)
	}
	os.WriteFile(engine.DataFolder()+"alias.json", getData, 0777)
}

func removeIntDuplicates(list []int) []int {
	seen := make(map[int]bool)
	var result []int
	for _, item := range list {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}
