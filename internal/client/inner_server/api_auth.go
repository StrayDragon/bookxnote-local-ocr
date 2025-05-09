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


type AuthAPI interface {

	/*
	Oauth20TokenPost 获取 OAuth token

	提供一个mock的OAuth token用于百度OCR API兼容

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return AuthAPIOauth20TokenPostRequest
	*/
	Oauth20TokenPost(ctx context.Context) AuthAPIOauth20TokenPostRequest

	// Oauth20TokenPostExecute executes the request
	//  @return HandlersAPITokenResp
	Oauth20TokenPostExecute(r AuthAPIOauth20TokenPostRequest) (*HandlersAPITokenResp, *http.Response, error)
}

// AuthAPIService AuthAPI service
type AuthAPIService service

type AuthAPIOauth20TokenPostRequest struct {
	ctx context.Context
	ApiService AuthAPI
	clientId *string
	clientSecret *string
	grantType *string
}

// Client ID
func (r AuthAPIOauth20TokenPostRequest) ClientId(clientId string) AuthAPIOauth20TokenPostRequest {
	r.clientId = &clientId
	return r
}

// Client secret
func (r AuthAPIOauth20TokenPostRequest) ClientSecret(clientSecret string) AuthAPIOauth20TokenPostRequest {
	r.clientSecret = &clientSecret
	return r
}

// Grant type
func (r AuthAPIOauth20TokenPostRequest) GrantType(grantType string) AuthAPIOauth20TokenPostRequest {
	r.grantType = &grantType
	return r
}

func (r AuthAPIOauth20TokenPostRequest) Execute() (*HandlersAPITokenResp, *http.Response, error) {
	return r.ApiService.Oauth20TokenPostExecute(r)
}

/*
Oauth20TokenPost 获取 OAuth token

提供一个mock的OAuth token用于百度OCR API兼容

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return AuthAPIOauth20TokenPostRequest
*/
func (a *AuthAPIService) Oauth20TokenPost(ctx context.Context) AuthAPIOauth20TokenPostRequest {
	return AuthAPIOauth20TokenPostRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return HandlersAPITokenResp
func (a *AuthAPIService) Oauth20TokenPostExecute(r AuthAPIOauth20TokenPostRequest) (*HandlersAPITokenResp, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *HandlersAPITokenResp
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AuthAPIService.Oauth20TokenPost")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauth/2.0/token"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.clientId == nil {
		return localVarReturnValue, nil, reportError("clientId is required and must be specified")
	}
	if r.clientSecret == nil {
		return localVarReturnValue, nil, reportError("clientSecret is required and must be specified")
	}
	if r.grantType == nil {
		return localVarReturnValue, nil, reportError("grantType is required and must be specified")
	}

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
	parameterAddToHeaderOrQuery(localVarFormParams, "client_id", r.clientId, "", "")
	parameterAddToHeaderOrQuery(localVarFormParams, "client_secret", r.clientSecret, "", "")
	parameterAddToHeaderOrQuery(localVarFormParams, "grant_type", r.grantType, "", "")
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
