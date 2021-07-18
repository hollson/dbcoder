-- 创建表
DROP TABLE IF EXISTS public.test_table;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.test_table
(
    id                    serial NOT NULL PRIMARY KEY,                      -- 主键(自动创建tableName_pkey自增序列)
    t_smallserial         smallserial,
    t_bigserial           bigserial,
    t_boolean             boolean                     DEFAULT FALSE,
    t_smallint            smallint                    DEFAULT 1,
    t_integer             integer                     DEFAULT 2,
    t_bigint              bigint                      DEFAULT 3,
    t_int                 int                         DEFAULT 5,
    t_int2                int2                        DEFAULT 6,
    t_int4                int4                        DEFAULT 7,
    t_int8                int8                        DEFAULT 8,
    t_uuid                uuid   NOT NULL             DEFAULT uuid_generate_v4(),
    t_varchar_n           varchar(20)                 DEFAULT 'a',
    t_character_varying   character varying           DEFAULT 'b',
    t_character_varying_n character varying(20)       DEFAULT 'b',
    t_char_n              char(20)                    DEFAULT 'c',
    t_character_n         character(20)               DEFAULT 'd',
    t_text                text                        DEFAULT 'e',
    t_real                real                        DEFAULT 32.00,        --单精度
    t_double              double precision            DEFAULT 64.00,        --双精度
    t_float               float                       DEFAULT 1.01,
    t_float8              float8                      DEFAULT 2.02,
    t_timestamp           timestamp                   DEFAULT CURRENT_TIMESTAMP::date,
    t_timestamp_z         timestamp WITH time zone    DEFAULT CURRENT_TIMESTAMP,
    t_timestamp_nz        timestamp WITHOUT time zone DEFAULT CURRENT_TIMESTAMP,
    t_inet                inet                        DEFAULT '192.168.10.1/24',
    t_cidr                cidr                        DEFAULT '192.168.10/24',
    t_macaddr             macaddr                     DEFAULT '08:00:2b:01:02:03',
    t_interval            interval                    DEFAULT '1 min',      --1s
    t_bytea               bytea                       DEFAULT '\134',       -- '\134',E'\\xFC'
    t_timeZ               time WITH time zone         DEFAULT CURRENT_TIMESTAMP,
    t_timeNZ              time WITHOUT time zone      DEFAULT now(),
    t_date                date                        DEFAULT CURRENT_DATE, -- ,CURRENT_TIMESTAMP::date
    t_decimal             decimal(18, 2)              DEFAULT 18.02,
    t_numeric             numeric(18, 2)              DEFAULT 18.02,
    t_money               money                       DEFAULT 99.99,
    t_bit                 bit                         DEFAULT '1',          --bit(1)
    t_bit_n               bit(3)                      DEFAULT '101'         --bit(N)
);

-- 添加注释
COMMENT ON TABLE public.test_table IS '我是一个类型演示表';
COMMENT ON COLUMN public.test_table.id IS '主键';
COMMENT ON COLUMN public.test_table.t_bit IS 'bit(1),即由0和1构成是字符串';
COMMENT ON COLUMN public.test_table.t_bit_n IS 'bit(n)';
COMMENT ON COLUMN public.test_table.t_real IS '单精度';
COMMENT ON COLUMN public.test_table.t_double IS '双精度';
COMMENT ON COLUMN public.test_table.t_numeric IS '同decimal';
COMMENT ON COLUMN public.test_table.t_money IS '同decimal';
COMMENT ON COLUMN public.test_table.t_inet IS 'ipv4/6';
COMMENT ON COLUMN public.test_table.t_cidr IS 'ipv4/6,不接受子网掩码';
COMMENT ON COLUMN public.test_table.t_macaddr IS 'MAC地址';
COMMENT ON COLUMN public.test_table.t_interval IS '时间间隔,year|month|day|hour|min|second';
COMMENT ON COLUMN public.test_table.t_bytea IS '字节数组';

-- 插入数据
INSERT INTO public.test_table(id)
VALUES (1);

-- 创建索引
CREATE INDEX idx_test ON public.test_table USING btree (t_bigint);
COMMENT ON INDEX public.idx_test IS '一个索引测试';

--------------------------------------------------------------------------------------------------------

