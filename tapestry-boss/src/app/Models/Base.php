<?php

namespace App\Models;

use App\Common\RedisKey;
use App\Handlers\SoftDeletes;
use App\Handlers\SoftDeletingScope;
use App\Library\Json;
use App\Traits\IdGenerateTrait;
use Illuminate\Database\Eloquent\Model;

abstract class Base extends Model
{
    /**
     * 使用软删除
     */
    use SoftDeletes, IdGenerateTrait;

    /**
     * redis实例
     *
     * @var \Illuminate\Redis\Connections\Connection
     */
    protected $redis;

    /**
     * 默认缓存时间 秒
     *
     * @var int
     */
    protected $cacheExpire = 3600;

    /**
     * 关闭主键ID自动转int类型（否则分布式ID插入后，返回的ID为0）
     *
     * @var bool
     */
    public $incrementing = false;

    /**
     * 设置日期时间格式
     *
     * @var string
     */
    protected $dateFormat = 'U';

    /**
     * 需要被转换日期时间格式的字段
     *
     * @var array
     */
    protected $dates = [
        'created_at',
        'updated_at',
        'deleted_at'
    ];

    public function __construct(array $attributes = [])
    {
        parent::__construct($attributes);

        $this->redis = redis();
    }

    /**
     * 执行update语句时过滤where条件中的deleted_at字段，需要跟在update语句的where条件之后
     *
     * @param $query
     * @return mixed
     */
    public function scopeFilterDeletedAt($query)
    {
        return $query->withoutGlobalScope(SoftDeletingScope::class);
    }

    public function __call($method, $parameters)
    {
        if (in_array($method, ['increment', 'decrement'])) {
            return $this->$method(...$parameters);
        }

        // 如果模型有主键,指定的主键值优先级高于这里的分布式ID赋值
        if ($this->primaryKey !== null) {
            if ($method === 'insert') {
                // 判断是新增一条记录还是多条
                if (! is_array(reset($parameters[0]))) {
                    if (!isset($parameters[0][$this->primaryKey])) {
                        $parameters[0][$this->primaryKey] = idGenerate();
                    }
                } else {
                    foreach ($parameters[0] as $key => $parameter) {
                        if (!isset($parameters[0][$key][$this->primaryKey])) {
                            $parameters[0][$key][$this->primaryKey] = idGenerate();
                        }
                    }
                }
            }
            if (in_array($method, ['create', 'firstOrCreate', 'firstOrNew', 'updateOrCreate', 'insertGetId'])) {
                if (!isset($parameters[0][$this->primaryKey])) {
                    $parameters[0][$this->primaryKey] = idGenerate();
                }
            }
        }

        return $this->newQuery()->$method(...$parameters);
    }

    public static function query()
    {
        return (new static);
    }

    /**
     * 序列化
     *
     * @var array
     */
    protected $casts = [
        'created_at' => 'integer',
        'updated_at' => 'integer',
        'deleted_at' => 'integer'
    ];

    /**
     * 获取模型数据数组(优先查缓存)
     *
     * @param $where
     * @return \Illuminate\Database\Eloquent\Builder|Model|mixed|object|string|null
     */
    public function getModelArray($where)
    {
        $key = $this->getModelCacheKey($where);

        $data = $this->redis->get($key);

        if ($data === null) {
            $data = $this->queryModel($where);

            if ($data !== null) {
                $user = $data->toArray();

                $this->redis->setex($key, $this->cacheExpire, Json::encode($user));
            }
        } else {
            $data = Json::decode($data);
        }

        return $data;
    }

    /**
     * 清理模型缓存
     *
     * @param array $where
     * @return int
     */
    public function clearModelCache($where)
    {
        $key = $this->getModelCacheKey($where);

        return $this->redis->del([$key]);
    }

    /**
     * 查询模型数据
     *
     * @param $where
     * @return \Illuminate\Database\Eloquent\Builder|Model|object|null
     */
    public function queryModel($where)
    {
        // 说明有调用方对时效性要求较高,尝试清理掉缓存便于重新生成
        $key = $this->getModelCacheKey($where);

        $this->redis->del([$key]);

        return self::query()->where($where)->first();
    }

    /**
     * 获取生成模型缓存key
     *
     * @param array $where eg: [['open_id', '=', 1], ['enterprise_id', '=', 1]]
     * @return string eg: {prefix}:{table}:open_id=1:enterprise_id=1:deleted_at=0
     */
    private function getModelCacheKey($where)
    {
        $suffix = '';
        foreach ($where as $key => $value) {
            if (is_array($value)) {
                // eg: [['id', '=', 1]]
                if (is_array($value[2])) {
                    $value[2] = implode(',', $value[2]);
                }
                $suffix .= $value[0] . $value[1] . $value[2] . ':';
            } else {
                // eg: ['id' => 1]
                $suffix .= $key . '=' . $value . ':';
            }
        }
        $suffix = substr($suffix, 0, -1);

        return sprintf(RedisKey::MODEL_ARRAY, $this->table, $suffix);
    }
}
