# swagger_client.AuthApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**admin_login_post**](AuthApi.md#admin_login_post) | **POST** /admin/login | login editor

# **admin_login_post**
> InlineResponse200 admin_login_post(body=body)

login editor

### Example
```python
from __future__ import print_function
import time
import swagger_client
from swagger_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = swagger_client.AuthApi()
body = swagger_client.LoginForm() # LoginForm | Login form (optional)

try:
    # login editor
    api_response = api_instance.admin_login_post(body=body)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling AuthApi->admin_login_post: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**LoginForm**](LoginForm.md)| Login form | [optional] 

### Return type

[**InlineResponse200**](InlineResponse200.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, aplication/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

