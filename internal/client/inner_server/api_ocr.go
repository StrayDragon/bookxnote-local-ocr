/*
BookxNote Local OCR API

This is a local OCR service for BookxNote that mimics Baidu OCR API

API version: 1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package inner_server

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
)


type OcrAPI interface {

	/*
	Rest20OcrV1AccurateBasicPost 对图片进行OCR识别

	使用OCR服务识别图片中的文字

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return OcrAPIRest20OcrV1AccurateBasicPostRequest
	*/
	Rest20OcrV1AccurateBasicPost(ctx context.Context) OcrAPIRest20OcrV1AccurateBasicPostRequest

	// Rest20OcrV1AccurateBasicPostExecute executes the request
	//  @return HandlersAPIAccurateOCRResp
	Rest20OcrV1AccurateBasicPostExecute(r OcrAPIRest20OcrV1AccurateBasicPostRequest) (*HandlersAPIAccurateOCRResp, *http.Response, error)
}

// OcrAPIService OcrAPI service
type OcrAPIService service

type OcrAPIRest20OcrV1AccurateBasicPostRequest struct {
	ctx context.Context
	ApiService OcrAPI
	detectDirection *string
	image *string
	languageType *string
	paragraph *string
	pdfFile *string
	pdfFileNum *string
	probability *string
	url *string
}

// (无效)是否检测文字方向
func (r OcrAPIRest20OcrV1AccurateBasicPostRequest) DetectDirection(detectDirection string) OcrAPIRest20OcrV1AccurateBasicPostRequest {
	r.detectDirection = &detectDirection
	return r
}

// Base64编码的图片数据
func (r OcrAPIRest20OcrV1AccurateBasicPostRequest) Image(image string) OcrAPIRest20OcrV1AccurateBasicPostRequest {
	r.image = &image
	return r
}

// (无效)语言类型
func (r OcrAPIRest20OcrV1AccurateBasicPostRequest) LanguageType(languageType string) OcrAPIRest20OcrV1AccurateBasicPostRequest {
	r.languageType = &languageType
	return r
}

// (无效)是否将文字组织成段落
func (r OcrAPIRest20OcrV1AccurateBasicPostRequest) Paragraph(paragraph string) OcrAPIRest20OcrV1AccurateBasicPostRequest {
	r.paragraph = &paragraph
	return r
}

// (无效)Base64编码的PDF文件
func (r OcrAPIRest20OcrV1AccurateBasicPostRequest) PdfFile(pdfFile string) OcrAPIRest20OcrV1AccurateBasicPostRequest {
	r.pdfFile = &pdfFile
	return r
}

// (无效)PDF页码
func (r OcrAPIRest20OcrV1AccurateBasicPostRequest) PdfFileNum(pdfFileNum string) OcrAPIRest20OcrV1AccurateBasicPostRequest {
	r.pdfFileNum = &pdfFileNum
	return r
}

// (无效)是否返回字符概率
func (r OcrAPIRest20OcrV1AccurateBasicPostRequest) Probability(probability string) OcrAPIRest20OcrV1AccurateBasicPostRequest {
	r.probability = &probability
	return r
}

// (无效)图片的URL
func (r OcrAPIRest20OcrV1AccurateBasicPostRequest) Url(url string) OcrAPIRest20OcrV1AccurateBasicPostRequest {
	r.url = &url
	return r
}

func (r OcrAPIRest20OcrV1AccurateBasicPostRequest) Execute() (*HandlersAPIAccurateOCRResp, *http.Response, error) {
	return r.ApiService.Rest20OcrV1AccurateBasicPostExecute(r)
}

/*
Rest20OcrV1AccurateBasicPost 对图片进行OCR识别

使用OCR服务识别图片中的文字

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return OcrAPIRest20OcrV1AccurateBasicPostRequest
*/
func (a *OcrAPIService) Rest20OcrV1AccurateBasicPost(ctx context.Context) OcrAPIRest20OcrV1AccurateBasicPostRequest {
	return OcrAPIRest20OcrV1AccurateBasicPostRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return HandlersAPIAccurateOCRResp
func (a *OcrAPIService) Rest20OcrV1AccurateBasicPostExecute(r OcrAPIRest20OcrV1AccurateBasicPostRequest) (*HandlersAPIAccurateOCRResp, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *HandlersAPIAccurateOCRResp
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OcrAPIService.Rest20OcrV1AccurateBasicPost")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/rest/2.0/ocr/v1/accurate_basic"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/x-www-form-urlencoded"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.detectDirection != nil {
		parameterAddToHeaderOrQuery(localVarFormParams, "detect_direction", r.detectDirection, "", "")
	}
	if r.image != nil {
		parameterAddToHeaderOrQuery(localVarFormParams, "image", r.image, "", "")
	}
	if r.languageType != nil {
		parameterAddToHeaderOrQuery(localVarFormParams, "language_type", r.languageType, "", "")
	}
	if r.paragraph != nil {
		parameterAddToHeaderOrQuery(localVarFormParams, "paragraph", r.paragraph, "", "")
	}
	if r.pdfFile != nil {
		parameterAddToHeaderOrQuery(localVarFormParams, "pdf_file", r.pdfFile, "", "")
	}
	if r.pdfFileNum != nil {
		parameterAddToHeaderOrQuery(localVarFormParams, "pdf_file_num", r.pdfFileNum, "", "")
	}
	if r.probability != nil {
		parameterAddToHeaderOrQuery(localVarFormParams, "probability", r.probability, "", "")
	}
	if r.url != nil {
		parameterAddToHeaderOrQuery(localVarFormParams, "url", r.url, "", "")
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 400 {
			var v HandlersErrorResp
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v HandlersErrorResp
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
