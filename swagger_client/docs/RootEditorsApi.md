# swagger_client.RootEditorsApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**admin_root_editor_id_delete**](RootEditorsApi.md#admin_root_editor_id_delete) | **DELETE** /admin/root/editor/{id} | 
[**admin_root_editor_id_get**](RootEditorsApi.md#admin_root_editor_id_get) | **GET** /admin/root/editor/{id} | 
[**admin_root_editors_get**](RootEditorsApi.md#admin_root_editors_get) | **GET** /admin/root/editors | 
[**admin_root_editors_post**](RootEditorsApi.md#admin_root_editors_post) | **POST** /admin/root/editors | 

# **admin_root_editor_id_delete**
> admin_root_editor_id_delete(id)



Delete editor by id

### Example
```python
from __future__ import print_function
import time
import swagger_client
from swagger_client.rest import ApiException
from pprint import pprint


# create an instance of the API class
api_instance = swagger_client.RootEditorsApi(swagger_client.ApiClient(configuration))
id = 789 # int | ID

try:
    api_instance.admin_root_editor_id_delete(id)
except ApiException as e:
    print("Exception when calling RootEditorsApi->admin_root_editor_id_delete: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| ID | 

### Return type

void (empty response body)

### Authorization

[rootAuth](../README.md#rootAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **admin_root_editor_id_get**
> Editor admin_root_editor_id_get(id)



Get editor by id

### Example
```python
from __future__ import print_function
import time
import swagger_client
from swagger_client.rest import ApiException
from pprint import pprint


# create an instance of the API class
api_instance = swagger_client.RootEditorsApi(swagger_client.ApiClient(configuration))
id = 789 # int | ID

try:
    api_response = api_instance.admin_root_editor_id_get(id)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling RootEditorsApi->admin_root_editor_id_get: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| ID | 

### Return type

[**Editor**](Editor.md)

### Authorization

[rootAuth](../README.md#rootAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **admin_root_editors_get**
> Editors admin_root_editors_get()



Get all users

### Example
```python
from __future__ import print_function
import time
import swagger_client
from swagger_client.rest import ApiException
from pprint import pprint


# create an instance of the API class
api_instance = swagger_client.RootEditorsApi(swagger_client.ApiClient(configuration))

try:
    api_response = api_instance.admin_root_editors_get()
    pprint(api_response)
except ApiException as e:
    print("Exception when calling RootEditorsApi->admin_root_editors_get: %s\n" % e)
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**Editors**](Editors.md)

### Authorization

[rootAuth](../README.md#rootAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **admin_root_editors_post**
> InlineResponse2001 admin_root_editors_post(body=body)



Create new editor

### Example
```python
from __future__ import print_function
import time
import swagger_client
from swagger_client.rest import ApiException
from pprint import pprint


# create an instance of the API class
api_instance = swagger_client.RootEditorsApi(swagger_client.ApiClient(configuration))
body = swagger_client.LoginForm() # LoginForm | register form (optional)

try:
    api_response = api_instance.admin_root_editors_post(body=body)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling RootEditorsApi->admin_root_editors_post: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**LoginForm**](LoginForm.md)| register form | [optional] 

### Return type

[**InlineResponse2001**](InlineResponse2001.md)

### Authorization

[rootAuth](../README.md#rootAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

