import os
import tempfile
import random
from datetime import (date, datetime, timedelta)

import telebot
from dotenv import load_dotenv

from client import client

TMP_DIR = 'tmp'
LOG_DIR = '.log'
SCHEDULE_STEP = 10

cl = client()

load_dotenv()
app = telebot.TeleBot(token=os.environ.get('TOKEN'))

login_candidates = {}
new_media_candidates = {}
status = {}
search_res = {}
schedule_cut = {}

LIBRARY = '📚🎶'
NEW_MEDIA = '🆕🎵'
SCHEDULE = '🗓🎼'
HELP_MESSAGE = 'хелп...'
MAIN_MENU = 'Главное меню'

@app.message_handler(commands=['start'])
def start(msg: telebot.types.Message) -> None:
    user_id = msg.from_user.id
    if cl.is_user_known(user_id):
        status[user_id] = 'logged'
        app.send_message(
            chat_id=user_id,
            text=f'Кастуй или пиздуй, {cl.get_user_name(user_id)}.',
            reply_markup=main_menu()
        )
    else:
        sended_msg = app.send_message(
            chat_id=user_id,
            text='Приветствую, радист, для начала надо зарегистрироваться.\n' + \
                'Введи свой логин.',
        )
        login_candidates[user_id] = {
            'login' : None,
            'pass' : None,
            'messages' : [msg.message_id, sended_msg.message_id]
        }
        status[user_id] = 'candidate'
        
@app.message_handler(regexp=MAIN_MENU)
def main_menu_handler(msg: telebot.types.Message) -> None:
    user_id = msg.from_user.id
    app.send_message(
        chat_id=user_id,
        text='Кастуй или пиздуй',
        reply_markup=main_menu()
    )

@app.message_handler(regexp=LIBRARY)
def library(msg: telebot.types.Message) -> None:
    user_id = msg.from_user.id

    if invalid_user(user_id):
        return

    try:
        cl.update_library(user_id)
    except Exception as e:
        log(user_id, 'failed to update library. %s' % e)
        fail_message(user_id)
        return
    
    try:
        cl.update_schedule(user_id)
    except Exception as e:
        log(user_id, 'failed to update schedule. %s' % e)
        fail_message(user_id)
        return
    
    status[user_id] = 'library-search'
    app.send_message(
        chat_id=user_id,
        text='Напиши часть названия или автора.',
        reply_markup=to_main_menu()
    )

@app.message_handler(regexp=SCHEDULE)
def schedule(msg: telebot.types.Message) -> None:
    user_id = msg.from_user.id

    if invalid_user(user_id):
        return

    try:
        cl.update_library(user_id)
    except Exception as e:
        log(user_id, 'failed to update library. %s' % e)
        fail_message(user_id)

    try:
        cl.update_schedule(user_id)
    except Exception as e:
        log(user_id, 'failed to update schedule. %s' % e)
        fail_message(user_id)

    schedule_cut[user_id] = get_actual_schedule()

    markup = telebot.types.InlineKeyboardMarkup(row_width=1).add(
        telebot.types.InlineKeyboardButton("Расписание", callback_data="0_schedule"),
        telebot.types.InlineKeyboardButton("Автодиджей", callback_data="autodj")
    )

    app.send_message(
        chat_id=user_id,
        text='Выбери день или авто диджея',
        reply_markup=markup
    )

@app.message_handler(regexp=NEW_MEDIA)
def new_media(msg: telebot.types.Message) -> None:
    user_id = msg.from_user.id
    if invalid_user(user_id):
        return

    app.send_message(
        chat_id=user_id,
        text='Хочешь добавить что-то новенькое? Для начала пришли мне файл в формате .mp3 (не более 20 МБ)',
    )
    status[msg.from_user.id] = 'new-media-upload-file'

@app.message_handler(regexp=HELP_MESSAGE)
def help_message(msg: telebot.types.Message) -> None:
    app.send_message(
        chat_id=msg.from_user.id,
        text='Сейчас я объясню свой функционал!',
        reply_markup=main_menu()
    )

