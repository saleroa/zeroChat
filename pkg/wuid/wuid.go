package wuid

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"

	"github.com/edwingeng/wuid/mysql/wuid"
)

// 生成唯一的 id，每次检查 mysql 防止出现重复id

var w *wuid.WUID

func Init(dsn string) {

	newDB := func() (*sql.DB, bool, error) {
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, false, err
		}
		return db, true, nil
	}

	w = wuid.NewWUID("default", nil)
	_ = w.LoadH28FromMysql(newDB, "wuid")
}

func GenUid(dsn string) string {
	if w == nil {
		Init(dsn)
	}

	return fmt.Sprintf("%#016x", w.Next())
}

// 将两个 id 组合成新 id
func CombineId(aid, bid string) string {
	ids := []string{aid, bid}

	sort.Slice(ids, func(i, j int) bool {
		a, _ := strconv.ParseUint(ids[i], 0, 64)
		b, _ := strconv.ParseUint(ids[j], 0, 64)
		return a < b
	})

	return fmt.Sprintf("%s_%s", ids[0], ids[1])
}