-- 数组类型
DROP TABLE IF EXISTS public.test_array;
CREATE TABLE public.test_array
(
    id                      serial NOT NULL PRIMARY KEY,
    arr_boolean             boolean[]                     DEFAULT ARRAY [true, false],
    arr_bytea               bytea[],
    arr_smallint            smallint[]                    DEFAULT ARRAY [1, 1],
    arr_integer             integer[]                     DEFAULT ARRAY [11, 22],
    arr_bigint              bigint[]                      DEFAULT ARRAY [111, 222],
    arr_real                real[]                        DEFAULT ARRAY [111, 222],
    arr_int2d               int[][]                       DEFAULT ARRAY [[1,2,3],[4,5,6]],
    arr_varchar_n           varchar(20)[]                 DEFAULT ARRAY ['a', 'a'],
    arr_character_varying_n character varying(20)[]       DEFAULT ARRAY ['b', 'b'],
    arr_char_n              char(20)[]                    DEFAULT ARRAY ['c', 'c'],
    arr_character_n         character(20)[]               DEFAULT ARRAY ['d', 'd'],
    arr_text                text[]                        DEFAULT ARRAY ['e', 'e'],
    arr_timestamp_z         timestamp WITH time zone[]    DEFAULT ARRAY [CURRENT_TIMESTAMP],
    arr_timestamp_nz        timestamp WITHOUT time zone[] DEFAULT ARRAY [CURRENT_TIMESTAMP],
    arr_bit                 bit[]                         DEFAULT ARRAY ['1'::bit,'0'::bit],
    arr_bit_n               bit(3)[]                      DEFAULT ARRAY ['1'::bit(3),'0'::bit(3)],
    arr_numeric             numeric(18, 2)[]              DEFAULT ARRAY [1.1, 1.2],
    arr_decimal             decimal(18, 2)[]              DEFAULT ARRAY [2.1, 2.2],
    arr_money               money[]                       DEFAULT ARRAY [10.01, 10.02],
    arr_time_z              time WITH time zone[]         DEFAULT ARRAY [CURRENT_TIMESTAMP],
    arr_time_nz             time WITHOUT time zone[]      DEFAULT ARRAY [CURRENT_TIMESTAMP],
    arr_date                date[]                        DEFAULT ARRAY [CURRENT_DATE,now()]
);
COMMENT ON TABLE public.test_array IS '演示数组类型';
COMMENT ON COLUMN public.test_array.arr_bytea IS '二维数组，即[][]byte';

INSERT INTO public.test_array(id)
VALUES (1);

--------------------------------------------------------------------------------------------------------

-- 查看表结构
SELECT a.attnum                                        as seq,
       a.attname                                       AS field_name,       -- 字段表名
       a.attlen                                        AS field_size,       -- 字段大小
       t.typcategory                                   AS field_type_group, -- 类型分组
       a.atttypid::regtype                             AS field_type,       -- 类型
       --format_type(a.atttypid, a.atttypmod)            AS field_type_raw,   -- 原始类型
       COALESCE(ct.contype = 'p', FALSE)               AS is_primary_key,   -- 是否主键
       a.attnotnull                                    AS not_null,         -- 是否为NULL
       COALESCE(pg_get_expr(ad.adbin, ad.adrelid), '') AS default_value,    -- 默认值
       COALESCE(b.description, '')                     AS comment           -- 注释
FROM pg_attribute a
         INNER JOIN ONLY pg_class C ON C.oid = a.attrelid
         INNER JOIN ONLY pg_namespace n ON n.oid = C.relnamespace
         LEFT JOIN pg_constraint ct ON ct.conrelid = C.oid AND a.attnum = ANY (ct.conkey) AND ct.contype = 'p'
         LEFT JOIN pg_attrdef ad ON ad.adrelid = C.oid AND ad.adnum = a.attnum
         LEFT JOIN pg_description b ON a.attrelid = b.objoid AND a.attnum = b.objsubid
         left join pg_type t on a.atttypid = t.oid
WHERE a.attisdropped = FALSE
  AND a.attnum > 0
  AND n.nspname = 'public'
  AND C.relname = 'test_table' -- 表名
ORDER BY a.attnum;


