from time import time
from datetime import (datetime, timedelta)
import jwt
import json

from swagger_client import (
    ApiClient,
    LoginForm,
    MediaRegister,
    SegmentRegister,
    AuthApi,
    LibraryMediaApi,
    ScheduleApi,
    Configuration
)
from swagger_client.rest import ApiException

class client:
    def __init__(self) -> None:
        self.cache_dir = '.cache'
        self.jwts = {}
        self.user_info = self.__recover_user_info()
        self.library = None
        self.schedule = None
        self.time_horizon = None

    def is_user_known(self, user_id: int) -> bool:
        return True if str(user_id) in self.user_info.keys() else False

    def get_user_name(self, user_id: int) -> str:
        return self.user_info[str(user_id)]['login']

    def __config(self, user_id: int) -> ApiClient:
        '''
            Custom config for jwt required api
        '''
        configuration = Configuration()
        configuration.api_key['Authorization'] = self.jwts[user_id].token
        configuration.api_key_prefix['Authorization'] = 'Bearer'
        configuration.temp_folder_path = 'tmp'
        return configuration

    def __add_token(self, user_id: int, token: str) -> None:
        payload =  jwt.decode(token, options={"verify_signature": False})
        timeout = payload['exp']
        self.jwts[user_id] = JWT(token, timeout)

    def __refresh_jwt_if_needed(self, user_id: int) -> None:
        now = time()
        if user_id not in self.jwts or now > self.jwts[user_id].timeout:
            self.login(user_id, *self.user_info[str(user_id)].values())

    def __recover_user_info(self) -> dict:
        try:
            r = open(self.cache_dir+'/users.json')
            return json.load(r)
        except FileNotFoundError:
            return {}
        except Exception as e:
            print("Exception when calling __recover_user_info: %s\n" % e)
            raise e

    def __save_user(self, user_id: int, login: str, password: str) -> None:
        self.user_info[str(user_id)] = {
            'login' : login,
            'pass' : password
        }
        with open(self.cache_dir+'/users.json', 'w') as wr:
            json.dump(self.user_info, wr)

    def get_media(self, media_id: int) -> dict:
        return next((x for x in self.library if x['id'] == media_id), None)

    def get_segment(self, segm_id: int) -> dict:
        return next((x for x in self.schedule if x['id'] == segm_id), None)

    def login(self, user_id: int, login: str, password: str) -> bool:
        api_instance = AuthApi(ApiClient())

        body = LoginForm(login=login, _pass=password)

        try:
            api_response = api_instance.admin_login_post(body=body)
            self.__add_token(user_id, api_response.token)
            self.__save_user(user_id, login, password)
            return True
        except ApiException as e:
            print("Exception when calling AuthApi->admin_login_post: %s\n" % e)
            if e.status == 400:
                return False
            return None
        
    def update_library(self, user_id: int) -> None:
        self.__refresh_jwt_if_needed(user_id)

        api_instance = LibraryMediaApi(
            ApiClient(self.__config(user_id)))

        try:
            api_response = api_instance.admin_library_media_get()
            self.library = api_response['library']
        except ApiException as e:
            print("Exception when calling LibraryMediaApi->admin_library_media_get: %s\n" % e)
            raise e

    def upload_media(self, user_id: int, name: str, author: str, source: str) -> None:
        self.__refresh_jwt_if_needed(user_id)

        api_instance = LibraryMediaApi(
            ApiClient(self.__config(user_id)))
        media = MediaRegister(name=name, author=author)

        try:
            api_instance.admin_library_media_post(media=media, source=source)
        except ApiException as e:
            print("Exception when calling LibraryMediaApi->admin_library_media_post: %s\n" % e)
            raise e

    def load_media(self, user_id: int, media_id: int) -> str:
        self.__refresh_jwt_if_needed(user_id)
        
        api_instance = LibraryMediaApi(
            ApiClient(self.__config(user_id)))

        try:
            api_response = api_instance.admin_library_source_id_get(media_id)
            # TODO: make tmp file and return its name
            with open('test.mp3', 'wb') as wr:
                wr.write(api_response.urllib3_response.data)

        except ApiException as e:
            print("Exception when calling LibraryMediaApi->admin_library_source_id_get: %s\n" % e)
            raise e

    def update_schedule(self, user_id: int) -> None:
        self.__refresh_jwt_if_needed(user_id)
        
        api_instance = ScheduleApi(
            ApiClient(self.__config(user_id)))

        # unix timestamp current day, 00:00
        today = datetime.now().date().strftime('%s')
        start = today

        try:
            api_response = api_instance.admin_schedule_get(start=start)
            self.schedule = api_response['segments']
            if len(self.schedule) > 0:
                start = datetime.strptime(self.schedule[-1]['start'], r'%Y-%m-%dT%H:%M:%S.%fZ')
                duration = timedelta(microseconds=self.schedule[-1]['stopCut']*1e-3)
                self.time_horizon = start + duration + timedelta(hours=3) # time zone shift
        except ApiException as e:
            print("Exception when calling ScheduleApi->admin_schedule_get: %s\n" % e)
            raise e

    def new_segment(self, user_id: int, media_id: int) -> None:
        self.__refresh_jwt_if_needed(user_id)
        
        api_instance = ScheduleApi(
            ApiClient(self.__config(user_id)))
        
        if self.time_horizon is None:
            self.time_horizon = datetime.now()

        media = self.get_media(media_id)
        if media is None:
            raise ValueError("Unknown media_id.")

        duration = timedelta(microseconds=media['duration']*1e-3)
        start = self.time_horizon

        body = SegmentRegister(
            media_id=media_id,
            start=start.strftime(r"%Y-%m-%dT%H:%M:%S.%f+03:00")
        )

        try:
            res = api_instance.admin_schedule_post(body=body)
            return res
        except ApiException as e:
            print("Exception when calling ScheduleApi->admin_schedule_post: %s\n" % e)
            raise e

class JWT:
    def __init__(self, token, timeout) -> None:
        self.token = token
        self.timeout = timeout
