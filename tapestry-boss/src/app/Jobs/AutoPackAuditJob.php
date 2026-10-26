<?php
/**
 * 小程序自动提审
 */

namespace App\Jobs;

use CrsSdk\Shop\AutoRelease;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Mail;
use App\Services\EnterpriseLogService;
use App\Models\EnterpriseMiniProgramAutoPack;
use App\Models\EnterpriseMiniProgramCodeLibrary;
use App\Services\OpenPlatform\MiniProgramService;

class AutoPackAuditJob extends Job
{
    /**
     * 队列名称
     *
     * @var string
     */
    public $queue = 'auto_pack_audit';

    /**
     * 参数
     *
     * @var $data
     */
    protected $data;

    public function __construct($data)
    {
        $this->data = $data;
    }

    /**
     * @param MiniProgramService $miniProgramService
     *
     * @return string
     */
    public function handle(MiniProgramService $miniProgramService)
    {
        Log::info('自动化打包，进入异步处理：' . json_encode($this->data));
        EnterpriseLogService::debug('自动化打包，进入异步处理', $this->data);
        if (isset($this->data['hash'])) {
            $_SERVER['enterprise_hash'] = $this->data['hash'];
            // 切换SHOP域名
            switchAppHost($this->data['hash']);
        }
        // 需要更新的数据
        $updateData = [
            'status' => $this->data['status'],
            'desc'   => is_string($this->data['statusDesc']) ? $this->data['statusDesc'] : json_encode($this->data['statusDesc'])
        ];
        switch ($this->data['status']) {
            case EnterpriseMiniProgramAutoPack::TASK_STATUS_FAILED:
                $updateData = $this->packFailed($updateData);
                break;
            case EnterpriseMiniProgramAutoPack::TASK_STATUS_NEED_LOGIN:
                $updateData = $this->packLogin($updateData);
                break;
            case EnterpriseMiniProgramAutoPack::TASK_STATUS_SUCCESS:
                $updateData = $this->packSuccess($updateData, $miniProgramService);
                break;
        }
        // 更新打包任务task状态
        EnterpriseMiniProgramAutoPack::query()
            ->where('task_id', '=', $this->data['id'])
            ->update($updateData);
        Log::info('自动化打包，异步处理 打包完成');
        // job貌似在内存中 全局变量不回回收 每次执行完手动清除
        if (isset($_SERVER['enterprise_hash'])) {
            unset($_SERVER['enterprise_hash']);
        }
        return 'success';
    }

    /**
     * 打包失败处理方法
     *
     * @param $updateData
     *
     * @return mixed
     */
    public function packFailed(&$updateData)
    {
        $this->sendSignalToShop([
            'id'     => $this->data['version_id'],
            'status' => EnterpriseMiniProgramAutoPack::SHOP_VERSION_AUDIT_ERROR,
            'reason' => $updateData['desc'] ?? '小程序打包失败，请重新提审'
        ]);
        Log::info('自动化打包，异步处理 小程序打包失败');
        return $updateData;
    }

    /**
     * @param                    $updateData
     * @param MiniProgramService $miniProgramService
     *
     * @return mixed
     */
    public function packSuccess(&$updateData, $miniProgramService)
    {
        $flag    = false;
        $message = '';
        try {
            //当前task对应版本号草稿提交至模版库
            $miniProgramService->addOneDraftToTemplate($this->data['wx_version']);
            // 提审版本号对应模版
            $miniProgramService->submitSpecifiedVersion(
                $this->data['app_id'],
                $this->data['wx_version'],
                true,
                EnterpriseMiniProgramCodeLibrary::ONLINE_TYPE_MANUALLY
            );
            $data = [
                'id'     => $this->data['version_id'],
                'status' => EnterpriseMiniProgramAutoPack::SHOP_VERSION_AUDIT_ING
            ];
            Log::info('自动化打包，异步处理 提审成功');
        } catch (\Exception $e) {
            Log::info('自动化打包，异步处理， 提审失败：' . $e->getMessage());

            EnterpriseLogService::error('自动化打包，异步处理， 提审失败', [
                'message' => $e->getMessage(),
                'file'    => $e->getFile(),
                'line'    => $e->getLine(),
                'data'    => $this->data
            ]);
            $data = [
                'id'     => $this->data['version_id'],
                'status' => EnterpriseMiniProgramAutoPack::SHOP_VERSION_AUDIT_ERROR,
                'reason' => $e->getMessage()
            ];

            $flag    = true;
            $message = $e->getMessage();
        }
        $updateData['wx_version'] = $this->data['wx_version'];
        $updateData['wx_desc']    = $this->data['wx_desc'];

        $this->sendSignalToShop($data);
        // 邮件通知相关人员
        $this->sendAuditMail($flag, $message);
        return $updateData;
    }

