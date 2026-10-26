<?php

namespace App\Services;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Library\Uc\Department;
use App\Models\StoreDepartment;
use App\Models\StoreEmployee;
use App\Models\StoreSftpErrorRecord;
use App\Models\StoreSftpLog;


class StoreService extends BaseService
{


    /**
     * 三方门店列表
     *
     * @param $params
     * @param bool $isPage
     * @param array $selectFields
     * @return array
     * User: zhaoyang@vchangyi.com
     * Date: 2022/2/7 15:33
     */
    public static function departmentList($params, bool $isPage = true, array $selectFields = ['*']): array
    {
        $page = $params['page'] ?? Constant::PAGE;
        $size = $params['size'] ?? Constant::SIZE;

        $sql   = StoreDepartment::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->when(!empty($params['code_or_name']), function ($query) use ($params) {
                $query->where(function ($query) use ($params) {
                    return $query->where('third_store_code', $params['code_or_name'])
                        ->orWhere('store_name', 'like', '%' . $params['code_or_name'] . '%');
                });
            })
            ->when(!empty($params['store_type']), function ($query) use ($params) {
                return $query->where('store_type', $params['store_type']);
            })
            ->when(!empty($params['close_status']), function ($query) use ($params) {
                return $query->where('closed_at', $params['close_status'] == 1 ? '=' : '>', 0);
            });
        $total = $sql->count();

        if ($isPage) {
            $sql = $sql->forPage($page, $size);
        }
        $list = $sql->orderBy('created_at', 'desc')
            ->get($selectFields)
            ->toArray();
        return [
            'list'  => $list,
            'total' => $total,
            'page'  => (int)$page,
            'size'  => (int)$size
        ];
    }

    /**
     * 三方门店员工列表
     *
     * @param $params
     * @param bool $isPage
     * @param array $selectFields
     * @return array
     * User: zhaoyang@vchangyi.com
     * Date: 2022/2/7 15:33
     */
    public static function employeeList($params, $isPage = true, $selectFields = ['*'])
    {
        $page = $params['page'] ?? Constant::PAGE;
        $size = $params['size'] ?? Constant::SIZE;

        $sql   = StoreEmployee::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->when(!empty($params['third_store_code']) && empty($params['no_department']), function ($query) use ($params) {
                return $query->where('third_store_code', $params['third_store_code']);
            })
            ->when(!empty($params['no_department']), function ($query) use ($params) {
                return $query->where('third_store_code', '')->where('status',StoreEmployee::STATUS_ENABLE);
            })
            ->when(empty($params['no_department']), function ($query) use ($params) {
                return $query->where('third_store_code', '!=', '');
            })
            ->when(!empty($params['mobile_or_name']), function ($query) use ($params) {
                $query->where(function ($query) use ($params) {
                    return $query->where('mobile', $params['mobile_or_name'])
                        ->orWhere('name', 'like', '%' . $params['mobile_or_name'] . '%')
                        ->orWhere('alias_name', 'like', '%' . $params['mobile_or_name'] . '%')
                        ->orWhere('employee_code', $params['mobile_or_name']);
                });
            })
            ->when(!empty($params['status']), function ($query) use ($params) {
                return $query->where('status', $params['status']);
            });
        $total = $sql->count();

        if ($isPage) {
            $sql = $sql->forPage($page, $size);
        }
        $list = $sql->orderBy('created_at', 'desc')
            ->get($selectFields)
            ->toArray();

        return [
            'list'  => $list,
            'total' => $total,
            'page'  => (int)$page,
            'size'  => (int)$size
        ];
    }


    /**
     * 修改门店信息
     *
     * @param $params
     * User: zhaoyang@vchangyi.com
     * Date: 2022/2/7 15:44
     */
    public static function saveDepartment($params)
    {
        $info = StoreDepartment::query()->find($params['id']);

        foreach ($params as $key => $val) {
            if (isset($info->$key) && $val != $info->$key) {
                $info->$key = $val;
                if ($key == 'store_name') {
                    $info->sync_status = 0;
                }
                if ($key == 'sa_check_role_status' && $val == StoreDepartment::SA_CHECK_ROLE_BY_SM) {
                    // 如果改为只允许店长查看 则直接清空其他人权限
                    $info->sa_check_role = '';
                }
                if ($key == 'sa_check_role' && !empty($val)) {
                    $info->sa_check_role = join(',', $val);
                }
                if ($key == "address" && !empty($val)) {
                    $res = file_get_contents('https://apis.map.qq.com/ws/geocoder/v1?key=AGFBZ-CM7L4-7F3UY-DY7NO-5P3KV-FGFIY&address=' . $val);
                    if (!empty($res)) {
                        $location = json_decode($res,true);
                        if (isset($location['result'])) {
                            if (isset($location['result']['location'])) {
                                $info->lat = $location['result']['location']['lat'] ?? '';
                                $info->lng = $location['result']['location']['lng'] ?? '';
                            }
                            if (isset($location['result']['address_components'])) {
                                $info->province_zh = $location['result']['address_components']['province'] ?? '';
                                $info->city_zh = $location['result']['address_components']['city'] ?? '';
                            }
                        }
                    }
                }
            }
        }

        $info->operate_name = $params['operate_name'];
        $info->save();
    }

