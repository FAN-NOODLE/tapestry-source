<?php

namespace App\Models;

use App\Handlers\SoftDeletes;
use App\Handlers\SoftDeletingScope;
use App\Traits\IdGenerateTrait;
use Exception;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Database\Eloquent\Relations\HasMany;
use Illuminate\Database\Eloquent\Relations\HasOne;
use Illuminate\Database\Eloquent\Relations\Relation;

abstract class Base extends Model
{
    /**
     * 使用软删除
     */
    use SoftDeletes, IdGenerateTrait;

    /**
     * 设置日期时间格式
     *
     * @var string
     */
    protected $dateFormat = 'U';

    /**
     * 关闭主键ID自动转int类型（否则分布式ID插入后，返回的ID为0）
     *
     * @var bool
     */
    public $incrementing = false;

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

    protected $hidden = [
        'deleted_at'
    ];

    protected $guarded = [];

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

        if ($this->incrementing === false) {
            // 如果模型有主键,指定的主键值优先级高于这里的分布式ID赋值
            if ($this->primaryKey !== null) {
                if ($method === 'insert') {
                    // 判断是新增一条记录还是多条
                    if (! is_array(reset($parameters[0]))) {
                        if (!empty($parameters[0]) && !isset($parameters[0][$this->primaryKey])) {
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
        }

        return $this->newQuery()->$method(...$parameters);
    }

    public static function query()
    {
        return (new static);
    }

    /**
     * 此方法为Eloquent模型 关联关系补充WhereHasIn
     * 若使用WhereHas时，主表的数据过大导致查询效率降低，可该用此方法
     * @param Builder $builder
     * @param $relationName
     * @param callable $callable
     * @return object
     * @throws Exception
     */
    public function scopeWhereHasIn(Builder $builder, $relationName, callable $callable)
    {
        $relationNames = explode('.', $relationName);
        $nextRelation = implode('.', array_slice($relationNames, 1));

        $method = $relationNames[0];

        $relation = Relation::noConstraints(function () use ($method) {
            return $this->$method();
        });

        if($nextRelation){
            $in = $relation->getQuery()->whereHasIn($nextRelation, $callable);
        } else {
            $in = $relation->getQuery()->where($callable);
        }

        if ($relation instanceof BelongsTo) {
            return $builder->whereIn($relation->getForeignKey(), $in->select($relation->getOwnerKey()));
        } elseif ($relation instanceof HasOne) {
            return $builder->whereIn($this->getKeyName(), $in->select($relation->getForeignKeyName()));
        } elseif ($relation instanceof HasMany) {
            return $builder->whereIn($this->getKeyName(), $in->select($relation->getForeignKeyName()));
        }

        throw new Exception(__METHOD__ . " 不支持 " . get_class($relation));
    }
}