    /**
     * @param $data
     *
     * @return bool
     */
    public function sendSignalToShop($data)
    {
        try {
            // 存在原因
            if (isset($data['reason'])) {

                if (is_array($data['reason'])) {
                    $reason = json_encode($data['reason'], JSON_UNESCAPED_UNICODE);
                } else {
                    $reason = $data['reason'];
                }

                if (strpos($reason, 'wx2b03c6e691cd7370')) {
                    $data['reason'] = '提审失败,请关闭直播功能,重新提审';
                }

                if (strpos($reason, 'wx56c8f077de74b07c')) {
                    $data['reason'] = '提审失败,请关闭好物圈功能,重新提审';
                }

                if (strpos($reason, 'wxf3f436ba9bd4be7b')) {
                    $data['reason'] = '提审失败,请关闭发券插件功能,重新提审';
                }

                if (strpos($reason, 'Operation timed')) {
                    $data['reason'] = '请求微信超时, 请重新提审';
                }

                if (strpos($reason, 'git fetch origin')) {
                    $data['reason'] = '提审失败,联系开发者处理' . $reason;
                }
            }

            //提审后 通知shop端 版本审核中状态
            (new AutoRelease())->saveStatus($data);
        } catch (\Exception $e) {
            Log::info('自动化打包，异步处理， shop通知失败：' . $e->getMessage());
            EnterpriseLogService::error('自动化打包，异步处理，shop通知失败', [
                'message' => $e->getMessage(),
                'line'    => $e->getLine(),
                'file'    => $e->getFile(),
                'data'    => $data
            ]);
            $this->sendAuditMail(true, $e->getMessage());
        }
        Log::info('自动化打包，异步处理， 通知shop完成');
        return true;
    }

    /**
     * 打包需要登录处理方法
     *
     * @param $updateData
     *
     * @return mixed
     */
    public function packLogin(&$updateData)
    {
        $updateData['login_url'] = isset($this->data['url']) ? $this->data['url'] : '';
        // 发送邮件给 侯英才 马雄伟
        $this->sendLoginMail($this->data['url']);
        Log::info('自动化打包，异步处理 build打包平台需要登录');
        return $updateData;
    }

    private function sendLoginMail($url)
    {
        // 发布成功邮件通知相关人员
        $content = "

                Dear all，

                build自动打包提示：开发者工具需要登录

                点击链接 扫描二维码->：" . $url . "
                ps： 请在微信中打开，此链接打开一次即失效：

                请知悉，谢谢
            ";


        $originConfig = config('client.api_url.shop');
        if (strpos($originConfig, 'gw-shop') == false) {
            Mail::raw($content, function ($message) {
                $message->subject('自动打包平台登录请求（' . date('Y-m-d') . ')');
                $message->to("gongguoxin@vchangyi.com");
                $message->cc([
                    "maxiongwei@vchangyi.com",
                    "wangliping@vchangyi.com",
                    "meijie@vchangyi.com",
                    "luoweiqing@vchangyi.com",
                    "liyunxiang@vchangyi.com",
                ]);
            });
        } else {
            Mail::raw($content, function ($message) {
                $message->subject('自动打包平台登录请求（' . date('Y-m-d') . ')');
                $message->to("houyingcai@vchangyi.com");
                $message->cc([
                    "maxiongwei@vchangyi.com",
                    "wangliping@vchangyi.com",
                    "meijie@vchangyi.com",
                    "luoweiqing@vchangyi.com",
                    "liyunxiang@vchangyi.com",
                ]);
            });
        }


        Log::info('自动化打包, 发送登录请求邮件完毕');
    }

    /**
     * @param bool $flag
     * @param string $message
     */
    private function sendAuditMail($flag = false, $message = '')
    {
        $tag    = isset($this->data['tag']) ? 'tag:' . $this->data['tag'] : 'tag:';
        $branch = isset($this->data['branch']) ? 'branch:' . $this->data['branch'] : 'branch:';
        // 发布成功邮件通知相关人员
        $title   = '【Success】';
        $empty   = "&emsp; &emsp; &emsp;";
        $content = "
                <br>
                <br>
                {$empty}Dear all，<br>
                <br>
                <div>{$empty} <font color=\"##008000\">小程序提审成功</font></div>
                {$empty}版本号：" . $this->data['wx_version'] . "<br>
                {$empty}版本描述：" . $this->data['wx_desc'] . "<br>
                {$empty}" . $tag . "<br>
                {$empty}" . $branch . "<br>
                {$empty}提审时间：" . date('Y-m-d H:i:s') . "<br>
                <br>
                {$empty}请知悉，谢谢<br>
            ";
        if ($flag) {
            $title   = '【Error】';
            $content = "
                <br>
                <br>
                {$empty}Dear all，<br>
                <br>
                <div>{$empty}<font color=\"#ff0000\">小程序提审失败</font></div>
                {$empty}版本号：" . $this->data['wx_version'] . "<br>
                {$empty}版本描述：" . $this->data['wx_desc'] . "<br>
                {$empty}" . $tag . "<br>
                {$empty}" . $branch . "<br>
                {$empty}shop版本状态更新失败：" . $message . "<br>
                {$empty}提审时间：" . date('Y-m-d H:i:s') . "<br>
                <br>
                {$empty}请知悉，谢谢<br>
            ";
        }
        $copyTo = [
            "maxiongwei@vchangyi.com",
            "chengyujia@vchangyi.com",
            "handong@vchangyi.com",
            "wangliping@vchangyi.com",
            "gongguoxin@vchangyi.com"
        ];
        Mail::raw($content, function ($message) use ($copyTo, $title) {
            $message->getHeaders()->addTextHeader('Content-Type', 'text/html');
            $message->subject($title . '商城小程序自动提审通知（' . date('Y-m-d') . ')');
            $message->to("houyingcai@vchangyi.com");
            $message->cc($copyTo);
        });
        Log::info('自动化打包, 发送提审通知邮件完毕');
    }
}
