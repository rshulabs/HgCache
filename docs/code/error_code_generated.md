# 错误码

!! 错误码列表，由 `codegen -type=int -doc` 命令生成，不要对此文件做任何更改。

## 功能说明

如果返回结果中存在 `code` 字段，则表示调用 API 接口失败。例如：

```json
{
  "code": 100101,
  "message": "Database error"
}
```

上述返回中 `code` 表示错误码，`message` 表示该错误的具体信息。每个错误同时也对应一个 HTTP 状态码，比如上述错误码对应了 HTTP 状态码 500(Internal Server Error)。

## 错误码列表

支持的错误码列表如下：

| Identifier | Code | HTTP Code | Description |
| ---------- | ---- | --------- | ----------- |
| ErrSuccess | 100001 | 200 | OK |
| ErrUnknown | 100002 | 500 | Internal server error |
| ErrValidation | 100003 | 400 | Validation failed |
| ErrNotFound | 100004 | 404 | Page not found |
| ErrTokenInvalid | 100005 | 401 | Token invalid |
| ErrBind | 100006 | 400 | Error occurred while binding the request body to the struct |
| ErrDatabase | 100101 | 500 | Database error |
| ErrEncrypt | 100201 | 401 | Error occurred while encrypting the user password |
| ErrSignatureInvalid | 100202 | 401 | Signature invalid |
| ErrInvalidAuthHeader | 100203 | 401 | Invalid authorization header |
| ErrMissingHeader | 100204 | 401 | Missing authorization header |
| ErrPasswordIncorrect | 100205 | 401 | Password incorrect |
| ErrPermissionDenied | 100206 | 403 | Permission denied |
| ErrBlackListCheck | 100207 | 401 | Black list check failed |
| ErrGuardTokenCheck | 100208 | 401 | Guard token check failed |
| ErrEncodingFailed | 100301 | 500 | Encoding failed due to an error with the data |
| ErrDecodingFailed | 100302 | 500 | Decoding failed due to an error with the data |
| ErrInvalidJSON | 100303 | 500 | Invalid json data |
| ErrEncodingJSON | 100304 | 500 | JSON data could not be encoded |
| ErrDecodingJSON | 100305 | 500 | JSON data could not be decoded |
| ErrInvalidYaml | 100306 | 500 | Invalid yaml data |
| ErrEncodingYaml | 100307 | 500 | YAML data could not be encoded |
| ErrDecodingYAML | 100308 | 500 | YAML data could not be decoded |

