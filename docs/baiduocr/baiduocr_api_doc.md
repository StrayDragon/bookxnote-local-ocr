---
updated: 2024-11-03
---

# API文档示例
> https://ai.baidu.com/ai-doc/OCR/1k3h7y3db

## 通用文字识别（标准版）

用户向服务请求识别某张图中的所有文字。

```python
    """ 读取文件 """
   def get_file_content(filePath):
      with open(filePath, "rb") as fp:
         return fp.read()

   image = get_file_content('文件路径')
   url = "https://www.x.com/sample.jpg"
   pdf_file = get_file_content('文件路径')

   # 调用通用文字识别（标准版）
   res_image = client.basicGeneral(image)
   res_url = client.basicGeneralUrl(url)
   res_pdf = client.basicGeneralPdf(pdf_file)
   print(res_image)
   print(res_url)
   print(res_pdf)

   # 如果有可选参数
   options = {}
   options["language_type"] = "CHN_ENG"
   options["detect_direction"] = "true"
   options["detect_language"] = "true"
   options["probability"] = "true"
   res_image = client.basicGeneral(image, options)
   res_url = client.basicGeneralUrl(url, options)
   res_pdf = client.basicGeneralPdf(pdf_file, options)
   print(res_image)
   print(res_url)
   print(res_pdf)

```

**通用文字识别 请求参数详情**

| 参数 | 是否必选 | 类型 | 可选值范围 | 说明 |
| --- | --- | --- | --- | --- |
| image | 和 url/pdf\_file 三选一 | string | \- | 图像数据，base64编码后进行urlencode，要求base64编码和urlencode后大小不超过4M，最短边至少15px，最长边最大4096px，支持jpg/jpeg/png/bmp格式   **优先级**：image > url > pdf\_file，当image字段存在时，url、pdf\_file字段失效 |
| url | 和 image/pdf\_file 三选一 | string | \- | 图片完整url，url长度不超过1024字节，url对应的图片base64编码后大小不超过4M，最短边至少15px，最长边最大4096px，支持jpg/jpeg/png/bmp格式   **优先级**：image > url > pdf\_file，当image字段存在时，url字段失效   **请注意关闭URL防盗链** |
| pdf\_file | 和 image/url 三选一 | string | \- | PDF文件，base64编码后进行urlencode，要求base64编码和urlencode后大小不超过4M，最短边至少15px，最长边最大4096px   **优先级**：image > url > pdf\_file，当image、url字段存在时，pdf\_file字段失效 |
| pdf\_file\_num | 否 | string | \- | 需要识别的PDF文件的对应页码，当 pdf\_file 参数有效时，识别传入页码的对应页面内容，若不传入，则默认识别第 1 页 |
| language\_type | 否 | string | CHN\_ENG   ENG   JAP   KOR   FRE   SPA   POR   GER   ITA   RUS | 识别语言类型，默认为CHN\_ENG   可选值包括：   \- CHN\_ENG：中英文混合   \- ENG：英文   \- JAP：日语   \- KOR：韩语   \- FRE：法语   \- SPA：西班牙语   \- POR：葡萄牙语   \- GER：德语   \- ITA：意大利语   \- RUS：俄语 |
| detect\_direction | 否 | string | true/false | 是否检测图像朝向，默认不检测，即：false。朝向是指输入图像是正常方向、逆时针旋转90/180/270度。可选值包括:   \- true：检测朝向；   \- false：不检测朝向。 |
| detect\_language | 否 | string | true/false | 是否检测语言，默认不检测。当前支持（中文、英语、日语、韩语） |
| paragraph | 否 | string | true/false | 是否输出段落信息 |
| probability | 否 | string | true/false | 是否返回识别结果中每一行的置信度 |

**通用文字识别 返回数据参数详情**

| 字段 | 是否必选 | 类型 | 说明 |
| --- | --- | --- | --- |
| direction | 否 | int32 | 图像方向，当 detect\_direction=true 时返回该字段。   \- - 1：未定义，   \- 0：正向，   \- 1：逆时针90度，   \- 2：逆时针180度，   \- 3：逆时针270度 |
| log\_id | 是 | uint64 | 唯一的log id，用于问题定位 |
| words\_result\_num | 是 | uint32 | 识别结果数，表示words\_result的元素个数 |
| words\_result | 是 | array\[\] | 识别结果数组 |
| \+ words | 否 | string | 识别结果字符串 |
| \+ probability | 否 | object | 识别结果中每一行的置信度值，包含average：行置信度平均值，variance：行置信度方差，min：行置信度最小值，当 probability=true 时返回该字段 |
| paragraphs\_result | 否 | array\[\] | 段落检测结果，当 paragraph=true 时返回该字段 |
| \+ words\_result\_idx | 否 | array\[\] | 一个段落包含的行序号，当 paragraph=true 时返回该字段 |
| language | 否 | int32 | 当 detect\_language=true 时返回该字段 |
| pdf\_file\_size | 否 | string | 传入PDF文件的总页数，当 pdf\_file 参数有效时返回该字段 |

