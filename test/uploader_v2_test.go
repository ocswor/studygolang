package test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/polaris1119/logger"
	"github.com/qiniu/api.v6/rs"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	gio "io"
	"io/ioutil"
	"sync"
	"testing"
	"time"
)

type UploaderV2 struct {
	bucketName string

	uptoken   string
	tokenTime time.Time
	locker    sync.RWMutex
}

func (this *UploaderV2) InitQiniu() {
	logger.Init("D:\\goproject\\studygolang\\log", "DEBUG")
	//ACCESS_KEY := "lCcyC4KgEa58sTjC4sllbw4lfSOuTP62lFVegwZQ"
	//SECRET_KEY := "06F0Pxu8EeSRbMyXYXD6JrjYOFEAXZ5N9GtaK3Wy"
	//UP_HOST := "http://cdn.studytrade.club/"
	//this.bucketName = "Blog"

}
func (this *UploaderV2) genUpToken() {
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
func (this *UploaderV2) uploadLocalFile(key string) {
	putPolicy := storage.PutPolicy{
		Scope: "Blog",
	}
	mac := qbox.NewMac("lCcyC4KgEa58sTjC4sllbw4lfSOuTP62lFVegwZQ", "06F0Pxu8EeSRbMyXYXD6JrjYOFEAXZ5N9GtaK3Wy")
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	resumeUploader := storage.NewResumeUploaderV2(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.RputV2Extra{}
	err := resumeUploader.PutFile(context.Background(), &ret, upToken, key, "C:\\Users\\eric\\Pictures\\Saved Pictures\\user_default.png", &putExtra)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret.Key, ret.Hash)

	return
}

func (this *UploaderV2) uploadMemoryFile(r gio.Reader, key string, size int64) {
	putPolicy := storage.PutPolicy{
		Scope: "Blog",
	}
	mac := qbox.NewMac("lCcyC4KgEa58sTjC4sllbw4lfSOuTP62lFVegwZQ", "06F0Pxu8EeSRbMyXYXD6JrjYOFEAXZ5N9GtaK3Wy")
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	//putExtra := storage.PutExtra{
	//	Params: map[string]string{
	//		"x:name": "github logo",
	//	},
	//}
	putExtra := storage.PutExtra{}
	err := formUploader.Put(context.Background(), &ret, upToken, key, r, size, &putExtra)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret.Key, ret.Hash)

	return
}

var Up2 = &UploaderV2{}

func TestUploadMemory2File(t *testing.T) {
	Up2.InitQiniu()
	buf, _ := ioutil.ReadFile("C:\\Users\\eric\\Pictures\\Saved Pictures\\user_default.png")
	reader := bytes.NewReader(buf)

	//Up2.uploadLocalFile("foo/user_default.png")
	Up2.uploadMemoryFile(reader, "foo/user_default2.png", int64(len(buf)))
}
