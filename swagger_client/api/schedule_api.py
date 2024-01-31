# coding: utf-8

"""
    Phystech Radio - OpenAPI 3.0

    This is a Phystech Radio specification.  # noqa: E501

    OpenAPI spec version: 0.0.1
    
    Generated by: https://github.com/swagger-api/swagger-codegen.git
"""

from __future__ import absolute_import

import re  # noqa: F401

# python 2 and python 3 compatibility library
import six

from swagger_client.api_client import ApiClient

import json


class ScheduleApi(object):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    Ref: https://github.com/swagger-api/swagger-codegen
    """

    def __init__(self, api_client=None):
        if api_client is None:
            api_client = ApiClient()
        self.api_client = api_client

    def admin_schedule_get(self, **kwargs):  # noqa: E501
        """admin_schedule_get  # noqa: E501

        Get schedule cut  # noqa: E501
        This method makes a synchronous HTTP request by default. To make an
        asynchronous HTTP request, please pass async_req=True
        >>> thread = api.admin_schedule_get(async_req=True)
        >>> result = thread.get()

        :param async_req bool
        :param int start: start cut
        :param int stop: stop cut
        :return: Segments
                 If the method is called asynchronously,
                 returns the request thread.
        """
        kwargs['_return_http_data_only'] = True
        if kwargs.get('async_req'):
            return self.admin_schedule_get_with_http_info(**kwargs)  # noqa: E501
        else:
            (data) = self.admin_schedule_get_with_http_info(**kwargs)  # noqa: E501
            return data

    def admin_schedule_get_with_http_info(self, **kwargs):  # noqa: E501
        """admin_schedule_get  # noqa: E501

        Get schedule cut  # noqa: E501
        This method makes a synchronous HTTP request by default. To make an
        asynchronous HTTP request, please pass async_req=True
        >>> thread = api.admin_schedule_get_with_http_info(async_req=True)
        >>> result = thread.get()

        :param async_req bool
        :param int start: start cut
        :param int stop: stop cut
        :return: Segments
                 If the method is called asynchronously,
                 returns the request thread.
        """

        all_params = ['start', 'stop']  # noqa: E501
        all_params.append('async_req')
        all_params.append('_return_http_data_only')
        all_params.append('_preload_content')
        all_params.append('_request_timeout')

        params = locals()
        for key, val in six.iteritems(params['kwargs']):
            if key not in all_params:
                raise TypeError(
                    "Got an unexpected keyword argument '%s'"
                    " to method admin_schedule_get" % key
                )
            params[key] = val
        del params['kwargs']

        collection_formats = {}

        path_params = {}

        query_params = []
        if 'start' in params:
            query_params.append(('start', params['start']))  # noqa: E501
        if 'stop' in params:
            query_params.append(('stop', params['stop']))  # noqa: E501

        header_params = {}

        form_params = []
        local_var_files = {}

        body_params = None
        # HTTP header `Accept`
        header_params['Accept'] = self.api_client.select_header_accept(
            ['application/json'])  # noqa: E501

        # Authentication setting
        auth_settings = ['editorAuth']  # noqa: E501
        header_params['Authorization'] = self.api_client.configuration.get_api_key_with_prefix('Authorization') # my insert for jwt

        return self.api_client.call_api(
            '/admin/schedule', 'GET',
            path_params,
            query_params,
            header_params,
            body=body_params,
            post_params=form_params,
            files=local_var_files,
            response_type='Segments',  # noqa: E501
            auth_settings=auth_settings,
            async_req=params.get('async_req'),
            _return_http_data_only=params.get('_return_http_data_only'),
            _preload_content=params.get('_preload_content', True),
            _request_timeout=params.get('_request_timeout'),
            collection_formats=collection_formats)

    def admin_schedule_id_delete(self, id, **kwargs):  # noqa: E501
        """admin_schedule_id_delete  # noqa: E501

        Delete segment by id  # noqa: E501
        This method makes a synchronous HTTP request by default. To make an
        asynchronous HTTP request, please pass async_req=True
        >>> thread = api.admin_schedule_id_delete(id, async_req=True)
        >>> result = thread.get()

        :param async_req bool
        :param int id: ID (required)
        :return: None
                 If the method is called asynchronously,
                 returns the request thread.
        """
        kwargs['_return_http_data_only'] = True
        if kwargs.get('async_req'):
            return self.admin_schedule_id_delete_with_http_info(id, **kwargs)  # noqa: E501
        else:
            (data) = self.admin_schedule_id_delete_with_http_info(id, **kwargs)  # noqa: E501
            return data

    def admin_schedule_id_delete_with_http_info(self, id, **kwargs):  # noqa: E501
        """admin_schedule_id_delete  # noqa: E501

        Delete segment by id  # noqa: E501
        This method makes a synchronous HTTP request by default. To make an
        asynchronous HTTP request, please pass async_req=True
        >>> thread = api.admin_schedule_id_delete_with_http_info(id, async_req=True)
        >>> result = thread.get()

        :param async_req bool
        :param int id: ID (required)
        :return: None
                 If the method is called asynchronously,
                 returns the request thread.
        """

        all_params = ['id']  # noqa: E501
        all_params.append('async_req')
        all_params.append('_return_http_data_only')
        all_params.append('_preload_content')
        all_params.append('_request_timeout')

        params = locals()
        for key, val in six.iteritems(params['kwargs']):
            if key not in all_params:
                raise TypeError(
                    "Got an unexpected keyword argument '%s'"
                    " to method admin_schedule_id_delete" % key
                )
            params[key] = val
        del params['kwargs']
        # verify the required parameter 'id' is set
        if ('id' not in params or
                params['id'] is None):
            raise ValueError("Missing the required parameter `id` when calling `admin_schedule_id_delete`")  # noqa: E501

        collection_formats = {}

        path_params = {}
        if 'id' in params:
            path_params['id'] = params['id']  # noqa: E501

        query_params = []

        header_params = {}

        form_params = []
        local_var_files = {}

        body_params = None
        # HTTP header `Accept`
        header_params['Accept'] = self.api_client.select_header_accept(
            ['application/json'])  # noqa: E501

        # Authentication setting
        auth_settings = ['editorAuth']  # noqa: E501
        header_params['Authorization'] = self.api_client.configuration.get_api_key_with_prefix('Authorization') # my insert for jwt

        return self.api_client.call_api(
            '/admin/schedule/{id}', 'DELETE',
            path_params,
            query_params,
            header_params,
            body=body_params,
            post_params=form_params,
            files=local_var_files,
            response_type=None,  # noqa: E501
            auth_settings=auth_settings,
            async_req=params.get('async_req'),
            _return_http_data_only=params.get('_return_http_data_only'),
            _preload_content=params.get('_preload_content', True),
            _request_timeout=params.get('_request_timeout'),
            collection_formats=collection_formats)

    def admin_schedule_id_get(self, id, **kwargs):  # noqa: E501
        """admin_schedule_id_get  # noqa: E501

        Get segment information by id  # noqa: E501
        This method makes a synchronous HTTP request by default. To make an
        asynchronous HTTP request, please pass async_req=True
        >>> thread = api.admin_schedule_id_get(id, async_req=True)
        >>> result = thread.get()

        :param async_req bool
        :param int id: ID (required)
        :return: Segment
                 If the method is called asynchronously,
                 returns the request thread.
        """
        kwargs['_return_http_data_only'] = True
        if kwargs.get('async_req'):
            return self.admin_schedule_id_get_with_http_info(id, **kwargs)  # noqa: E501
        else:
            (data) = self.admin_schedule_id_get_with_http_info(id, **kwargs)  # noqa: E501
            return data

    def admin_schedule_id_get_with_http_info(self, id, **kwargs):  # noqa: E501
        """admin_schedule_id_get  # noqa: E501

        Get segment information by id  # noqa: E501
        This method makes a synchronous HTTP request by default. To make an
        asynchronous HTTP request, please pass async_req=True
        >>> thread = api.admin_schedule_id_get_with_http_info(id, async_req=True)
        >>> result = thread.get()

        :param async_req bool
        :param int id: ID (required)
        :return: Segment
                 If the method is called asynchronously,
                 returns the request thread.
        """

        all_params = ['id']  # noqa: E501
        all_params.append('async_req')
        all_params.append('_return_http_data_only')
        all_params.append('_preload_content')
        all_params.append('_request_timeout')

        params = locals()
        for key, val in six.iteritems(params['kwargs']):
            if key not in all_params:
                raise TypeError(
                    "Got an unexpected keyword argument '%s'"
                    " to method admin_schedule_id_get" % key
                )
            params[key] = val
        del params['kwargs']
        # verify the required parameter 'id' is set
        if ('id' not in params or
                params['id'] is None):
            raise ValueError("Missing the required parameter `id` when calling `admin_schedule_id_get`")  # noqa: E501

        collection_formats = {}

        path_params = {}
        if 'id' in params:
            path_params['id'] = params['id']  # noqa: E501

        query_params = []

        header_params = {}

        form_params = []
        local_var_files = {}

        body_params = None
        # HTTP header `Accept`
        header_params['Accept'] = self.api_client.select_header_accept(
            ['application/json'])  # noqa: E501

        # Authentication setting
        auth_settings = ['editorAuth']  # noqa: E501
        header_params['Authorization'] = self.api_client.configuration.get_api_key_with_prefix('Authorization') # my insert for jwt

        return self.api_client.call_api(
            '/admin/schedule/{id}', 'GET',
            path_params,
            query_params,
            header_params,
            body=body_params,
            post_params=form_params,
            files=local_var_files,
            response_type='Segment',  # noqa: E501
            auth_settings=auth_settings,
            async_req=params.get('async_req'),
            _return_http_data_only=params.get('_return_http_data_only'),
            _preload_content=params.get('_preload_content', True),
            _request_timeout=params.get('_request_timeout'),
            collection_formats=collection_formats)

    def admin_schedule_post(self, **kwargs):  # noqa: E501
        """admin_schedule_post  # noqa: E501

        Create new segment  # noqa: E501
        This method makes a synchronous HTTP request by default. To make an
        asynchronous HTTP request, please pass async_req=True
        >>> thread = api.admin_schedule_post(async_req=True)
        >>> result = thread.get()

        :param async_req bool
        :param SegmentRegister body: ...
        :return: InlineResponse2001
                 If the method is called asynchronously,
                 returns the request thread.
        """
        kwargs['_return_http_data_only'] = True
        if kwargs.get('async_req'):
            return self.admin_schedule_post_with_http_info(**kwargs)  # noqa: E501
        else:
            (data) = self.admin_schedule_post_with_http_info(**kwargs)  # noqa: E501
            return data

    def admin_schedule_post_with_http_info(self, **kwargs):  # noqa: E501
        """admin_schedule_post  # noqa: E501

        Create new segment  # noqa: E501
        This method makes a synchronous HTTP request by default. To make an
        asynchronous HTTP request, please pass async_req=True
        >>> thread = api.admin_schedule_post_with_http_info(async_req=True)
        >>> result = thread.get()

        :param async_req bool
        :param SegmentRegister body: ...
        :return: InlineResponse2001
                 If the method is called asynchronously,
                 returns the request thread.
        """

        all_params = ['body']  # noqa: E501
        all_params.append('async_req')
        all_params.append('_return_http_data_only')
        all_params.append('_preload_content')
        all_params.append('_request_timeout')

        params = locals()
        for key, val in six.iteritems(params['kwargs']):
            if key not in all_params:
                raise TypeError(
                    "Got an unexpected keyword argument '%s'"
                    " to method admin_schedule_post" % key
                )
            params[key] = val
        del params['kwargs']

        collection_formats = {}

        path_params = {}

        query_params = []

        header_params = {}

        form_params = []
        local_var_files = {}

        body_params = {
            'segment' : {
                'mediaID' : params['body'].media_id,
                'start' : params['body'].start
            }
        }

        # HTTP header `Accept`
        header_params['Accept'] = self.api_client.select_header_accept(
            ['application/json'])  # noqa: E501

        # HTTP header `Content-Type`
        header_params['Content-Type'] = self.api_client.select_header_content_type(  # noqa: E501
            ['application/json'])  # noqa: E501

        # Authentication setting
        auth_settings = ['editorAuth']  # noqa: E501
        header_params['Authorization'] = self.api_client.configuration.get_api_key_with_prefix('Authorization') # my insert for jwt

        return self.api_client.call_api(
            '/admin/schedule', 'POST',
            path_params,
            query_params,
            header_params,
            body=body_params,
            post_params=form_params,
            files=local_var_files,
            response_type='InlineResponse2001',  # noqa: E501
            auth_settings=auth_settings,
            async_req=params.get('async_req'),
            _return_http_data_only=params.get('_return_http_data_only'),
            _preload_content=params.get('_preload_content', True),
            _request_timeout=params.get('_request_timeout'),
            collection_formats=collection_formats)