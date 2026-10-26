<?php
/**
 * Created by PhpStorm.
 * Function: 企业员工业绩排行
 * User: 龚国新
 * Date: 2020/2/14
 * Time: 11:43
 */

namespace App\Http\Controllers\App\V1\StaffOrder;


use App\Http\Controllers\App\V1\BaseController;
use App\Models\StaffOrderRelation;
use App\Services\StaffManagerService;
use App\Services\StaffOrderService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;

class AchievementRankListController extends BaseController
{
    //  protected $requireLogin = false;

    /**
     * 企业员工排行榜显示top但是要显示当前内购用户排名
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $type   = $request->input('type', 0);
        $memUid = $request->input('mem_uid', '');
        $startTime = $request->input('start_time', '');
        $endTime   = $request->input('end_time', '');

        if (null == $this->memUid) {
            $this->memUid = $memUid;
        }

        // 如果不是自定时间
        if ($type != 4) {
            $search = [
                'type' => $type ?? 0,
            ];
            list($startTime, $endTime) = (new StaffOrderService())->orderTimeWhere($search);
        }


        $obj = StaffOrderRelation::query()
            ->select(DB::raw('sum(if(order_no=1,0,1)) as all_num,mem_uid,min(created_at) as time'))
            ->when(!empty($startTime), function ($query) use ($startTime) {
                return $query->where('created_at', '>=', $startTime);
            })
            ->when(!empty($endTime), function ($query) use ($endTime) {
                return $query->where('created_at', '<', $endTime);
            })
            ->where('enterprise_id', $this->enterpriseId)
            ->groupBy('mem_uid')
            ->orderBy('id')
            ->get()
            ->toArray();

        $data = $this->format($obj);

        return success($data);
    }

    /**
     * 格式化代码
     *
     * @param array $data
     * @param array $search
     * @param int   $order_num
     *
     * @return array
     */
    public function format($data = [], $order_num = 1)
    {
        if (empty($data)) {
            return ['list' => [], 'my_rank' => []];
        }

        $memUids = array_column($data, 'mem_uid');

        $memUids = array_unique($memUids);

        // 分页查询导购排行
        $chunkMemUids = array_chunk($memUids, 500);

        $userData = $list = $myRank = [];

        foreach ($chunkMemUids as $v) {
            sort($v);
            $params['memUids'] = $v;

            $userList = (new StaffManagerService())->getStaffListFromUc($this->qyDomain, $this->appIdentifier, $params);
            $userData = array_merge($userData, $userList['list']);
        }

        $userData = array_column($userData, null, 'memUid') ?? [];

        foreach ($data as $value) {
            $user = $userData[$value['mem_uid']] ?? [];

            $orderData = [
                'order_num'    => $order_num,
                'mem_username' => $user['memUsername'] ?? '',
                'mem_face'     => $user['memFace'] ?? '',
                'all_num'      => intval($value['all_num']) ?? 0,
                'time'         => $value['time']
            ];

            // 如果是当前用户则记录排名
            if ($this->memUid == $value['mem_uid']) {
                $myRank = $orderData;
            }

            // 如果大于50名则跳出
            if ($order_num <= 50) {
                $list[] = $orderData;
            }

            $order_num++;
        }

        array_multisort(array_column($list, 'all_num'), SORT_NUMERIC, SORT_DESC, array_column($list, 'time'), SORT_NUMERIC, SORT_ASC, $data);

        return ['data' => $list, 'my_rank' => $myRank];
    }
}