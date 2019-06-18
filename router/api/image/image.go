package image

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"git.yitum.com/saas/shop-admin/model"
	"git.yitum.com/saas/shop-admin/model/mysql"
	"git.yitum.com/saas/shop-admin/model/trans"
	"git.yitum.com/saas/shop-admin/pkg/bootstrap"
	"git.yitum.com/saas/shop-admin/router/api"
	"git.yitum.com/saas/shop-admin/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Add(c *gin.Context) {
	// base64.StdEncoding.DecodeString(datasource)

	reqModel := trans.ReqImageAdd{}
	if err := c.Bind(&reqModel); err != nil {
		model.Logger.Info(err.Error(), zap.Int("code", 1))
		api.JSON(c, api.MsgErr, "bad request")
		return
	}
	ext := ".jpg"
	dataType := reqModel.Image[strings.IndexByte(reqModel.Image, ':')+1 : strings.IndexByte(reqModel.Image, ';')]
	b64data := reqModel.Image[strings.IndexByte(reqModel.Image, ',')+1:]

	fmt.Println(dataType)
	b64dataDecode, _ := base64.StdEncoding.DecodeString(b64data) // 成图片文件并把文件写入到buffer

	fileName := generateUniqueMd5()
	rootPath, month := generatePath()

	filePath := filepath.Join(rootPath, fileName+ext)

	path := filepath.Dir(filePath)

	os.MkdirAll(path, os.ModePerm)

	err := ioutil.WriteFile(filePath, b64dataDecode, 0666) // buffer输出到jpg文件中（不做处理，直接写到文件）
	if err != nil {
		fmt.Println(err.Error())
		api.JSON(c, api.MsgErr, "bad create file")
		return
	}
	// fmt.Println(filePath)
	showPath := bootstrap.Conf.Image.Domain + bootstrap.Conf.Image.Space + "/" + month + "/" + fileName + ext

	var faImage = mysql.Image{
		Name: fileName,
		Type: dataType,
		Url:  showPath + "/200_200",
	}
	if err := service.Image.Create(c, model.Db, &faImage); err != nil {
		api.JSON(c, api.MsgErr, "image create file error")
		return
	}

	api.JSONOK(c, gin.H{
		"bmind": trans.RespImage{
			Name: fileName,
			Path: showPath + "/120_120",
		},
		"origin": trans.RespImage{
			Name: fileName,
			Path: showPath + "/200_200",
		},
		"small": trans.RespImage{
			Name: fileName,
			Path: showPath + "/60_60",
		},
	})
}

func List(c *gin.Context) {
	req := trans.ReqImageList{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request list params is error")
		return
	}
	total, list := service.Image.ListPage(c, mysql.Conds{}, req.ReqPage)
	api.JSONList(c, list, total)

}

func GoodsList(c *gin.Context) {
	req := trans.ReqGoodsImageList{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request list params is error")
		return
	}
	// total, list := service.GoodsImage.ListPage(c, req.ReqPage, auth.Default(c).Id)
	total, list := service.GoodsImage.ListPage(c, mysql.Conds{}, req.ReqPage)
	api.JSONList(c, list, total)

}

func generatePath() (string, string) {
	month := time.Now().Format("200601")
	return fmt.Sprintf("%s/%s/", bootstrap.Conf.Image.Path, month), month
}

func generateUniqueMd5() string {
	date := time.Now().Format("20060102150405")
	uniqueID := GenerateUniqueID()
	sno := date + bootstrap.Conf.System.Hostname + string(bootstrap.Conf.System.Pid) + uniqueID

	return fmt.Sprintf("%x", md5.Sum([]byte(sno)))
}

func GenerateUniqueID() string {
	b := make([]byte, 16)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		model.Logger.Error("GenerateUniqueId error", zap.String("err", err.Error()))
	}

	return hex.EncodeToString(b)
}
