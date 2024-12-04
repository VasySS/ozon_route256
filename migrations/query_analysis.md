# Orders

## Запрос с id

Выполняемый запрос:

```sql
EXPLAIN (analyse, verbose, buffers) SELECT id, user_id, weight, price, expiry_date
FROM user_order
WHERE id = $1;
```

### Без индексов

```
Seq Scan on public.user_order  (cost=0.00..209.00 rows=1 width=40) (actual time=0.393..0.519 rows=1 loops=1)"
  Output: id, user_id, weight, price, expiry_date
  Filter: (user_order.id = 792126562)
  Rows Removed by Filter: 9999
  Buffers: shared hit=84

Planning Time: 0.051 ms
Execution Time: 0.529 ms
```

### Primary key (b-tree)

```
Index Scan using user_order_pkey on public.user_order  (cost=0.29..8.30 rows=1 width=40) (actual time=0.017..0.018 rows=1 loops=1)
  Output: id, user_id, weight, price, expiry_date
  Index Cond: (user_order.id = 981568712)
  Buffers: shared hit=3

Planning Time: 0.065 ms
Execution Time: 0.037 ms
```

### Primary key (b-tree) + hash

`CREATE INDEX user_order_id_custom_idx ON user_order USING hash (id);`

```
Index Scan using user_order_id_custom_idx on public.user_order  (cost=0.00..8.02 rows=1 width=40) (actual time=0.014..0.014 rows=1 loops=1)
  Output: id, user_id, weight, price, expiry_date
  Index Cond: (user_order.id = 854733011)
  Buffers: shared hit=2

Planning Time: 0.050 ms
Execution Time: 0.024 ms
```

### Вывод

Цена выполнения запроса без индексов - `0.00..209.00`, c b-tree - `0.29..8.30` и hash - `0.00..8.02`. Время выполнения запроса уменьшилось в 14 раз - с `0.529 мс` до `0.037 мс`.

Hash индекс работает чуть быстрее и занимает меньше места, но он не может быть использован для сортировки, а также он не может быть UNIQUE. B-tree является оптимальным выбором.

## Запрос с id и user_id

Выполняемый запрос:

```sql
SELECT id, user_id, weight, price, expiry_date
FROM user_order
WHERE user_id = $1
ORDER BY id DESC
LIMIT $2
```

### Primary key (b-tree) для id

```
Limit  (cost=209.01..209.01 rows=1 width=40) (actual time=0.535..0.535 rows=0 loops=1)
  Output: id, user_id, weight, price, expiry_date
  Buffers: shared hit=84
  ->  Sort  (cost=209.01..209.01 rows=1 width=40) (actual time=0.534..0.534 rows=0 loops=1)
        Output: id, user_id, weight, price, expiry_date
        Sort Key: user_order.id DESC
        Sort Method: quicksort  Memory: 25kB
        Buffers: shared hit=84
        ->  Seq Scan on public.user_order  (cost=0.00..209.00 rows=1 width=40) (actual time=0.526..0.527 rows=0 loops=1)
              Output: id, user_id, weight, price, expiry_date
              Filter: (user_order.user_id = 854733011)
              Rows Removed by Filter: 10000
              Buffers: shared hit=84

"Planning Time: 0.069 ms"
"Execution Time: 0.557 ms"
```

### Primary key (b-tree) для id + b-tree для user_id

`CREATE INDEX user_order_user_id_idx ON user_order USING btree (user_id);`

```
Limit  (cost=8.31..8.32 rows=1 width=40) (actual time=0.022..0.023 rows=0 loops=1)
  Output: id, user_id, weight, price, expiry_date
  Buffers: shared hit=2
  ->  Sort  (cost=8.31..8.32 rows=1 width=40) (actual time=0.022..0.022 rows=0 loops=1)
        Output: id, user_id, weight, price, expiry_date
        Sort Key: user_order.id DESC
        Sort Method: quicksort  Memory: 25kB
        Buffers: shared hit=2
        ->  Index Scan using user_order_user_id_idx on public.user_order  (cost=0.29..8.30 rows=1 width=40) (actual time=0.015..0.015 rows=0 loops=1)
              Output: id, user_id, weight, price, expiry_date
              Index Cond: (user_order.user_id = 704498569)
              Buffers: shared hit=2

Planning Time: 0.085 ms
Execution Time: 0.048 ms
```

### Вывод

Цена выполнения запроса без индекса для user_id - `209.01..209.01`, c b-tree - `8.31..8.32`.

Как и в первом запросе, добавление индекса сильно ускорило выполнение запроса - примерно в 11 раз, с `0.557 мс` до `0.048 мс`.

# Order returns

Выполняемый запрос:

```sql
SELECT user_id, order_id
FROM order_return
WHERE order_id = $1
```

### Без индексов

```
Seq Scan on public.order_return  (cost=0.00..180.00 rows=1 width=16) (actual time=0.402..0.500 rows=1 loops=1)
  Output: user_id, order_id
  Filter: (order_return.order_id = 645917160)
  Rows Removed by Filter: 9999
  Buffers: shared hit=55

Planning Time: 0.050 ms
Execution Time: 0.510 ms
```

### B-tree индекс для order_id

```
Index Scan using order_return_order_id_idx on public.order_return  (cost=0.29..8.30 rows=1 width=16) (actual time=0.019..0.020 rows=1 loops=1)
  Output: user_id, order_id
  Index Cond: (order_return.order_id = 729538074)
  Buffers: shared hit=3

Planning Time: 0.062 ms
Execution Time: 0.040 ms
```

### Вывод

Цена выполнения запроса без индекса для user_id - `0.00..180.00`, c b-tree - `0.29..8.30`.

Время выполнения уменьшилось в 12 раз - с `0.510 мс` до `0.040 мс`.
