<?php
/**
 * 删除推送消息
 * User: tony
 * Date: 20200220
 * Time: 11:20
 */

namespace App\Http\Controllers\V1\CustomerMaintenance\Msg;

use App\Models\QYwxEmsg;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;

class DeleteController extends BaseController
{
    public function __construct(Request $request)
    {
        parent::__construct();
        $this->request = $request;
    }

    /**
     * 删除推送消息
     *
     * @param string $tsid
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function index($tsid = "")
    {
        if (!$tsid) notice(1000005);

        $where = [
            ['tsid', $tsid],
            ['status', 0]
        ];
        if ($this->qyDomain) {
            $where[] = ['uc_eid', $this->qyDomain];
        }
        try {
            if (!QYwxEmsg::query()->where($where)->delete()) {
                error(1000006);
            }
            $id = $tsid;
        } catch (\Exception $e) {
            notice(1000006);
        }
        return success(['id' => $id]);
    }
}
