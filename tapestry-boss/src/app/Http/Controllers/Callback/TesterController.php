<?php

namespace App\Http\Controllers\Callback;

use App\Models\EnterpriseOpenPlatform;
use CrsSdk\Boss\Auth;
use CrsSdk\MicroAuthority\WX;
use CrsSdk\MicroGoodsV1\Goods;
use CrsSdk\Shop\Invoice;
use App\Exceptions\Error;
use App\Models\Enterprise;
use App\Traits\ToolsTrait;
use EasyWeChat\OpenPlatform\Application;
use GuzzleHttp\Client;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use CrsSdk\Shop\AutoRelease;
use App\Services\EnterpriseLogService;
use App\Models\EnterpriseThirdPartAccount;
use App\Models\EnterpriseMiniProgramAutoPack;
use App\Models\EnterpriseMiniProgramCodeLibrary;
use App\Services\OpenPlatform\MiniProgramService;
use EasyWeChat\Kernel\Exceptions\InvalidConfigException;
use Overtrue\LaravelWeChat\Facade as LaravelWechat;

class TesterController extends BaseController
{
    use ToolsTrait;

    protected $miniProgramService;

    /**
     * 是否将草稿提交至模版库
     */
    const ADD_YES = 1; // 添加
    const ADD_NO = 2; // 不添加

    /**
     * 微信对应的商品状态
     */
    const STATUS_UP = 1;      // 在售
    const STATUS_DOWN = 0;    // 下架
    const STATUS_DELETE = -1; // 页面对应的商品不存在

    /**
     * @var EnterpriseThirdPartAccount
     */
    protected $thirdPartAccount;

    /**
     * TesterController constructor.
     *
     * @param MiniProgramService $miniProgramService
     * @param EnterpriseThirdPartAccount $thirdPartAccount
     */
    public function __construct(MiniProgramService $miniProgramService, EnterpriseThirdPartAccount $thirdPartAccount)
    {
        $this->miniProgramService = $miniProgramService;
        $this->thirdPartAccount = $thirdPartAccount;
    }

