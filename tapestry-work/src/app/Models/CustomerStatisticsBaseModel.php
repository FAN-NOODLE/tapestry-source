<?php


namespace App\Models;


use Illuminate\Database\Eloquent\Model;
use Illuminate\Support\Facades\DB;
use phpDocumentor\Reflection\Types\Static_;

class CustomerStatisticsBaseModel extends Model
{

    /**
     * 设置日期时间格式
     *
     * @var string
     */
    public $dateFormat = 'U';

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
     * 统计字段列表
     * @var array
     */
    protected static $statFields = [];


    /**
     * 递增/递减更新统计字段
     * @param array|int $ids
     * @param array $setData
     * @param array $additionalWhere
     * @return int
     */
    public static function incrStatFieldsByIds($ids, array $setData, array $additionalWhere = []): int
    {
        $query = static::query();
        if (is_array($ids)) {
            if (count($ids) == 1) {
                $query->where('id', $ids[0]);
            } else {
                $query->whereIn('id', $ids);
            }
        } else {
            $query->where('id', $ids);
        }
        if (!empty($additionalWhere)) {
            $query->where($additionalWhere);
        }
        return $query->update(self::generateStatFieldsIncrSetData($setData));
    }

    /**
     * 根据更新数据生成统计字段递增/递减更新数据
     * @param array $setData
     * @return array
     */
    public static function generateStatFieldsIncrSetData(array $setData)
    {
        if (!isset($setData['updated_at'])) {
            $setData['updated_at'] = time();
        }
        $values = [];
        foreach ($setData as $field => $value) {
            if (in_array($field, static::$statFields)) {
                $values[$field] = DB::raw("{$field} + {$value}");
            } else {
                $values[$field] = $value;
            }
        }
        return $values;
    }

    /**
     * 生成统计字段聚合查询数据
     * @param array $statFields
     * @return array
     */
    public static function generateStatFieldsAggColumns($statFields = [])
    {
        if (empty($statFields)) {
            $statFields = static::$statFields;
        }
        $columns = [];
        foreach ($statFields as $statField) {
            array_push($columns, DB::raw("sum( $statField ) AS {$statField}"));
        }
        return $columns;
    }

    /**
     * 统计数据查询结果格式化
     * @param array $aggStatisticList
     */
    public static function formatStatColumn(array &$aggStatisticList)
    {
        foreach ($aggStatisticList as $key => $statistic) {
            foreach ($statistic as $field => $value) {
                if (in_array($field, static::$statFields)) {
                    if ($value === null) {
                        $aggStatisticList[$key][$field] = 0;
                    } else {
                        $aggStatisticList[$key][$field] = intval($value);
                    }
                }
            }
        }
    }
}