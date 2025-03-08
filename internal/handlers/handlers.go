package handlers

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/straydragon/bookxnote-local-ocr/internal/service"
)

// 从 gin.Context 获取服务实例的key
const serviceKey = "service"

// GetService 从上下文获取服务实例
func GetService(c *gin.Context) *service.Service {
	return c.MustGet(serviceKey).(*service.Service)
}

type APITokenReq struct {
	GrantType    string `form:"grant_type" binding:"required"`
	ClientID     string `form:"client_id" binding:"required"`
	ClientSecret string `form:"client_secret" binding:"required"`
}

type APITokenResp struct {
	AccessToken   string `json:"access_token"`
	ExpiresIn     int64  `json:"expires_in"`
	RefreshToken  string `json:"refresh_token"`
	Scope         string `json:"scope"`
	SessionKey    string `json:"session_key"`
	SessionSecret string `json:"session_secret"`
}

type APIAccurateOCRReq struct {
	Image           string `form:"image" binding:"required_without_all=URL PDFFile"`
	URL             string `form:"url" binding:"required_without_all=Image PDFFile"`
	PDFFile         string `form:"pdf_file" binding:"required_without_all=Image URL"`
	PDFFileNum      string `form:"pdf_file_num"`
	LanguageType    string `form:"language_type"`
	DetectDirection string `form:"detect_direction"`
	Paragraph       string `form:"paragraph"`
	Probability     string `form:"probability"`
}

type APIAccurateOCRRespWordResult struct {
	Words string `json:"words"`
}

type APIAccurateOCRResp struct {
	LogID          uint64                          `json:"log_id"`
	Direction      int32                           `json:"direction,omitempty"`
	WordsResult    []*APIAccurateOCRRespWordResult `json:"words_result"`
	WordsResultNum uint32                          `json:"words_result_num"`
	XBackend       string                          `json:"x_backend"`
}

// TokenHandler 处理token请求(Hook Baidu OCR API)
// @Summary 获取 OAuth token
// @Description 提供一个mock的OAuth token用于百度OCR API兼容
// @Tags auth
// @Accept x-www-form-urlencoded
// @Produce json
// @Param grant_type formData string true "Grant type"
// @Param client_id formData string true "Client ID"
// @Param client_secret formData string true "Client secret"
// @Success 200 {object} APITokenResp
// @Router /oauth/2.0/token [post]
func TokenHandler(c *gin.Context) {
	c.JSON(200, APITokenResp{
		AccessToken:   "24.460da4889caad24cccf1fbbeb6608.2592000.1458530384.282335-1234567",
		ExpiresIn:     2592000,
		RefreshToken:  "25.92dc5c24c6b507cc54d70e33890d92.315360000.1771798384.282335-1234567",
		Scope:         "public brain_all_scope brain_ocr_general brain_ocr_general_basic",
		SessionKey:    "9mzdXdrN3dKEIW05X7fvX",
		SessionSecret: "197c4081538730d1b3692b7e9faa9f1f",
	})
}

// AccurateOCRHandler 处理OCR请求(Hook Baidu OCR API)
// @Summary 对图片进行OCR识别
// @Description 使用UmiOCR识别图片中的文字
// @Tags ocr
// @Accept x-www-form-urlencoded
// @Produce json
// @Param image formData string false "Base64编码的图片数据"
// @Param url formData string false "(无效)图片的URL"
// @Param pdf_file formData string false "(无效)Base64编码的PDF文件"
// @Param pdf_file_num formData string false "(无效)PDF页码"
// @Param language_type formData string false "(无效)语言类型"
// @Param detect_direction formData string false "(无效)是否检测文字方向"
// @Param paragraph formData string false "(无效)是否将文字组织成段落"
// @Param probability formData string false "(无效)是否返回字符概率"
// @Success 200 {object} APIAccurateOCRResp
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /rest/2.0/ocr/v1/accurate_basic [post]
func AccurateOCRHandler(c *gin.Context) {
	svc := GetService(c)

	// NOTE: 测试 BookxNote OCR 设置的超时时间, 可能连接成功后header中就有timeout, 但是懒得抓包23333
	// for i := 1; i <= 60; i++ {
	// 	log.Printf("OCR识别中, 已进行 %d 秒", i)
	// 	time.Sleep(1 * time.Second)
	// }

	var req APIAccurateOCRReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, ErrInvalidParamResp)
		return
	}

	resp, err := svc.UmiOCRClient.Recognize(req.Image)
	if err != nil {
		log.Printf("OCR recognition failed: %v", err)
		c.JSON(500, ErrInternalServerResp)
		return
	}

	wordsResult := make([]*APIAccurateOCRRespWordResult, 0, len(resp.Data))
	for _, item := range resp.Data {
		wordsResult = append(wordsResult, &APIAccurateOCRRespWordResult{
			Words: item.Text + "\n",
		})
	}

	c.JSON(200, &APIAccurateOCRResp{
		LogID:          uint64(time.Now().UnixNano()),
		WordsResult:    wordsResult,
		WordsResultNum: uint32(len(wordsResult)),
	})
}

// CatchAllHandler 处理未定义的路由
// @Summary 处理未定义的路由
// @Description 返回一个空响应用于任何未定义的路由
// @Tags misc
// @Produce json
// @Success 200 {object} object
// @Router /catch-all [get]
func CatchAllHandler(c *gin.Context) {
	c.JSON(200, gin.H{})
}
