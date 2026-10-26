<?php
/**
 * Created by PhpStorm.
 * User: wpf
 * Date: 2020/4/12
 * Time: 16:52
 */

namespace App\Services;

use App\Models\EnterpriseDataRange;

class EnterpriseDataRangeService extends BaseService
{
    /**
     *  1-区域  2-门店 3-商城
     */
    const TYPE_DISTRICT = 1;
    const TYPE_STORE = 2;
    const TYPE_SHOP = 3;

    /**
     * 插入账号范围数据
     * @param $input
     * @return bool
     */
    public function addEnterpriseDataRange($input)
    {
        // 插入
        $data = [];
        foreach ($input['range'] as $value) {
            $data[] = [
                'enterprise_id'       => $input['enterprise_id'],
                'enterprise_admin_id' => $input['enterprise_admin_id'],
                'type'                => $value['type'],
                'object_id'           => $value['object_id'],
                'created_at'          => time(),
                'updated_at'          => time(),
            ];
        }
        return EnterpriseDataRange::query()->insert($data);
    }

    /**
     * 编辑账号范围
     * @param $input
     * @return bool
     * @throws \Exception
     */
    public function editDataRange($input, $isDel = true)
    {
        if (!isset($input['range'])){
            return false;
        }

        // 获取旧数据
        $list      = $this->getDataRangeListByAdminId($input['id']);
        $objectIds = array_column($list, 'object_id');

        // 新提交数据
        $inputObjectIds = array_column($input['range'], 'object_id');

        $insertObjectIds = array_diff($inputObjectIds, $objectIds);

        // 取应当删除的数据差集
        $delObjectIds = array_diff($objectIds, $inputObjectIds);

        // 相互都没有差集 不需要修改
        if (empty($insertObjectIds) && empty($delObjectIds)) {
            return true;
        }

        // 插入增量数据
        if (!empty($insertObjectIds)) {
            $data = [];
            foreach ($input['range'] as $value) {
                if (in_array($value['object_id'], $insertObjectIds)) {
                    $data[] = [
                        'enterprise_id'       => $input['enterprise_id'],
                        'enterprise_admin_id' => $input['id'],
                        'type'                => $value['type'],
                        'object_id'           => $value['object_id'],
                        'created_at'          => time(),
                        'updated_at'          => time(),
                    ];
                }
            }
            EnterpriseDataRange::query()->insert($data);
        }

        // 删除减量数据-isdel时才删除减量数据
        if (!empty($delObjectIds) && $isDel) {
            $this->delDataRangeByAdminId($input['id'], $delObjectIds);
        }

        return true;
    }

    /**
     * 根据管理员ID获取数据范围
     * @param $adminId
     * @return array
     */
    public function getDataRangeListByAdminId($adminId)
    {
        return EnterpriseDataRange::query()
            ->where('enterprise_admin_id', $adminId)
            ->get()
            ->toArray();
    }

    /**
     * 根据管理员ID删除管理员所有账号范围数据
     * @param $adminId
     * @return bool|mixed|null
     * @throws \Exception
     */
    public function delDataRangeByAdminId($adminId, $objectIds = [])
    {
        return EnterpriseDataRange::query()
            ->when(!empty($objectIds), function ($query) use ($objectIds) {
                return $query->whereIn('object_id', $objectIds);
            })
            ->where('enterprise_admin_id', $adminId)
            ->delete();
    }

    /**
     * 获取范围列表数据
     * @param $enterprise_id
     * @param null $object_id
     * @return array
     */
    public function getDataRangeList($enterprise_id, $object_id = null, $adminIds = []) {
        return EnterpriseDataRange::query()
            ->when($object_id != null && is_array($object_id), function ($query) use ($object_id) {
                return $query->whereIn('object_id', $object_id);
            })
            ->when($object_id != null && is_numeric($object_id), function ($query) use ($object_id) {
                return $query->where('object_id', $object_id);
            })
            ->when(!empty($adminIds), function ($query) use ($adminIds) {
                return $query->whereIn('enterprise_admin_id', $adminIds);
            })
            ->where('enterprise_id', $enterprise_id)
            ->get()
            ->toArray();
    }
}
