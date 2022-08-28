-- +goose Up

CREATE TABLE IF NOT EXISTS t_tasks (
    key_hash    UUID primary key,
    task_id     VARCHAR(24),
    hash        UUID,
    create_ts   timestamp default now(),
    owner_login text,

    task_name   text,
    description text,

    status      SMALLINT,
    end_ts      timestamp,
    email_list  text[],
    actual      BOOLEAN
);

comment on table t_tasks is 'Общая информация о задачах';

comment on column t_tasks.task_id is 'Идентификатор задачи';
comment on column t_tasks.hash is 'Хэш от структуры таски';
comment on column t_tasks.create_ts is 'Дата и время создания задачи';
comment on column t_tasks.owner_login is 'Создатель задачи';
comment on column t_tasks.task_name is 'Наименование задачи';
comment on column t_tasks.description is 'Описание задачи';
comment on column t_tasks.status is 'Статус задачи: 1 - согласована; 2 - Не согласована; 3 - В процессе согласования';
comment on column t_tasks.end_ts is 'Дата и время согласования/не согласования';
comment on column t_tasks.email_list is 'Список согласователей';
comment on column t_tasks.actual is 'Флаг актуальности задачи: 1 - Задача актуальна; 2 - задача удалена';

CREATE UNIQUE INDEX IF NOT EXISTS t_tasks_task_id_idx ON t_tasks (task_id);


CREATE TABLE IF NOT EXISTS t_task_status (
     key_hash  UUID primary key ,
     task_id   VARCHAR(24)      ,
     create_ts timestamp        ,
     email     text             ,
     status    SMALLINT         ,
     diff_time NUMERIC(20, 2)
);

CREATE INDEX IF NOT EXISTS t_task_status_task_id_sum_dif_idx ON t_task_status (task_id);

comment on table t_task_status is 'Информация по согласованию задач';

comment on column t_task_status.key_hash is 'Хэш от task_id + email';
comment on column t_task_status.task_id is 'Идентификатор задачи';
comment on column t_task_status.create_ts is 'Дата и время нажатия';
comment on column t_task_status.email is 'Согласователь';
comment on column t_task_status.status is 'Статус задачи: 1 - согласована; 2 - Не согласована';
comment on column t_task_status.diff_time is 'Время в секундах. Для первого согласователя: между датой и временем создания задачи и датой нажатия. Для остальных: между текущим и предыдущим нажатием';


CREATE TABLE IF NOT EXISTS t_mails (
    hash  UUID primary key,
    header     TEXT,
    body       TEXT,
    create_ts  timestamp,
    email_list text[]
);