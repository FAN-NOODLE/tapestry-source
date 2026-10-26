<?php
/**获取内容
 * Created by PhpStorm.
 * User: wangpengfei
 * Date: 2019/5/6
 * Time: 9:54 AM
 */

namespace App\Http\Controllers\V1\ContentManage;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\ContentManage;

class InfoController extends BaseController
{

    public function index()
    {
        // 接值
        $params = $this->request->only([
            'type',
        ]);

        $rules    = [
            'type'          => 'required',
        ];
        $messages = [
            'type.required'          => '类型不能为空',
        ];
        $this->validate($this->request, $rules, $messages);

        $params['enterprise_id'] = $this->enterpriseId;

        return success((new ContentManage())->info($params));
    }
}