**通用文字识别 返回示例**

```json
{
"log_id": 2471272194,
"words_result_num": 2,
"words_result":
    [
        {"words": " TSINGTAO"},
        {"words": "青島睥酒"}
    ]
}
```

## 通用文字识别（高精度版）

用户向服务请求识别某张图中的所有文字，相对于通用文字识别该产品精度更高，但是识别耗时会稍长。

```python
   """ 读取文件 """
   def get_file_content(filePath):
      with open(filePath, "rb") as fp:
         return fp.read()

   image = get_file_content('文件路径')
   url = "https://www.x.com/sample.jpg"
   pdf_file = get_file_content('文件路径')

	# 调用通用文字识别（高精度版）
   res_image = client.basicAccurate(image)
   res_url = client.basicAccurateUrl(url)
   res_pdf = client.basicAccuratePdf(pdf_file)
   print(res_image)
   print(res_url)
   print(res_pdf)

	# 如果有可选参数
   options = {}
   options["detect_direction"] = "true"
   options["probability"] = "true"
   res_image = client.basicAccurate(image, options)
   res_url = client.basicAccurateUrl(url, options)
   res_pdf = client.basicAccuratePdf(pdf_file, options)
   print(res_image)
   print(res_url)
   print(res_pdf)
```

**通用文字识别（高精度版） 请求参数详情**

| 参数 | 是否必选 | 类型 | 可选值范围 | 说明 |
| --- | --- | --- | --- | --- |
| image | 和 url/pdf\_file 三选一 | string | \- | 图像数据，base64编码后进行urlencode，要求base64编码和urlencode后大小不超过10M，最短边至少15px，最长边最大8192px，支持jpg/jpeg/png/bmp格式   **优先级**：image > url > pdf\_file，当image字段存在时，url、pdf\_file字段失效 |
| url | 和 image/pdf\_file 三选一 | string | \- | 图片完整url，url长度不超过1024字节，url对应的图片base64编码后大小不超过10M，最短边至少15px，最长边最大8192px，支持jpg/jpeg/png/bmp格式   **优先级**：image > url > pdf\_file，当image字段存在时，url字段失效   **请注意关闭URL防盗链** |
| pdf\_file | 和 image/url 三选一 | string | \- | PDF文件，base64编码后进行urlencode，要求base64编码和urlencode后大小不超过10M，最短边至少15px，最长边最大8192px   **优先级**：image > url > pdf\_file，当image、url字段存在时，pdf\_file字段失效 |
| pdf\_file\_num | 否 | string | \- | 需要识别的PDF文件的对应页码，当 pdf\_file 参数有效时，识别传入页码的对应页面内容，若不传入，则默认识别第 1 页 |
| language\_type | 否 | string | auto\_detect   CHN\_ENG   ENG   JAP   KOR   FRE   SPA   POR   GER   ITA   RUS   DAN   DUT   MAL   SWE   IND   POL   ROM   TUR   GRE   HUN | 识别语言类型，默认为CHN\_ENG   可选值包括：   \- auto\_detect：自动检测语言，并识别   \- CHN\_ENG：中英文混合   \- ENG：英文   \- JAP：日语   \- KOR：韩语   \- FRE：法语   \- SPA：西班牙语   \- POR：葡萄牙语   \- GER：德语   \- ITA：意大利语   \- RUS：俄语   \- DAN：丹麦语   \- DUT：荷兰语   \- MAL：马来语   \- SWE：瑞典语   \- IND：印尼语   \- POL：波兰语   \- ROM：罗马尼亚语   \- TUR：土耳其语   \- GRE：希腊语   \- HUN：匈牙利语   \- THA：泰语   \- VIE：越南语   \- ARA：阿拉伯语   \- HIN：印地语 |
| detect\_direction | 否 | string | true/false | 是否检测图像朝向，默认不检测，即：false。朝向是指输入图像是正常方向、逆时针旋转90/180/270度。可选值包括:   \- true：检测朝向；   \- false：不检测朝向 |
| paragraph | 否 | string | true/false | 是否输出段落信息 |
| probability | 否 | string | true/false | 是否返回识别结果中每一行的置信度 |

**通用文字识别（高精度版） 返回数据参数详情**