    /**
     * 修改门店信息
     *
     * @param $params
     * User: zhaoyang@vchangyi.com
     * Date: 2022/2/7 15:44
     */
    public static function addDepartment($params)
    {
        $exists = StoreDepartment::query()->where("enterprise_id", $params['enterprise_id'])
            ->where(function ($query) use ($params) {
                $query->where("store_name", $params['store_name'])
                    ->orWhere("third_store_code", $params['store_code']);
            })->first();
        if ($exists) {
            throw new Error(600009, '门店名称或编码已存在');
        }
        $saveArr = [
            'store_name'            => $params['store_name'],
            'store_type'            => $params['store_type'],
            'third_store_code'      => $params['store_code'],
            'operate_name'          => $params['operate_name'],
            'enterprise_id'         => $params['enterprise_id'],
            'address'               => $params['address'] ?? '',
            'email'                 => $params['email'] ?? '',
            'city'                  => $params['city'] ?? '',
            'phone'                 => $params['phone'] ?? '',
            'inventory_name'        => $params['inventory_name'] ?? '',
            'popup_store_status'    => $params['popup_store_status'] ?? 0,
            'popup_store_keys'      => $params['popup_store_keys'] ?? '',
            'sa_share_order_url'    => $params['sa_share_order_url'] ?? '',
            'sa_check_order_url'    => $params['sa_check_order_url'] ?? '',
            'sa_check_role_status'  => $params['sa_check_role_status'] ?? 1,
            'sa_check_role'         => isset($params['sa_check_role']) ? join(',', $params['sa_check_role']) : '',
        ];
        StoreDepartment::query()->create($saveArr);
    }

    /**
     * 修改门店信息
     *
     * @param $params
     * @return bool
     * @throws \Exception
     */
    public static function deleteDepartment($params): bool
    {
        $deptInfo = StoreDepartment::query()->find($params['department_id']);
        if (empty($deptInfo)) return false;

        if (!empty($deptInfo->wechat_department_id)) {
            //如果已经同步企微则直接调用企微接口删除部门
            $res = (new Department())->delete($params['qy_domain'], 'communalContact', ['dpId' => $deptInfo->uc_department_id]);
            if (!is_array($res)) {
                $error = json_decode($res, true);
                throw new Error($error['code'], $error['msg']);
            }
        }

        //删除本地数据
        $deptInfo->delete();
        //更新关联员工数据
        StoreEmployee::query()->where('third_store_code', $deptInfo->third_store_code)
            ->update(['third_store_code' => '']);
        return false;
    }

    /**
     * 修改门店信息
     *
     * @param $params
     * @return void
     */
    public static function updateEmployeeDepartment($params)
    {
        $deptInfo = StoreDepartment::query()->find($params['department_id']);
        if (empty($deptInfo)) {
            throw new Error(600009, '门店不存在');
        }

        //更新关联员工数据
        StoreEmployee::query()->whereIn('id', $params['employee_id_arr'])
            ->where('enterprise_id', $params['enterprise_id'])
            ->update([
                'third_store_code' => $deptInfo->third_store_code,
                'sync_status'      => 0
            ]);

    }

    /**
     * 获取people hub 拉取文件记录
     * @param $params
     * @param $isPage
     * @param array $fields
     * @return array
     */
    public static function getSftpLogList($params, $isPage, $fields = ['*'])
    {
        $page  = $params['page'] ?? Constant::PAGE;
        $size  = $params['size'] ?? Constant::SIZE;
        $sql   = StoreSftpLog::query()->where('enterprise_id', $params['enterprise_id'])
            ->when(!empty($params['start_time']), function ($query) use ($params) {
                return $query->where('created_at', '>', $params['start_time']);
            })
            ->when(!empty($params['end_time']), function ($query) use ($params) {
                return $query->where('created_at', '<', $params['end_time']);
            })
            ->when(!empty($params['file_name']), function ($query) use ($params) {
                return $query->Where('file_name', 'like', '%' . $params['file_name'] . '%');
            })
            ->when(!empty($params['status']), function ($query) use ($params) {
                return $query->where('status', $params['status']);
            });
        $total = $sql->count();

        if ($isPage) {
            $sql = $sql->forPage($page, $size);
        }
        $list = $sql->select($fields)
            ->withCount(['errorRecord as error_count'])
            ->orderBy('created_at', 'desc')
            ->get()
            ->toArray();
        foreach ($list as &$val) {
            if ($val['error_count'] > 0 ) {
                $val['error_count'] --;
            }
        }
        return [
            'list'  => $list,
            'total' => $total,
            'page'  => (int)$page,
            'size'  => (int)$size
        ];
    }


    /**
     * 获取错误日志记录
     *
     * @param $fileLogId
     * @return array
     * User: zhaoyang@vchangyi.com
     * Date: 2022/3/23 16:46
     */
    public static function getErrorRecordList($fileLogId)
    {
        $fileInfo = StoreSftpLog::query()->find($fileLogId);
        if (empty($fileInfo)) {
            throw new Error(600009, '文件记录不存在');
        }
        $errorRecordList = StoreSftpErrorRecord::query()->where('file_log_id', $fileLogId)
            ->pluck('data_json')
            ->toArray();
        return ['file_name' => $fileInfo->file_name, 'error_record' => $errorRecordList];
    }


    /**
     * @param $params
     *
     * @return array
     * @author: wangle@vchangyi.com
     * @Time: 2022/5/18
     */
    public static function departmentDetail($params): array
    {

        $list = StoreDepartment::query()
            ->select("store_department.third_store_code", "store_department.wechat_department_id", "store_department.store_name", "store_department.address")
            ->leftJoin("store_employee", "store_employee.third_store_code", "=", "store_department.third_store_code")
            ->where('store_department.enterprise_id', $params['enterprise_id'])
            ->where('store_department.deleted_at', 0)
            ->when(!empty($params['mobile']), function ($query) use ($params) {
                return $query->where('store_employee.mobile', $params['mobile']);
            })->get();

        return !empty($list) ? $list->toArray() : [];
    }

}
