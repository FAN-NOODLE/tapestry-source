<?php
/**
 * 删除素材
 * User: tony
 * Date: 20200220
 * Time: 11:20
 */

namespace App\Http\Controllers\V1\CustomerMaintenance\Media;

use App\Models\QYwxMedia;
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
     * 删除素材
     *
     * @param string $scid
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function index($scid = "")
    {
        if (!$scid) notice(1100002);
        $data = [
            'deleted_at' => time(),
            'status'     => -1
        ];

        $where = [['scid', $scid]];
        if ($this->qyDomain) {
            $where[] = ['uc_eid', $this->qyDomain];
        }
        try {
            if (!QYwxMedia::query()->where($where)->update($data)) {
                error(1100005);
            }
            $id = $scid;
        } catch (\Exception $e) {
            notice(1100005);
        }
        return success(['id' => $id]);
    }
}