| 字段 | 是否必选 | 类型 | 说明 |
| --- | --- | --- | --- |
| log\_id | 是 | uint64 | 唯一的log id，用于问题定位 |
| direction | 否 | int32 | 图像方向，当 detect\_direction=true 时返回该字段。   \- - 1：未定义，   \- 0：正向，   \- 1：逆时针90度，   \- 2：逆时针180度，   \- 3：逆时针270度 |
| words\_result | 是 | array\[\] | 识别结果数组 |
| words\_result\_num | 是 | uint32 | 识别结果数，表示words\_result的元素个数 |
| \+ words | 否 | string | 识别结果字符串 |
| paragraphs\_result | 否 | array\[\] | 段落检测结果，当 paragraph=true 时返回该字段 |
| \+ words\_result\_idx | 否 | array\[\] | 一个段落包含的行序号，当 paragraph=true 时返回该字段 |
| \+ probability | 否 | object | 识别结果中每一行的置信度值，包含average：行置信度平均值，variance：行置信度方差，min：行置信度最小值，当 probability=true 时返回该字段 |
| pdf\_file\_size | 否 | string | 传入PDF文件的总页数，当 pdf\_file 参数有效时返回该字段 |

**通用文字识别（高精度版） 返回示例**

**参考通用文字识别（标准版）返回示例**


# 错误返回格式
> https://ai.baidu.com/ai-doc/OCR/zkibizyhz

若请求错误，服务器将返回的JSON文本包含以下参数：

- **error\_code**：错误码。
- **error\_msg**：错误描述信息，帮助理解和解决发生的错误。

## 错误码

| 错误码 | 错误信息 | 描述 |
| --- | --- | --- |
| 4 | Open api request limit reached | 集群超限额 |
| 6 | No permission to access data | 无权限访问该用户数据，创建应用时未勾选相关接口，请登录百度云控制台，找到对应的应用，编辑应用，勾选上相关接口，然后重试调用 |
| 14 | IAM Certification failed | IAM鉴权失败，建议用户参照文档自查生成sign的方式是否正确，或换用控制台中ak sk的方式调用 |
| 17 | Open api daily request limit reached | 每天流量超限额 |
| 18 | Open api qps request limit reached | QPS超限额 |
| 19 | Open api total request limit reached | 请求总量超限额 |
| 100 | Invalid parameter | 无效参数 |
| 110 | Access token invalid or no longer valid | Access Token失效 |
| 111 | Access token expired | Access token过期 |
| 282000 | internal error | 服务器内部错误，如果您使用的是高精度接口，报这个错误码的原因可能是您上传的图片中文字过多，识别超时导致的，建议您对图片进行切割后再识别，其他情况请再次请求， 如果持续出现此类错误，请通过QQ群（631977213）或工单联系技术支持团队。 |
| 216100 | invalid param | 请求中包含非法参数，请检查后重新尝试 |
| 216101 | not enough param | 缺少必须的参数，请检查参数是否有遗漏 |
| 216102 | service not support | 请求了不支持的服务，请检查调用的url |
| 216103 | param too long | 请求中某些参数过长，请检查后重新尝试 |
| 216110 | appid not exist | appid不存在，请重新核对信息是否为后台应用列表中的appid |
| 216200 | empty image | 图片为空，请检查后重新尝试 |
| 216201 | image format error | 上传的图片格式错误，现阶段我们支持的图片格式为：PNG、JPG、JPEG、BMP，请进行转码或更换图片 |
| 216202 | image size error | 上传的图片大小错误，现阶段我们支持的图片大小为：base64编码后小于4M，分辨率不高于4096\*4096 px，请重新上传图片 |
| 216630 | recognize error | 识别错误，请再次请求，如果持续出现此类错误，请通过QQ群（631977213）或工单联系技术支持团队。 |
| 216631 | recognize bank card error | 识别银行卡错误，出现此问题的原因一般为：您上传的图片非银行卡正面，上传了异形卡的图片或上传的银行卡正品图片不完整 |
| 216633 | recognize idcard error | 识别身份证错误，出现此问题的原因一般为：您上传了非身份证图片或您上传的身份证图片不完整 |
| 216634 | detect error | 检测错误，请再次请求，如果持续出现此类错误，请通过QQ群（631977213）或工单联系技术支持团队。 |
| 282003 | missing parameters: {参数名} | 请求参数缺失 |
| 282005 | batch  processing error | 处理批量任务时发生部分或全部错误，请根据具体错误码排查 |
| 282006 | batch task  limit reached | 批量任务处理数量超出限制，请将任务数量减少到10或10以下 |
| 282110 | urls not exit | URL参数不存在，请核对URL后再次提交 |
| 282111 | url format illegal | URL格式非法，请检查url格式是否符合相应接口的入参要求 |
| 282112 | url download timeout | url下载超时，请检查url对应的图床/图片无法下载或链路状况不好，您可以重新尝试以下，如果多次尝试后仍不行，建议更换图片地址 |
| 282113 | url response invalid | URL返回无效参数 |
| 282114 | url size error | URL长度超过1024字节或为0 |
| 282808 | request id: xxxxx not exist | request id xxxxx 不存在 |
| 282809 | result type error | 返回结果请求错误（不属于excel或json） |
| 282810 | image recognize error | 图像识别错误 |