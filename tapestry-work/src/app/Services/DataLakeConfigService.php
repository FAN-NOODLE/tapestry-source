<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/26 0026
 * Time: 15:37
 */

namespace App\Services;


use App\Models\DataLakeConfig;
use Illuminate\Database\Eloquent\Builder;

class DataLakeConfigService
{
    public function __construct()
    {
        $this->dataLakeConfig = DataLakeConfig::query();
    }

    public function add($param)
    {
        $this->dataLakeConfig->create($param);
    }

    public function edit($param, $where)
    {
        $this->dataLakeConfig->where($where)->update($param);
    }

    public function delete($where)
    {
        $this->dataLakeConfig->where($where)->delete();
    }

    public function list($page, $size, $where = [])
    {
        $paganator = $this->dataLakeConfig
            ->when($where, function (Builder $query) use ($where) {
                return $query->where($where);
            })
            ->paginate($size, ['*'], $pageName = '', $page);
        $list = [];
        foreach ($paganator->items() as $item){
            /**
             * @var DataLakeConfig $item
             */
            $list[] = $item->toArray();
        }
        return [
            'page'=>$page,
            'size'=>$size,
            'total'=>$paganator->total(),
            'list'=>$list
        ];
    }
}