@app.callback_query_handler(lambda call: call.data.endswith('_schedule'))
def schedule_list(call: telebot.types.CallbackQuery) -> None:
    user_id = call.from_user.id

    if invalid_user(user_id):
        return

    data = call.data.removesuffix('_schedule')

    if data == 'main':
        schedule_cut[user_id] = get_actual_schedule()
        markup = telebot.types.InlineKeyboardMarkup(row_width=1).add(
            telebot.types.InlineKeyboardButton("Расписание", callback_data="0_schedule"),
            telebot.types.InlineKeyboardButton("Автодиджей", callback_data="autodj")
        )
        app.send_message(
            chat_id=user_id,
            text='Выбери день или авто диджея',
            reply_markup=markup
        )
        return

    i = int(data)

    result = [segment_pretty(x) for x in schedule_cut[user_id][i*SCHEDULE_STEP:(i+1)*SCHEDULE_STEP]]

    text = "\n".join(result) if len(result) != 0 else "Расписание пусто."

    app.edit_message_text(
        chat_id=user_id,
        message_id=call.message.message_id,
        text=text,
        reply_markup=arrows_menu(i, len(result) == 0)
    )

@app.callback_query_handler(lambda call: call.data == 'autodj')
def autodj_handler(call: telebot.types.CallbackQuery) -> None:
    user_id = call.from_user.id

    if len(cl.library) == 0:
        app.send_message(
            chat_id=user_id,
            text='Упс, а у меня ничегошеньки нет, чтобы в расписание вставить...'
        )
        return

    status[user_id] = 'autodj-hours'

    app.send_message(
        chat_id=user_id,
        text='Сколько часов ты хочешь забить в расписании?'
    )

@app.callback_query_handler(lambda call: call.data == 'new-segment')
def new_segment_ask_id(call: telebot.types.CallbackQuery) -> None:
    user_id = call.from_user.id

    status[user_id] = 'new-segment-id'

    app.send_message(
        chat_id=user_id,
        text='Введи номер композиции, которую ты хочешь добавить.'
    )

@app.callback_query_handler(lambda call: call.data == 'new-media-send')
def new_media_send(call: telebot.types.CallbackQuery) -> None:
    user_id = call.from_user.id

    if invalid_user(user_id):
        return    

    status[user_id] = None

    sended_msg = app.edit_message_text(
        chat_id=user_id,
        message_id=call.message.id,
        text='Отправляю...'
    )

    try:
        cl.upload_media(
            user_id, 
            new_media_candidates[user_id]['name'],
            new_media_candidates[user_id]['author'],
            new_media_candidates[user_id]['source']
        )
    except Exception as e:
        log(user_id, 'failed to upload media. %s' %e)
        fail_message(user_id)
        return

    try:
        cl.update_library(user_id)
    except Exception as e:
        log(user_id, 'failed to update library. %s' % e)
        fail_message(user_id)
        return

    try:
        os.remove(new_media_candidates[user_id]['source'])
    except Exception as e:
        log(user_id, 'failef to delete tmp file. %s' % e)
    finally:
        del new_media_candidates[user_id]

    app.edit_message_text(
        chat_id=user_id,
        message_id=sended_msg.message_id,
        text='Загружено!',
        # reply_markup=main_menu()
    )

@app.callback_query_handler(lambda call: call.data == 'new-media-name')
def new_media_name(call: telebot.types.CallbackQuery) -> None:
    user_id = call.from_user.id

    if invalid_user(user_id):
        return
    
    status[user_id] = 'new-media-name'
    new_media_candidates[user_id]['name'] = None
    new_media_candidates[user_id]['author'] = None

    app.edit_message_text(
        chat_id=user_id,
        message_id=call.message.message_id,
        text="Отлично, теперь скажи мне название."
    )

@app.message_handler(content_types=['text'])
def text_message_handler(msg: telebot.types.Message) -> None:    
    user_id = msg.from_user.id

    if status[user_id] == 'candidate':
        authorization(msg)
        return

    if invalid_user(user_id):
        return
    
    match status[user_id]:
        case 'library-search':
            library_search(msg)
        case 'new-media-name':
            new_media_name(msg)
        case 'new-media-author':
            new_media_author(msg)
        case 'new-segment-id':
            new_segment_id(msg)
        case 'autodj-hours':
            autodj_hours(msg)