    public function bindToOpenPlatform(Request $request, Application $openPlatform)
    {
        $appid = $request->input('app_id');
        $open_appid = $request->input('open_appid');
        $enterpriseId = $request->input('enterprise_id');

        $accessToken = (new WX())->accessToken(['enterprise_id' => $enterpriseId]);

        $client = new Client();
        try {
            $result = $client->request('POST', 'https://api.weixin.qq.com/cgi-bin/open/bind?access_token='REDACTED_FOR_SOURCE_SHARING'access_token'], [
                'headers'     => [
                    'Accept'       => 'application/json',
                    'Content-Type' => 'application/json',
                ],
                'body'        => json_encode([
                    'appid'      => $appid,
                    'open_appid' => $open_appid
                ]),
                'http_errors' => false
            ]);
            return json_decode($result->getBody(), true);
        } catch (\Exception $e) {
            throw new Error($e->getCode(), $e->getMessage());
        }
    }

    /**
     * 获取体验二维码
     *
     * @param Request $request
     * @throws InvalidConfigException
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function qrcode(Request $request)
    {
        $appId = $request->input('app_id');
        if (!$appId) {
            throw new Error(100000, '非正确的体验参数');
        }

        $this->miniProgramService->showExperienceQrcode($appId);
    }

    /**
     * 添加体验者.
     *
     * @param Request $request
     * @return array
     * @throws InvalidConfigException
     */
    public function createTester(Request $request)
    {
        $appId = $request->input('app_id');
        $wechatId = $request->input('wechat_id');
        $remark = $request->input('remark');
        $name = $request->input('name');

        if (!$appId || !$wechatId) {
            throw new Error(100000, '非正确的体验参数');
        }

        if (!$remark) {
            throw new Error(100000, '请填写体验者备注');
        }

        if (!$name) {
            throw new Error(100000, '请填写体验者姓名');
        }

        $this->miniProgramService->addTester($wechatId, $appId, ['name' => $name, 'remark' => $remark]);

        return [
            'data' => [
                'message' => 'success'
            ]
        ];
    }

    /**
     * 获取小程序分类.
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function category(Request $request)
    {
        $appId = $request->input('app_id');
        if (!$appId) {
            throw new Error(100000, '非正确的体验参数');
        }

        return $this->miniProgramService->getCategory($appId);
    }

    /**
     * 获取小程序页面.
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function page(Request $request)
    {
        $appId = $request->input('app_id');
        if (!$appId) {
            throw new Error(100000, '非正确的体验参数');
        }

        return $this->miniProgramService->getPage($appId);
    }

    /**
     * 提交小程序审核.
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function submitAudit(Request $request)
    {
        $appId = $request->input('app_id');
        if (!$appId) {
            throw new Error(100000, '非正确的体验参数');
        }

        $auditLogStatus = $this->miniProgramService->submitAudit($appId);

        return [
            'data' => [
                'status' => $auditLogStatus
            ]
        ];
    }

    /**
     * 发布指定版本
     *
     * @param Request $request
     * @return array
     * @throws InvalidConfigException
     */
    public function submitSpecified(Request $request)
    {
        $appId = $request->input('app_id');
        if (!$appId) {
            throw new Error(100000, '非正确的体验参数');
        }

        $version = $request->input('version');
        if (!$version) {
            throw new Error(100000, '非正确的版本参数');
        }

        $isChecked = $request->input('is_checked', false);
        $result = $this->miniProgramService->submitSpecifiedVersion($appId, $version, $isChecked);

        return [
            'data' => [
                'result' => $result
            ]
        ];
    }

    /**
     * 完整提交小程序审核流程
     *
     * @param Request $request
     * @return array
     * @throws Error
     */
    public function submitVersion(Request $request)
    {
        $appId = $request->input('app_id');
        // 是否进行版本提审
        $isChecked = $request->input('is_checked', false);

        if (!$appId) {
            throw new Error(100000, '非正确的体验参数');
        }

        $thirdPartAccount = $this->thirdPartAccount->where('app_id', $appId)->where('status', 1)->first();
        if (!$thirdPartAccount) {
            throw new Error(100000, '非正确的参数');
        }

        try {
            if (config('app.env') != 'prod') {
                $result = $this->miniProgramService->submitVersion($thirdPartAccount, $isChecked);
            } else {
                $result = $this->miniProgramService->submitVersion($thirdPartAccount, true);
            }
        } catch (\Exception $e) {
            app('log')->error($e);
            $result = false;
        }

        return [
            'data' => [
                'result' => $result
            ]
        ];
    }

    /**
     * 更新可信域名
     *
     * @param Request $request
     * @return array
     * @throws InvalidConfigException
     */
    public function setDomains(Request $request)
    {
        $appId = $request->input('app_id');
        if (!$appId) {
            throw new Error(100000, '非正确的体验参数');
        }

        $thirdPartAccount = $this->thirdPartAccount->where('app_id', $appId)->first();
        if (!$thirdPartAccount) {
            throw new Error(100000, '非正确的参数');
        }

        // 更新可信域名
        $result = $this->miniProgramService->setDomain($thirdPartAccount->app_id, 'set', $thirdPartAccount->refresh_token, $thirdPartAccount->enterprise_id);

        return $result;
    }

    /**
     * 根据app_id获取审核情况.
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function getAuditLog(Request $request)
    {
        $appId = $request->input('app_id');
        if (!$appId) {
            throw new Error(100000, '非正确的体验参数');
        }

        return $this->miniProgramService->recentAudit($appId);
    }

    /**
     * 根据审核ID获取审核情况.
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function getAuditIdLog(Request $request)
    {
        $appId = $request->input('app_id');
        $auditId = $request->input('audit_id');
        if (!$appId || !$auditId) {
            throw new Error(100000, '非正确的体验参数');
        }

        return $this->miniProgramService->getAuditInfoByWXAuditId($appId, $auditId);
    }

    /**
     * 上线小程序.
     *
     * @param Request $request
     * @return array
     * @throws InvalidConfigException
     */
    public function release(Request $request)
    {
        $appId = $request->input('app_id');
        $auditId = $request->input('audit_id');
        if (!$appId || !$auditId) {
            throw new Error(100000, '非正确的体验参数');
        }

        $res = $this->miniProgramService->release($appId, $auditId);

        return [
            'data' => [
                'status' => $res
            ]
        ];
    }

    /**
     * 上、下线小程序.
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function offline(Request $request)
    {
        $appId = $request->input('app_id');
        $action = $request->input('action');
        if (!$appId) {
            throw new Error(100000, '非正确的体验参数');
        }

        $action = $action ? $action : 'OFFLINE';
        $res = $this->miniProgramService->changeOnlineStatus($appId, $action);

        return [
            'data' => [
                'status' => $res
            ]
        ];
    }

    /**
     * 获取草稿箱第一个草稿
     *
     * @return mixed
     * @throws Error
     */
    public function draft()
    {
        $drafts = $this->miniProgramService->getDrafts();
        $drafts = multi_collect($drafts)->sortByDesc('create_time')->values();

        return $drafts->first();
    }

    /**
     * 获取模板列表.
     *
     * @return array
     * @throws Error
     */
    public function template()
    {
        return $this->miniProgramService->getTemplates();
    }

    /**
     * 撤回小程序审核.
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function cancelAudit(Request $request)
    {
        $appId = $request->input('app_id');
        if (!$appId) {
            throw new Error(100000, '非正确的体验参数');
        }

        return $this->miniProgramService->cancelAudit($appId);
    }

    /**
     * 批量添加体验者
     *
     * @param Request $request
     * @return array
     */
    public function allTester(Request $request)
    {
        $appId = $request->input('app_id');
        $infos = [
            [
                'wechat_id' => 'LR815648863',
                'name'      => '刘瑞',
                'remark'    => '刘瑞'
            ],
            [
                'wechat_id' => 'HLY9694',
                'name'      => '胡龙阳',
                'remark'    => '胡龙阳'
            ]
        ];

        $errors = [];
        foreach ($infos as $info) {
            try {
                $this->miniProgramService->addTester($info['wechat_id'], $appId, [
                    'name'   => $info['name'],
                    'remark' => $info['remark']
                ]);
            } catch (\Exception $e) {
                $errors = [
                    'msg'       => $e->getMessage(),
                    'wechat_id' => $info['wechat_id']
                ];
            }
        }

        return [
            'data' => [
                'message' => 'success',
                'errors'  => $errors
            ]
        ];
    }

    /**
     * 手动打包提审代码
     *
     * @param Request $request
     * @param MiniProgramService $miniProgramService
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function handPackAudit(Request $request, MiniProgramService $miniProgramService)
    {
        $rules = [
            'is_add'     => 'integer|required|in:1,2',
            'status'     => 'integer|required|in:2',
            'app_id'     => 'required',
            'wx_version' => 'required',
            'version_id' => 'required'
        ];
        $messages = [
            'is_add.*'     => 'is_add类型有误',
            'status.*'     => 'status不符合要求',
            'app_id.*'     => 'app_id类型有误',
            'wx_version.*' => 'wx_version类型有误',
            'version_id.*' => 'version_id类型有误'
        ];
        $this->validate($request, $rules, $messages);

        $input = $request->all();
        EnterpriseLogService::debug('手动执行打包提审代码');

        $packMsg = '';
        try {
            // 是否需要将草稿添加到模板库
            if ($input['is_add'] == self::ADD_YES) {
                //当前task对应版本号草稿提交至模版库
                $miniProgramService->addOneDraftToTemplate($input['wx_version']);
            }

            // 提审版本号对应模版
            $miniProgramService->submitSpecifiedVersion($input['app_id'], $input['wx_version'], true, EnterpriseMiniProgramCodeLibrary::ONLINE_TYPE_MANUALLY);
            $data = [
                'id'     => $input['version_id'],
                'status' => EnterpriseMiniProgramAutoPack::SHOP_VERSION_AUDIT_ING
            ];

            // 此处调用切换域名在submitSpecifiedVersion方法中
            // ----提审后 通知shop端 版本审核中状态
            (new AutoRelease())->saveStatus($data);
        } catch (\Exception $e) {
            EnterpriseLogService::error('自动化打包，异步处理， 提审失败', [
                'message' => $e->getMessage(),
                'file'    => $e->getFile(),
                'line'    => $e->getLine(),
                'data'    => $input
            ]);
            $data = [
                'id'     => $input['version_id'],
                'status' => EnterpriseMiniProgramAutoPack::SHOP_VERSION_AUDIT_ERROR,
                'reason' => $e->getMessage()
            ];
            $packMsg = $e->getMessage();
        }

        $toShopMsg = '';
        try {
            //提审后 通知shop端 版本审核中状态
            (new AutoRelease())->saveStatus($data);
        } catch (\Exception $e) {
            EnterpriseLogService::error('手动打包，shop通知失败', [
                'message' => $e->getMessage(),
                'line'    => $e->getLine(),
                'file'    => $e->getFile(),
                'data'    => $data
            ]);
            $toShopMsg = $e->getMessage();
        }

        return success([
            'pack_audit_msg' => $packMsg,
            'to_shop_msg'    => $toShopMsg
        ]);
    }

    /**
     * 微信事件回调
     *
     * @param Request $request
     * @return bool|JsonResponse
     */
    public function wxEvent(Request $request)
    {
        $str = file_get_contents('php://input') ?? '';
        $goodsInput = json_decode($str, true);
        $input = $this->xmlToArray($str);
        $goodsInputType = $request->input('type', '');

        if (!empty($_GET['echostr']) || $goodsInputType == 'merchant') {
            if (!$this->checkWxCallback()) {
                EnterpriseLogService::error('微信事件回调验证失败', [$str]);
                return success(['data' => '验证失败']);
            }
        }
        $result = ['data' => ''];
        // 发票授权回调
        if (!empty($input) && !empty($input['Event']) && $input['Event'] == 'user_authorize_invoice') {
            $order = $input['SuccOrderId'] ?? '';
            if (empty($order)) {
                $order = is_array($input['SuccOrderId']) ? $input['SuccOrderId'][0] : $input['SuccOrderId'];
            }
            $enterpriseId = explode('_', $order)[0];
            $enterpriseInfo = Enterprise::query()->where('id', $enterpriseId)->first(['hash']);
            $_SERVER['enterprise_hash'] = $enterpriseInfo->hash ?? '';

            // 切换SHOP域名
            switchAppHost($enterpriseId);

            EnterpriseLogService::info('开票授权回调通知Shop', $input);
            $result = (new Invoice())->auth($input);
            EnterpriseLogService::info('开票授权回调通知Shop结果', $result);
        }

        // 更新商品数据
        try {
            if (!empty($goodsInput) && !empty($goodsInput['type']) && $goodsInput['type'] == 'merchant') {
                $appId = $goodsInput['appid'];
                $thirdPartAccount = EnterpriseThirdPartAccount::query()
                    ->where('app_id', '=', $appId)
                    ->where('status', '=', EnterpriseThirdPartAccount::STATUS_ACTIVE)
                    ->first(['enterprise_id']);
                if (empty($thirdPartAccount)) {
                    EnterpriseLogService::error('微信回调更新商品数据失败', $goodsInput);
                    return true;
                }
                $enterpriseData = Enterprise::query()->where(['id' => $thirdPartAccount->enterprise_id])->first(['hash']);
                $_SERVER['enterprise_hash'] = $enterpriseData->hash ?? '';

                // 切换SHOP域名
                switchAppHost($thirdPartAccount->enterprise_id);

                $spuIds = [];
                foreach ($goodsInput['req_data_list'] as $value) {
                    parse_str($value['query'], $queryArr);
                    if (isset($queryArr['id'])) {
                        $spuIds[] = $queryArr['id'];
                    }
                }

                if (empty($spuIds)) {
                    return success(['data' => '未查询到商品']);
                }

                $spuData = (new Goods())->spuDetail([
                    'enterprise_id' => $thirdPartAccount->enterprise_id,
                    'spu_ids'       => $spuIds
                ]);

                if (empty($spuData['list'])) {
                    return success(['data' => '未查询到商品']);
                }

                $resp_data_list = [];
                foreach ($spuData['list'] as $value) {
                    if (empty($value['sku'])) {
                        continue;
                    }
                    foreach ($value['sku'] as $sku) {
                        // 没上架就跳过
                        if ($sku['status'] !== 2) {
                            continue;
                        }
                        $resp_data_list[] = [
                            'path'    => 'pages/goods-detail/goods-detail',
                            'query'   => 'id=' . $value['id'] . '&sku_id=' . $sku['id'],
                            'uniq_id' => "{$sku['id']}",
                            'status'  => self::STATUS_UP,
                            'price'   => $sku['on_line_price'] / 100
                        ];
                        // 结束当前商品的循环
                        break;
                    }
                }
                $result = [
                    'type'           => 'merchant',
                    'resp_data_list' => $resp_data_list
                ];
            }
        } catch (\Exception $e) {
            EnterpriseLogService::error('微信回调更新商品数据失败', [
                'msg'          => $e->getMessage(),
                'line'         => $e->getLine(),
                'file'         => $e->getFile(),
                'origin_input' => $str
            ]);
        }

        return success($result);
    }

    /**
     * 设置公众号回调需要验证token
     *
     * @param Request $request
     * @param MiniProgramService $miniProgramService
     * @return array
     */
    private function checkWxCallback()
    {
        $env = config('app.env');
        $token = config('wechat.official.check_token.' . $env);

        $signature = $_GET["signature"];
        $timestamp = $_GET["timestamp"];
        $nonce = $_GET["nonce"];

        $tmpArr = array($token, $timestamp, $nonce);
        sort($tmpArr, SORT_STRING);
        $tmpStr = implode($tmpArr);
        $tmpStr = sha1($tmpStr);

        if ($tmpStr == $signature) {
            if (!empty($_GET['echostr'])) {
                $echoStr = $_GET["echostr"];
                echo $echoStr;
                die();
            } else {
                return true;
            }
        } else {
            return false;
        }
    }

    /**
     * 更新小程序配置域名
     *
     * @param Request $request
     * @return JsonResponse|void
     * @throws \Illuminate\Validation\ValidationException
     * @date 2020/12/24
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public function updateDomainJson(Request $request)
    {
        $params = $request->only([
            'id',
            'key',
            'value',
            'reset'
        ]);
        $rules = [
            'id'    => 'integer|required',
            'key'   => 'required|in:request_domain,wx_request_domain,upload_domain,download_domain,biz_domain',
            'value' => 'required',
        ];
        $messages = [
            'id.*'    => 'id类型有误',
            'key.*'   => 'key不符合要求',
            'value.*' => 'value类型有误',
        ];
        $this->validate($request, $rules, $messages);

        $config = EnterpriseOpenPlatform::query()->where('id', $params['id'])->where('deleted_at', 0)->first();

        if (empty($config)) {
            return notice(620004);
        }

        $domainArr = $old = json_decode($config->domain_json, true);
        $domainArr[$params['key']] .= ',' . $params['value'];
        // 重置域名
        if (isset($params['reset']) && $params['reset'] == 1) {
            $domainArr[$params['key']] = $params['value'];
        }
        $data = [
            'domain_json' => stripslashes(json_encode($domainArr)),
        ];

        EnterpriseOpenPlatform::query()->where('id', $params['id'])->update($data);

        foreach ($old as $key => &$value) {
            $value = explode(',', $value);
            $domainArr[$key] = explode(',', $domainArr[$key]);
        }
        return success([
            'old' => $old,
            'new' => $domainArr
        ]);
    }

    public function deployOpenConfig(Request $request)
    {
        $params = $request->only([
            'action',
            'app_id',
            'secret',
            'token',
            'aes_key',
            // 小程序打包发版配置信息
            'ext_build_json',
            // 小程序域名配置
            'domain_json'
        ]);
        $rule = [
            'action'         => 'required|in:add,update',
            'app_id'         => 'required',
            'ext_build_json' => 'array',
            'domain_json'    => 'array'
        ];
        $message = [
            'action.*' => 'key不符合要求',
            'app_id.*' => 'app_id不能为空'
        ];
        if ($params['action'] == 'add') {
            $rule = array_merge($rule, [
                'secret'  => 'required',
                'token'   => 'required',
                'aes_key' => 'required',
            ]);

            $message = array_merge($message, [
                'secret.*'  => 'secret不能为空',
                'token.*'   => 'token不能为空',
                'aes_key.*' => 'aes_key不能为空',
            ]);
        }
        $this->validate($request, $rule, $message);

        $config = EnterpriseOpenPlatform::query()->where('app_id', $params['app_id'])->where('deleted_at', 0)->first();

        if (empty($config)) {
            $data = [
                'app_id'         => $params['app_id'],
                'secret'         => $params['secret'],
                'token'          => $params['token'],
                'aes_key'        => $params['aes_key'],
                'ext_build_json' => empty($params['ext_build_json']) ? '' : stripslashes(json_encode($params['ext_build_json'])),
                'domain_json'    => empty($params['domain_json']) ? '' : stripslashes(json_encode($params['domain_json']))
            ];
            EnterpriseOpenPlatform::query()->create($data);
        } else {
            $data = [];
            if (!empty($params['ext_build_json'])) {
                $data['ext_build_json'] = json_encode($params['ext_build_json']);
            }

            if (!empty($params['domain_json'])) {
                $data['domain_json'] = json_encode($params['domain_json']);
            }

            if (empty($data)) {
                return success([
                    'msg' => '没有要更新的数据'
                ]);
            }
            EnterpriseOpenPlatform::query()->where('id', $config->id)->update($data);

            // 更新所有关联企业的build_json
            if (!empty($params['ext_build_json'])) {
                $enterpriseIds = Enterprise::query()->where('open_platform_id', $config->id)->pluck('id')->toArray();
                EnterpriseThirdPartAccount::query()
                    ->whereIn('enterprise_id', $enterpriseIds)
                    ->update([
                        'ext_build_json' => $data['ext_build_json'],
                    ]);
            }
        }

        $config = EnterpriseOpenPlatform::query()->where('app_id', $params['app_id'])->where('deleted_at', 0)->first();

        return success([
            'id' => $config->id,
        ]);

    }
}
