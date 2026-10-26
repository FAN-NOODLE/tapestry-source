<?php
/**
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/8/31
 * Time: 16:21
 */

namespace App\Http\Controllers\Callback;

use App\Cache\EnterpriseCache;
use App\Exceptions\Error;
use App\Models\Enterprise;
use Illuminate\Http\Request;

class ContactController
{
    /**
     * 企业model
     *
     * @var Enterprise
     */
    protected $enterprise;

    /**
     * @var EnterpriseCache
     */
    protected $enterpriseCache;

    /**
     * ContactController constructor.
     *
     * @param Enterprise      $enterprise
     * @param EnterpriseCache $enterpriseCache
     */
    public function __construct(Enterprise $enterprise, EnterpriseCache $enterpriseCache)
    {
        $this->enterprise = $enterprise;
        $this->enterpriseCache = $enterpriseCache;
    }

    /**
     * 回调入口
     *
     * @param Request $request
     *
     * @return mixed
     * @author 陈朔 chenshuo@vchangyi.com
     * @date   2018/12/29 4:00 PM
     */
    public function index(Request $request)
    {
        $enterpriseId = $request->get('enterprise_id');
        $corpId = $request->get('corpid');

        if (empty($enterpriseId)) {
            throw new Error(100000, '企业ID不能为空');
        }

        if (empty($corpId)) {
            throw new Error(100000, 'corpID不能为空');
        }

        // 去uc根据corp_id查询企业信息
        $ucEnterprise = (new \CrsSdk\Work\Enterprise())->detailByCorpId(
            ['corp_id' => $corpId]
        );

        if (!is_array($ucEnterprise) || empty($ucEnterprise['epEnumber'])) {
            throw new Error(100000, 'UC返回的企业数据错误');
        }

        $this->enterprise->where('id', $enterpriseId)->update([
            'corp_id' => $corpId,
            'domain'  => $ucEnterprise['epEnumber'],
            'contact_status' => 1 //通讯录已授权
        ]);

        // 此时需要更新一下企业的缓存 修复时间: 2018-12-29
        $this->enterpriseCache->update()->detail('id', (string)$enterpriseId);

        $enterpriseInfo = $this->enterprise->where('id', $enterpriseId)->first();
        $hash = '';
        if ($enterpriseInfo) {
            if ($enterpriseInfo->is_split_flow == 1) {
                // 分流企业
                $hash = $enterpriseInfo->hash;
            } else {
                // 不分流企业
                $hash = 'public';
            }
        }
        // 跳转至企业管理后台
        return redirect(config('app.console_host') . '/' . $hash . '/myapp/organization/setting?domain='
            . $ucEnterprise['epEnumber']);
    }
}
