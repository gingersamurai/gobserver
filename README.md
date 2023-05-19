# gobserver
Консольная утилита, которое позволяет следить за изменениями в различных директориях
и выполнять произвольный набор консольных команд

## Установка и конфигурация
+ Склонировать репозиторий
    ```bash
    git clone https://github.com/gingersamurai/gobserver.git
    cd gobserver
    ```
+ Запустить СУБД _PostgreSQL_
    ```
    make local_postgres_init 
    ```
+ Установить переменные окружения для миграций и связи приложения с СУБД 
    ```bash
    set -a && source .env && set +a
    ```
+ Запустить миграции
    ```bash
    make migrate
    ```
+ Скомпилировать приложение
    ```bash
    make build
    ```
  В папке `gobserver/build/` появится исполняемый файл `gobserver-cli`
+ настроить файл конфигурации `gobserver/config.yaml`
+ Запустить `gobserver-cli`
    ```bash
    ./build/gobserver-cli
    ```

с примером работы можно ознакомиться по [ссылке](https://www.youtube.com/watch?v=wahTz_VXRMM).

## Архитектура
С архитекутрой приложения можно ознакомиться по [ссылке](https://viewer.diagrams.net/?tags=%7B%7D&highlight=0000ff&edit=_blank&layers=1&nav=1&title=gobserever#Uhttps%3A%2F%2Fdrive.google.com%2Fuc%3Fid%3D1UAzhKiBOSijiNY3n6azzFZglPbTcUVQ4%26export%3Ddownload)


## Список проблем, с которыми пришлось столкнуться:
+ Как бесконечно мониторить состояние файла?\
    **Решение:** Использовать командную утилиту `inotifywait`.
+ Как без головной боли работать с `inotifywait` в Go?\
    **Решение:**  Воспользоваться [готовой оболочкой](https://github.com/fsnotify/fsnotify).

+ Как запускать команды _shell_ формата из конфига в _exec_ формате? \
    **Решение:** Небольшой костыль в формате `['bash', '-c', command]`
    