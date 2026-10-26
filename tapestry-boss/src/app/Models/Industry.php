<?php

namespace App\Models;

class Industry extends Base
{
    /**
     * 表名称:行业表
     *
     * @var string
     */
    protected $table = 'industry';

    protected $guarded = [];

    public function lists($page, $size = 10)
    {
        $total = self::query()->count('id');

        $list = self::query()->orderBy('id', 'DESC')->forPage($page, $size)->get()->toArray();

        return [$total, $list];
    }

    public function allData()
    {
        return self::query()
            ->orderBy('id', 'DESC')
            ->pluck('name', 'id')
            ->toArray();
    }

    /**
     * ID为key,名字为值
     * @param $industryIds
     * @return array
     */
    public function idKeyName($industryIds)
    {
        return self::query()
            ->select(['id', 'name'])
            ->whereIn('id', $industryIds)
            ->pluck('name', 'id')
            ->toArray();
    }

    /**
     * 验证ID是否存在
     * @param $id
     * @return bool
     */
    public static function idExists($id)
    {
        return self::query()
            ->where('id', '=', $id)
            ->exists();
    }
}
