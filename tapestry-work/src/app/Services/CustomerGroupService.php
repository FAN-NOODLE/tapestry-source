<?php


namespace App\Services;

use App\Common\Constant;
use App\Library\Cmq\CmqQueue;
use App\Models\CustomerGroup;
use App\Models\CustomerGroupRelation;
use Illuminate\Support\Facades\Log;

class CustomerGroupService
{

    /**
     * 插入用户分组数据
     * @param $input
     * @return array
     */
    public function insertCustomerGroup($input)
    {
        $creat = (new CustomerGroup())->insertCustomerGroup($input);

        $data = [
            'id'            => (int)$creat['id'],
            'name'          => $creat['name'],
            'enterprise_id' => (int)$creat['enterprise_id'],
            'user_visit'    => json_decode($input['user_visit'], true),
            'coupon'        => json_decode($input['coupon'], true),
            'order'         => json_decode($input['order'], true),
        ];

        // 发送队列消息
        $this->sendCMQMsg($data);

        return $data;
    }

    /**
     * 获取用户分组列表
     * @param $input
     * @return array
     */
    public function ListCustomerGroup($input)
    {
        $list = CustomerGroup::query()
            ->when(($input['keyword'] != '' && ($input['select'] == 0 || $input['select'] == 1)), function ($query) use ($input) {
                $query->where(function ($query) use ($input) {
                    $query->where('id', 'like', '%' . $input['keyword'] . '%')
                        ->orWhere('name', 'like', '%' . $input['keyword'] . '%');
                });
            })
            ->where('enterprise_id', $input['enterprise_id'])
            ->orderBy('id')
            ->forpage($input['page'] ?? Constant::PAGE, $input['size'] ?? Constant::SIZE)
            ->get(['id', 'name', 'user_visit', 'coupon', 'order', 'created_at'])
            ->toArray();

        if (empty($list)) {
            return [];
        }

        // 计算分组下人数
        if ($input['select'] == 0) {
            $groupIds = array_column($list, 'id');

            $countList = CustomerGroupRelation::query()
                ->selectRaw('count(id) as user_count, customer_group_id')
                ->whereIn('customer_group_id', $groupIds)
                ->groupBy(['customer_group_id'])
                ->get()
                ->toArray();

            $countList = array_column($countList, null, 'customer_group_id');

            foreach ($list as &$item) {
                $item['user_count'] = !empty($countList) ? (isset($countList[$item['id']]) ? $countList[$item['id']]['user_count'] : 0) : 0;
            }

            unset($item);
        }

        return $list;
    }

    /**
     * 获取分组记录总数
     * @param $input
     * @return int
     */
    public function customerGroupTotal($input)
    {
        $count = CustomerGroup::query()
            ->when(($input['keyword'] != '' && ($input['select'] == 0 || $input['select'] == 1)), function ($query) use ($input) {
                $query->where(function ($query) use ($input) {
                    $query->where('id', 'like', '%' . $input['keyword'] . '%')
                        ->orWhere('name', 'like', '%' . $input['keyword'] . '%');
                });
            })
            ->where('enterprise_id', $input['enterprise_id'])
            ->count();
        return $count ?? 0;
    }

    /**
     * 获取用户分组详情
     * @param $id
     * @return array
     */
    public function customerGroupDetail($id)
    {
        $detail = CustomerGroup::query()
            ->where('id', $id)
            ->first();

        return $detail ? $detail->toArray() : [];
    }

    /**
     * 根据用户名字获取分组详情-排除名字重复
     * @param $name
     * @param int $id
     * @return array
     */
    public function customerGroupDetailByName($name, $id = 0)
    {
        $detail = CustomerGroup::query()
            ->where('name', $name)
            ->when(!empty($id), function ($query) use ($id) {
                $query->where('id', '<>', $id);
            })
            ->first();

        return $detail ? $detail->toArray() : [];
    }

    /**
     * 根据条件获取详情-排除条件重复
     * @param $input
     * @return array
     */
    public function customerGroupDetailByCondition($input)
    {
        $detail = CustomerGroup::query()
            ->where([
                ['user_visit', $input['user_visit']],
                ['coupon', $input['coupon']],
                ['order', $input['order']]
            ])
            ->first();

        return $detail ? $detail->toArray() : [];
    }

    /**
     * 删除用户分组
     * @param $id
     * @return mixed
     */
    public function customerGroupDel($id)
    {
        return CustomerGroup::query()->where('id', $id)->delete();
    }

    /**
     * 更新用户分组
     * @param $input
     * @return int
     */
    public function customerGroupEdit($input)
    {
        $data = [
            'name'       => $input['name'],
            'user_visit' => $input['user_visit'],
            'coupon'     => $input['coupon'],
            'order'      => $input['order'],
        ];

        $res = CustomerGroup::query()->where('id', $input['id'])->update($data);

        // 发送队列消息
        $data['enterprise_id'] = (int)$input['enterprise_id'];
        $data['id']            = (int)$input['id'];
        $data['user_visit']    = json_decode($input['user_visit'], true);
        $data['coupon']        = json_decode($input['coupon'], true);
        $data['order']         = json_decode($input['order'], true);
        $this->sendCMQMsg($data);

        return $res;
    }

    public function sendCMQMsg($data, $action = 'customer_group')
    {
        $dataCMQ['source'] = 'crs';
        $dataCMQ['action'] = $action;
        $dataCMQ['query']  = $data;

        try {

            (new CmqQueue(config('cmq_setting.tag_customer_group')))
                ->sendMessage(json_encode($dataCMQ));

        } catch (\Exception $exception) {

            Log::info('发送队列action:' . $action . ',错误信息:' . $exception->getMessage());

        }

    }

}