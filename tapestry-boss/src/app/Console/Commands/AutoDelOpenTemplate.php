<?php
/**
 * Created by Phpstorm.
 * User : zoulongbo
 * Date : 2019-09-27
 * Tome : 15:06
 */

namespace App\Console\Commands;

use App\Models\EnterpriseMiniProgramCodeLibrary;
use App\Models\EnterpriseOpenPlatform;
use App\Services\EnterpriseLogService;
use App\Services\OpenPlatform\MiniProgramService;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Mail;

class AutoDelOpenTemplate extends Command
{
    /**
     * @var string
     */
    const REFUSED = EnterpriseMiniProgramCodeLibrary::REFUSED;
    const WITHDRAW = EnterpriseMiniProgramCodeLibrary::WITHDRAW;
    const PENDING = EnterpriseMiniProgramCodeLibrary::PENDING;
    const AUDITED = EnterpriseMiniProgramCodeLibrary::AUDITED;
    const ONLINE = EnterpriseMiniProgramCodeLibrary::ONLINE;

    /**
     * @var string
     */
    protected $programService;

    /**
     * @var string
     */
    protected $signature = 'auto:del_open_template';

    /**
     * @var string
     */
    protected $description = '定时删除开放平台小程序模版';

    public function handle(MiniProgramService $programService)
    {
        EnterpriseLogService::info('开始执行自动删除开发平台模板');

        $openPlatformConfig = EnterpriseOpenPlatform::query()->get()->toArray();
        foreach ($openPlatformConfig as $value) {
            // 获取开放平台模版列表
            try {
                changeOpenPlatformConfig(0, $programService->openPlatform, $value['id']);
                $templateList = $programService->getTemplates();
            } catch (\Exception $e) {
                Log::info('自动删除开放平台模版 模版列表获取失败: ' . $e->getMessage());
                EnterpriseLogService::error('自动删除开放平台模版 模版列表获取失败', [
                    'message' => $e->getMessage(),
                    'file'    => $e->getFile(),
                    'line'    => $e->getLine()
                ]);
                return false;
            }

            EnterpriseLogService::info('自动删除开放平台模版，获取的微信模板列表', $templateList);
            if (empty($templateList)) {
                Log::info('自动删除开放平台模版 模版列表为空 data:' . json_encode($templateList));
                return false;
            }

            $templateIds = array_column($templateList, 'template_id');
            if (empty($templateIds)) {
                Log::info('自动删除开放平台模版 模版id集合为空 data:' . json_encode($templateList));
                return false;
            }

            $templateInfo = EnterpriseMiniProgramCodeLibrary::query()
                ->whereIn('template_id', $templateIds)
                ->whereIn('audit_status', [self::REFUSED, self::WITHDRAW, self::AUDITED])
                ->get(['account_id', 'template_id', 'audit_status', 'online_status', 'version', 'description'])
                ->toArray();
            if (empty($templateInfo)) {
                Log::info('自动删除开放平台模版 模版审核记录为空');
                return false;
            }

            // 即将删除或发通知的模版
            Log::info('即将自动删除或邮件通知的开放平台模版 DB模版数据 data:' . json_encode($templateInfo));
            EnterpriseLogService::info('即将自动删除或邮件通知的开放平台模版', $templateInfo);

            $this->programService = $programService;

            // 删除审核失败、撤审、已上线的模版
            $this->clearTemplate($templateInfo);
            // 删除企业作废旧重复模版
            $this->delOldTemplate($templateInfo);

            Log::info('自动删除开放平台模版 完成');
        }


        return true;
    }

    /**
     * 删除审核失败、撤审、已上线的模版
     * @param $templateInfo
     */
    private function clearTemplate(&$templateInfo)
    {
        // 审核通过待通知发布的邮件模版描述
        $waitReleaseTpl = "";

        foreach ($templateInfo as $key => $template) {
            try {

                // 删除审核失败、撤审、已上线的
                if (self::ONLINE == $template['online_status']  || in_array($template['audit_status'], [self::REFUSED, self::WITHDRAW])) {

                    $this->programService->deleteTemplate($template['template_id']);
                    unset($templateInfo[$key]);

                    Log::info('自动删除开放平台模版 DB模版删除成功：'.$template['version'] . '| template_id:'.$template['template_id']);
                    EnterpriseLogService::debug('自动删除开放平台模版-删除success', [
                        'version'     => $template['version'],
                        'template_id' => $template['template_id']
                    ]);
                }

                // 发邮件逻辑，暂时去除
                /*elseif (self::AUDITED == $template['audit_status'] && self::PENDING == $template['online_status']) {

                    // 待发布的发邮件提醒(需求要变更，暂停邮件服务)
                    $waitReleaseTpl .= "       【 模版ID：" . $template['template_id'] . "  -  " . $template['description']
                        . " 】 \n";
                }*/

            } catch (\Exception $e) {
                Log::info('自动删除开放平台模版 DB模版删除失败:'. $template['version']. '| template_id:'.$template['template_id']. '| no.'. $key. ' :' . $e->getMessage());
                EnterpriseLogService::error('自动删除开放平台模版-删除error', [
                    'version'     => $template['version'],
                    'template_id' => $template['template_id'],
                    'message'     => $e->getMessage(),
                    'file'        => $e->getFile(),
                    'line'        => $e->getLine()
                ]);
            }
        }

        // 发送提醒邮件
        /*if (!empty($waitReleaseTpl)) {
            @$this->sendTemplateMail($waitReleaseTpl);
        }*/
    }


    /**
     * 删除企业废旧模版
     * @param $templateInfo
     */
    private function delOldTemplate($templateInfo)
    {
        // 先取出所有企业，计算每个企业现占有的模版数量
        $arr = [];
        foreach ($templateInfo as $v) {
            $arr[$v['account_id']][] = $v['template_id'];
        }

        // 如果模版数量大于1，删除旧的模版
        foreach ($arr as $key => $value) {
            if (1 < count($value)) {
                $this->programService->deleteTemplate(min($value));
                Log::info('自动删除开放平台模版 企业废旧模版删除成功：| template_id:'.min($value));
                EnterpriseLogService::debug('自动删除开放平台模版-删除废旧模版success', ['template_id' => min($value)]);
            }
        }
    }

    /**
     * 发送提醒邮件
     * @param $waitReleaseTpl
     */
    private function sendTemplateMail($waitReleaseTpl)
    {
        $content = "

        Dear all，

        系统监控到以下小程序模版已提审通过，请联系相关品牌方负责人及时更新发版： \n"

            . $waitReleaseTpl .  " \n ".


            "      请知悉，谢谢！
        ";

        Mail::raw($content, function ($message) {
            $message->subject('小程序发版提醒（' . date('Y-m-d').')');
            $message->to("xuyang@vchangyi.com");
            $message->cc(["liupengwei@vchangyi.com", "weimengjie@vchangyi.com", "maxiongwei@vchangyi.com"]);
        });
        Log::info('小程序发版提醒, 发送邮件完成');
    }
}
