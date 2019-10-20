package handle

import (
	"itflow/bug/asset"
	"itflow/bug/bugconfig"
	"errors"
	"itflow/gadb"
	"itflow/gaencrypt"
	"github.com/hyahm/golog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func headers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/x-www-form-urlencoded,application/json; charset=UTF-8")

	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Access-Token,X-Token")

}

var NotFoundToken = errors.New("not found token")

func logtokenmysql(r *http.Request) (*gadb.Db, string, error) {
	mc := gadb.NewSqlConfig()

	mdb, err := mc.ConnDB()
	if err != nil {
		return mdb, "", err
	}
	a := r.Header.Get("X-Token")
	destoken, err := gaencrypt.RsaDecrypt(a, bugconfig.PrivateKey, true)
	if err != nil {
		golog.Error(err.Error())
		return mdb, "", NotFoundToken
	}
	nickname, err := asset.Getvalue(string(destoken))
	if err != nil {
		golog.Error(err.Error())
		return mdb, "", NotFoundToken

	}
	err = asset.Settimeout(a)
	if err != nil {
		return mdb, "", err
	}
	return mdb, string(nickname), nil
}

func sortpermlist(permlist []string) []string {
	l := len(bugconfig.CacheSidStatus)

	newlist := make([]string, 0)
	//["aaaa", "cccc", "dddd"]
	//["aaaa", "bbbb", "cccc", "dddd"]
	for i := 0; i < l; i++ {
		for _, v := range permlist {
			if bugconfig.CacheSidStatus[int64(i)] == v {
				newlist = append(newlist, v)
			}
		}
	}
	return newlist
}

// []string{ "admin(admin)" }  前端传过来的数据换成数组插入到data
func formatUserlistToData(userlist []string, bid int) (string, []string, [][]interface{}) {
	ul := ""
	l := len(userlist)
	nicknamelist := make([]string, 0)
	args := make([][]interface{}, 0)
	for j, v := range userlist {

		onearg := make([]interface{}, 0)
		i := strings.Index(v, "(")

		nickname := v[:i]
		if j == l-1 {
			ul = ul + strconv.FormatInt(bugconfig.CacheNickNameUid[v[:i]], 10)
		} else {
			ul = ul + strconv.FormatInt(bugconfig.CacheNickNameUid[v[:i]], 10) + ","
		}
		onearg = append(onearg, bid)
		onearg = append(onearg, bugconfig.CacheNickNameUid[v[:i]])
		args = append(args, onearg)
		nicknamelist = append(nicknamelist, nickname)
	}

	return ul, nicknamelist, args
}

// bugs表中的spuser返回昵称和真实用户名的组合（直接显示在前端）
func formatUserlistToShow(userlist string) []string {
	al := make([]string, 0)
	ul := strings.Split(userlist, ",")
	for _, v := range ul {
		uid, _ := strconv.Atoi(v)
		al = append(al, bugconfig.CacheUidNickName[int64(uid)]+"("+bugconfig.CacheUidRealName[int64(uid)]+")")
	}
	return al
}

// bugs表中的spuser返回真实用户名
func formatUserlistToRealname(userlist string) []string {
	al := make([]string, 0)
	ul := strings.Split(userlist, ",")
	for _, v := range ul {
		uid, _ := strconv.Atoi(v)
		al = append(al, bugconfig.CacheUidRealName[int64(uid)])
	}
	return al
}

// 插入到log表中
func insertlog(conn *gadb.Db, classify string, content string, r *http.Request) error {
	logsql := "insert into log(exectime,classify,content,ip) values(?,?,?,?)"
	ip := strings.Split(r.RemoteAddr, ":")[0]
	if ip != "127.0.0.1" {
		_, err := conn.Insert(logsql, time.Now().Unix(), classify, content, ip)
		if err != nil {
			return err
		}
	}

	return nil
}