@app.message_handler(content_types=['audio'])
def upload_media(msg: telebot.types.Message) -> None:
    user_id = msg.from_user.id
    audio = msg.audio

    if invalid_user(user_id):
        return    

    if status[user_id] != 'new-media-upload-file':
        app.send_message(
            chat_id=user_id,
            text='Не ожидал от тебя такого...'
        )
        return

    if audio.mime_type != "audio/mpeg":
        app.send_message(
            chat_id=user_id,
            text='Я могу принять только .mp3 файл.'
        )
        return

    sended_msg = app.send_message(
        chat_id=user_id,
        text="Качаю..."
    )

    file_info = app.get_file(audio.file_id)
    downloaded_file = app.download_file(file_info.file_path)

    _, filename = tempfile.mkstemp(dir=TMP_DIR, prefix='audio-')

    with open(filename, 'wb') as new_file:
        new_file.write(downloaded_file)

    new_media_candidates[user_id] = {
        'name' : None,
        'author' : None,
        'source' : filename,
    }

    if msg.audio.file_name.count(' - ') == 1:
        author, name = msg.audio.file_name.split(' - ')
        name = name.removesuffix('.mp3')
        new_media_candidates[user_id]['name'] = name
        new_media_candidates[user_id]['author'] = author
        status[user_id] = None
        app.edit_message_text(
            chat_id=user_id,
            message_id=sended_msg.message_id,
            text=f"Название песни \"{name}\", а автор - \"{author}\"?"
        )
        app.edit_message_reply_markup(
            chat_id=user_id,
            message_id=sended_msg.message_id,
            reply_markup=telebot.types.InlineKeyboardMarkup(row_width=2).add(
                telebot.types.InlineKeyboardButton(text='Да', callback_data='new-media-send'),
                telebot.types.InlineKeyboardButton(text='Нет', callback_data='new-media-name')
            )
        )
        return

    app.edit_message_text(
        chat_id=user_id,
        message_id=sended_msg.message_id,
        text="Отлично, теперь скажи мне название."
    )
    status[user_id] = 'new-media-name'

def main_menu() -> telebot.types.ReplyKeyboardMarkup:
    return telebot.types.ReplyKeyboardMarkup(resize_keyboard=True, row_width=2).add(
        telebot.types.KeyboardButton(LIBRARY),
        telebot.types.KeyboardButton(SCHEDULE),
        telebot.types.KeyboardButton(NEW_MEDIA),
        telebot.types.KeyboardButton(HELP_MESSAGE),
    )

def to_main_menu() -> telebot.types.ReplyKeyboardMarkup:
    return telebot.types.ReplyKeyboardMarkup(resize_keyboard=True).add(
        telebot.types.KeyboardButton(MAIN_MENU),
    )

def arrows_menu(i: int, empty: bool) -> telebot.types.InlineKeyboardMarkup:
    if i > 0:
        if empty:
            return telebot.types.InlineKeyboardMarkup(row_width=1).add(
                telebot.types.InlineKeyboardButton(text='🔙', callback_data=f'{i-1}_schedule'),
                telebot.types.InlineKeyboardButton(text='Назад', callback_data=f'main_schedule')
            )
        return telebot.types.InlineKeyboardMarkup(row_width=2).add(
            telebot.types.InlineKeyboardButton(text='🔙', callback_data=f'{i-1}_schedule'),
            telebot.types.InlineKeyboardButton(text='🔜', callback_data=f'{i+1}_schedule'),
            telebot.types.InlineKeyboardButton(text='Назад', callback_data=f'main_schedule')
        )
    else:
        if empty:
            return telebot.types.InlineKeyboardMarkup(row_width=1).add(
                telebot.types.InlineKeyboardButton(text='Назад', callback_data=f'main_schedule')
            )
        return telebot.types.InlineKeyboardMarkup(row_width=1).add(
            telebot.types.InlineKeyboardButton(text='🔜', callback_data=f'{i+1}_schedule'),
            telebot.types.InlineKeyboardButton(text='Назад', callback_data=f'main_schedule')
        )

