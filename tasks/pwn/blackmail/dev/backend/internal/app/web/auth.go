// Auth middleware
package web

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"cbs.dev/brics/droidchat/internal/app"
	"github.com/gin-gonic/gin"
)

const AuthCtxKey = "droidchat.auth"

var Key []byte = []byte("dr01dch4t-tok3N-kEy @_@")

func AppendKey(key []byte) {
	Key = append(Key, key...)
}

func getAuth(c *gin.Context) (userId app.Uid, err error) {
	if token, ok := c.GetQuery("token"); !ok || token == "" {
		return -1, errors.New("'token' query missing/empty")
	} else {
		return parseToken(token)
	}
}

func GetUid(c *gin.Context) app.Uid {
	uid, ok := c.Get(AuthCtxKey)
	if !ok {
		panic(errors.New("must include auth middleware to use!"))
	}
	return uid.(app.Uid)
}

func RequireAuth(c *gin.Context) {
	uid, err := getAuth(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.Set(AuthCtxKey, uid)
	c.Next()
}

func MakeToken(u app.Uid) string {
	message := make([]byte, 8)
	binary.LittleEndian.PutUint64(message, uint64(u))

	mac := hmac.New(sha256.New, Key)
	mac.Write(message)
	h := hex.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf("%v:%064s", u, h)
}

func parseToken(token string) (userId app.Uid, err error) {
	parts := strings.Split(token, ":")
	fmt.Printf("%q\n", token)
	if len(parts) != 2 {
		return -1, errors.New("invalid token")
	}

	var uid int
	if parsedUid, err := strconv.ParseUint(parts[0], 10, 31); err != nil || parsedUid < 1 {
		return -1, errors.New("invalid id")
	} else {
		uid = int(parsedUid)
	}

	mac, err := hex.DecodeString(parts[1])
	if err != nil {
		return -1, errors.New("invalid token hex")
	}

	if !verifyHmac(uid, mac) {
		return -1, errors.New("invalid token mac")
	}

	return app.Uid(uid), nil
}

func verifyHmac(userId int, tokenMac []byte) bool {
	message := make([]byte, 8)
	binary.LittleEndian.PutUint64(message, uint64(userId))

	mac := hmac.New(sha256.New, Key)
	mac.Write(message)
	return hmac.Equal(tokenMac, mac.Sum(nil))
}
