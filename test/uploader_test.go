package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/polaris1119/logger"
	"github.com/qiniu/api.v6/conf"
	"github.com/qiniu/api.v6/io"
	"github.com/qiniu/api.v6/rs"
	gio "io"
	"io/ioutil"
	"sync"
	"testing"
	"time"
)

type Uploader struct {
	bucketName string

	uptoken   string
	tokenTime time.Time
	locker    sync.RWMutex
}

func (this *Uploader) InitQiniu() {
	logger.Init("D:\\goproject\\studygolang\\log", "DEBUG")
	conf.ACCESS_KEY = "lCcyC4KgEa58sTjC4sllbw4lfSOuTP62lFVegwZQ"
	conf.SECRET_KEY = "06F0Pxu8EeSRbMyXYXD6JrjYOFEAXZ5N9GtaK3Wy"
	conf.UP_HOST = "http://cdn.studytrade.club/"
	this.bucketName = "Blog"
}
func (this *Uploader) genUpToken() {
	// 避免服务器时间不同步，45分钟就更新 token
	if this.uptoken != "" && this.tokenTime.Add(45*time.Minute).After(time.Now()) {
		return
	}

	putPolicy := rs.PutPolicy{
		Scope: this.bucketName,
		// CallbackUrl:  callbackUrl,
		// CallbackBody: callbackBody,
		// ReturnUrl:    returnUrl,
		// ReturnBody:   returnBody,
		// AsyncOps:     asyncOps,
		// EndUser:      endUser,
		// 指定上传凭证有效期（默认1小时）
		// Expires:      expires,
	}

	this.locker.Lock()
	this.uptoken = putPolicy.Token(nil)
	this.locker.Unlock()
	this.tokenTime = time.Now()
}
func (this *Uploader) uploadMemoryFile(r gio.Reader, key string, size int) (err error) {
	this.genUpToken()

	var ret io.PutRet
	var extra = &io.PutExtra{
		// Params:   params,
		// MimeType: mieType,
		// Crc32:    crc32,
		// CheckCrc: CheckCrc,
	}

	// ret       变量用于存取返回的信息，详情见 io.PutRet
	// uptoken   为业务服务器端生成的上传口令
	// key       为文件存储的标识
	// r         为io.Reader类型，用于从其读取数据
	// extra     为上传文件的额外信息,可为空， 详情见 io.PutExtra, 可选
	err = io.Put2(nil, &ret, this.uptoken, key, r, int64(size), extra)

	// 上传产生错误
	if err != nil {
		logger.Errorln("io.Put failed:", err)

		errInfo := make(map[string]interface{})
		err = json.Unmarshal([]byte(err.Error()), &errInfo)
		if err != nil {
			logger.Errorln("io.Put Unmarshal failed:", err)
			return
		}

		code, ok := errInfo["code"]
		if ok && code == 614 {
			err = nil
		}

		return
	}

	// 上传成功，处理返回值
	logger.Debugln(ret.Hash, ret.Key)

	return
}

var Up = &Uploader{}

func TestUploadMemoryFile(t *testing.T) {
	Up.InitQiniu()
	buf, _ := ioutil.ReadFile("C:\\Users\\eric\\Pictures\\Saved Pictures\\user_default.png")
	reader := bytes.NewReader(buf)

	err := Up.uploadMemoryFile(reader, "test", len(buf))
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		print("ok")
	}
}