def authorization(msg: telebot.types.Message) -> None:
    user_id = msg.from_user.id
    if cl.is_user_known(user_id):
        return

    if login_candidates[user_id]['login'] is None:
        login = msg.text
        if login is None or login == '':
            app.send_message(
                chat_id=user_id,
                text='Логин не может быть пустым'
            )
            return
        
        login_candidates[user_id]['login'] = login
        sended_msg = app.send_message(
            chat_id=user_id,
            text='Отлично, теперь введи свой пароль'
        )
        login_candidates[user_id]['messages'].append(msg.message_id)
        login_candidates[user_id]['messages'].append(sended_msg.message_id)
        return
    elif login_candidates[user_id]['pass'] is None:
        password = msg.text
        if password is None or password == '':
            app.send_message(
                chat_id=user_id,
                text='Пароль не может быть пустым'
            )
            return
        
        login_candidates[user_id]['pass'] = password
        login_candidates[user_id]['messages'].append(msg.message_id)
        app.delete_messages(user_id, login_candidates[user_id]['messages'])
        
        try:
            res = cl.login(
                user_id,
                login_candidates[user_id]['login'],
                login_candidates[user_id]['pass']    
            )
        except Exception as e:
            log(user_id, 'Error occured during client->login. %s' % e)
            fail_message(user_id)

        if res:
            del login_candidates[user_id]
            status[user_id] = None
            app.send_message(
                chat_id=user_id,
                text=f'Привет, {cl.get_user_name(user_id)}! Начнем?',
                reply_markup=main_menu()
            )
        else:
            login_candidates[user_id] = login_candidates[user_id] = {
                'login' : None,
                'pass' : None,
                'messages' : [msg.message_id]
            }
            app.send_message(
                chat_id=user_id,
                text='Неверный логин или пароль, попробуй еще раз.'
            )
        return
    else:
        fail_message(user_id)

def new_media_name(msg: telebot.types.Message) -> None:
    user_id = msg.from_user.id
    name = msg.text

    if name is None or name == '':
        app.send_message(
            chat_id=user_id,
            text='Название не может быть пустым',
        )
    if new_media_candidates[user_id]['name'] is not None:
        fail_message(user_id)
        return

    new_media_candidates[user_id]['name'] = name
    status[user_id] = 'new-media-author'

    app.send_message(
        chat_id=user_id,
        text='Отлично, теперь назови автора.'
    )

def new_media_author(msg: telebot.types.Message) -> None:
    user_id = msg.from_user.id
    author = msg.text

    if author is None or author == '':
        app.send_message(
            chat_id=user_id,
            text='Имя автора не может быть пустым',
        )
    if new_media_candidates[user_id]['author'] is not None:
        fail_message(user_id)

    new_media_candidates[user_id]['author'] = author
    status[user_id] = None

    sended_msg = app.send_message(
        chat_id=user_id,
        text='Отправляю...'
    )

    try:
        cl.upload_media(
            user_id, 
            new_media_candidates[user_id]['name'],
            new_media_candidates[user_id]['author'],
            new_media_candidates[user_id]['source']
        )
    except Exception as e:
        log(user_id, 'failed to upload media. %s' %e)
        fail_message(user_id)
        return

    try:
        cl.update_library(user_id)
    except Exception as e:
        log(user_id, 'failed to update library. %s' % e)
        fail_message(user_id)
        return

    try:
        os.remove(new_media_candidates[user_id]['source'])
    except Exception as e:
        log(user_id, 'failef to delete tmp file. %s' % e)
    finally:
        del new_media_candidates[user_id]

    app.edit_message_text(
        chat_id=user_id,
        message_id=sended_msg.message_id,
        text='Загружено!',
        # reply_markup=main_menu()
    )

def new_segment_id(msg: telebot.types.Message) -> None:
    user_id = msg.from_user.id

    try:
        array_id = int(msg.text)
    except Exception as e:
        app.send_message(
            chat_id=user_id,
            text='Ты отправил мне что-то не то...'
        )
        return
    
    if not ( 1 <= array_id <= len(search_res[user_id]) ):
        app.send_message(
            chat_id=user_id,
            text='Некорректное число.'
        )
        return

    try:
        res = cl.new_segment(user_id, search_res[user_id][array_id - 1]['id'])
        segm_id = res.id
    except Exception as e:
        log(user_id, 'failed to create new segment. %s' %e)
        fail_message(user_id)
        return
    
    status[user_id] = 'library-search'

    try:
        cl.update_schedule(user_id)
    except Exception as e:
        log(user_id, 'failed to update schedule. %s' %e)
        fail_message(user_id)
        return

    app.send_message(
        chat_id=user_id,
        text=f'Добавлено {segment_pretty(cl.get_segment(segm_id))}\n' +\
            'Можешь добавить что-то еще.'
    )

