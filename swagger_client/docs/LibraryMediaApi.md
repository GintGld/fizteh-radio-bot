# swagger_client.LibraryMediaApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**admin_library_media_get**](LibraryMediaApi.md#admin_library_media_get) | **GET** /admin/library/media | 
[**admin_library_media_id_delete**](LibraryMediaApi.md#admin_library_media_id_delete) | **DELETE** /admin/library/media/{id} | 
[**admin_library_media_id_get**](LibraryMediaApi.md#admin_library_media_id_get) | **GET** /admin/library/media/{id} | 
[**admin_library_media_post**](LibraryMediaApi.md#admin_library_media_post) | **POST** /admin/library/media | 
[**admin_library_source_id_get**](LibraryMediaApi.md#admin_library_source_id_get) | **GET** /admin/library/source/{id} | 

# **admin_library_media_get**
> MediaArray admin_library_media_get()



Get all info media

### Example
```python
from __future__ import print_function
import time
import swagger_client
from swagger_client.rest import ApiException
from pprint import pprint


# create an instance of the API class
api_instance = swagger_client.LibraryMediaApi(swagger_client.ApiClient(configuration))

try:
    api_response = api_instance.admin_library_media_get()
    pprint(api_response)
except ApiException as e:
    print("Exception when calling LibraryMediaApi->admin_library_media_get: %s\n" % e)
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**MediaArray**](MediaArray.md)

### Authorization

[editorAuth](../README.md#editorAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **admin_library_media_id_delete**
> admin_library_media_id_delete(id)



Delete media and its source by id

### Example
```python
from __future__ import print_function
import time
import swagger_client
from swagger_client.rest import ApiException
from pprint import pprint


# create an instance of the API class
api_instance = swagger_client.LibraryMediaApi(swagger_client.ApiClient(configuration))
id = 789 # int | ID

try:
    api_instance.admin_library_media_id_delete(id)
except ApiException as e:
    print("Exception when calling LibraryMediaApi->admin_library_media_id_delete: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| ID | 

### Return type

void (empty response body)

### Authorization

[editorAuth](../README.md#editorAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **admin_library_media_id_get**
> Media admin_library_media_id_get(id)



Get media information by id

### Example
```python
from __future__ import print_function
import time
import swagger_client
from swagger_client.rest import ApiException
from pprint import pprint


# create an instance of the API class
api_instance = swagger_client.LibraryMediaApi(swagger_client.ApiClient(configuration))
id = 789 # int | ID

try:
    api_response = api_instance.admin_library_media_id_get(id)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling LibraryMediaApi->admin_library_media_id_get: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| ID | 

### Return type

[**Media**](Media.md)

### Authorization

[editorAuth](../README.md#editorAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **admin_library_media_post**
> admin_library_media_post(media=media, source=source)



Post media with its source

### Example
```python
from __future__ import print_function
import time
import swagger_client
from swagger_client.rest import ApiException
from pprint import pprint


# create an instance of the API class
api_instance = swagger_client.LibraryMediaApi(swagger_client.ApiClient(configuration))
media = swagger_client.MediaRegister() # MediaRegister |  (optional)
source = 'source_example' # str |  (optional)

try:
    api_instance.admin_library_media_post(media=media, source=source)
except ApiException as e:
    print("Exception when calling LibraryMediaApi->admin_library_media_post: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **media** | [**MediaRegister**](.md)|  | [optional] 
 **source** | **str**|  | [optional] 

### Return type

void (empty response body)

### Authorization

[editorAuth](../README.md#editorAuth)

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **admin_library_source_id_get**
> str admin_library_source_id_get(id)



Get media source by id

### Example
```python
from __future__ import print_function
import time
import swagger_client
from swagger_client.rest import ApiException
from pprint import pprint


# create an instance of the API class
api_instance = swagger_client.LibraryMediaApi(swagger_client.ApiClient(configuration))
id = 789 # int | ID

try:
    api_response = api_instance.admin_library_source_id_get(id)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling LibraryMediaApi->admin_library_source_id_get: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| ID | 

### Return type

**str**

### Authorization

[editorAuth](../README.md#editorAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: audio/mpeg, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