def library_search(msg: telebot.types.Message) -> None:
    user_id = msg.from_user.id
    pattern = msg.text
    if pattern is None or pattern == '':
        app.send_message(
            chat_id=user_id,
            text='Ключевое слово не может быть пустым'
        )
        return
    filtered = filter(
        lambda x: search_alg(pattern, x['name']) or search_alg(pattern, x['author']),
        cl.library    
    )
    search_res[user_id] = list(filtered)
    result = [f'{i+1}) {media_pretty(x)}' for i, x in enumerate(search_res[user_id])]

    text = "\n".join(result) if len(result) != 0 else 'К сожалению, ничего не найдено, попробуйте поискать что-то другое.'

    markup = telebot.types.InlineKeyboardMarkup().add(
        telebot.types.InlineKeyboardButton(text='Добавить в расписание', callback_data='new-segment')
    ) if len(result) != 0 else None

    app.send_message(
        chat_id=user_id,
        text=text,
        reply_markup=markup
    )

def autodj_hours(msg: telebot.types.Message) -> None:
    user_id = msg.from_user.id

    try:
        hours = float(msg.text)
    except Exception as e:
        app.send_message(
            chat_id=user_id,
            text='Ты отправил мне что-то не то...'
        )
        return
    
    if hours < 0:
        app.send_message(
            chat_id=user_id,
            text='Ты куда в минуса полез...'
        )
        return
    if hours == 0:
        app.send_message(
            chat_id=user_id,
            text='Ну, ноль так ноль. Делать ничего не буду.'
        )
        return
    if hours > 24*3:
        app.send_message(
            chat_id=user_id,
            text='Больше трех дней за раз, силен конечно. Давай в несколько заходов.'
        )
        return

    sended_msg = app.send_message(
        chat_id=user_id,
        text='Работаем, радисты...'
    )

    media_to_send = []
    duration_sum = timedelta(0)

    indices = list(range(len(cl.library)))
    random.shuffle(indices)
    while duration_sum < timedelta(hours=hours):
        media_id = indices.pop()
        if len(indices) == 0:
            indices = list(range(len(cl.library)))
            random.shuffle(indices)
        media_to_send.append(media_id)

        duration = cl.lirary[media_id]['duration']
        duration_sum += timedelta(microseconds=duration*1e-3)

    for m in media_to_send:
        try:
            cl.new_segment(user_id, m)
        except Exception as e:
            log(user_id, f'failed to create new segment id={m} %s' %e)
            fail_message(user_id)
            return

    try:
        cl.update_schedule(user_id)
    except Exception as e:
        log(user_id, f'failed to update schedule %s' %e)
        fail_message(user_id)
        return

    app.delete_message(user_id, sended_msg.message_id)
    app.send_message(
        chat_id=user_id,
        text=f'Расписание забито до {cl.time_horizon.strftime(r"%Y-%m-%d %H:%M:%S")}'
    )

    status[user_id] = 'library-search'

def search_alg(pattern: str, text: str) -> bool:
    return pattern.lower() in text.lower()

def get_actual_schedule() -> list[dict]:
    start = datetime.today()
    sch = []

    for segm in cl.schedule:
        segm_start = datetime.strptime(segm['start'], r'%Y-%m-%dT%H:%M:%S.%fZ') + timedelta(hours=3)
        segm_stop = segm_start + timedelta(microseconds=segm['stopCut']*1e-3)
        if start < segm_stop:
            sch.append(segm)
    return sch

def media_pretty(media: dict) -> str:
    return f'{media["name"]} \u2014 {media["author"]}'

def segment_pretty(segm: dict) -> str:
    segm_start = datetime.strptime(segm['start'], r'%Y-%m-%dT%H:%M:%S.%fZ') + timedelta(hours=3)
    segm_stop = segm_start + timedelta(microseconds=segm['stopCut']*1e-3)

    media = cl.get_media(segm['mediaID'])

    time_diap = segm_start.strftime('%H:%M:%S') + '-' + segm_stop.strftime('%H:%M:%S')
    return time_diap + '  ' + media_pretty(media)

def invalid_user(user_id: int) -> bool:
    if not cl.is_user_known(user_id):
        app.send_message(
            chat_id=user_id,
            text='Я тебя не знаю, разбойник!'
        )
        return True
    if user_id not in status.keys():
        status[user_id] = None
    return False

def fail_message(user_id: int) -> None:
    app.send_message(
        chat_id=user_id,
        text='Произошла ошибка, тыкайте админа.'
    )

def log(user_id: int, *args) -> None:
    now = datetime.now()
    args = now.strftime(f'[%Y-%m-%d %H:%M:%S, {user_id}] ') + ... + args
    with open(LOG_DIR+'/bot.log', 'a') as wr:
        print(*args, file=wr)

if __name__ == '__main__':
    app.polling(non_stop=True)