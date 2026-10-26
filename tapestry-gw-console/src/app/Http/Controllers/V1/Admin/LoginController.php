<?php

namespace App\Http\Controllers\V1\Admin;

use App\Exceptions\Error;
use CrsSdk\MicroAuthority\Account;
use Exception;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Support\Facades\Config;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;

class LoginController extends BaseController
{

    /**
     * 非必须登录
     *
     * @var bool
     */
    protected $requireLogin = false;

    /**
     * 企业ID
     */
    protected $enterpriseId = null;


    /**
     * 企业端登录接口
     *
     * @return JsonResponse|void
     * @throws Exception
     */
    public function index()
    {
        $account   = trim($this->request->post('account'));
        $password  = trim($this->request->post('password'));
        $loginType = trim($this->request->post('login_type', 1));

        if (empty($account)) {
            return error(451001);
        }

        if (empty($password)) {
            return error(451002);
        }

        // 调用微服务的登录接口， 返回企业列表
        try {
            $loginReturnData = (new Account())->login([
                'account'    => $account,
                'password'   => $password,
                'login_type' => $loginType,
                'client'     => self::CLIENT,
            ]);
        } catch (\Exception $e) {
            if ($e->getCode() == 700005) {
                if ($loginType == 1) {
                    notice(100008);
                } else {
                    notice(100012);
                }
            }

            if ($e->getCode() == 700006) {
                if ($loginType == 1) {
                    notice(100009);
                } else {
                    notice(100010);
                }
            }

            throw new Error($e->getCode(), $e->getMessage(), [], 422);
        }

        $enterpriseList  = $loginReturnData['enterprise_admin_list'];
        $enterpriseTotal = count($enterpriseList);

        $hash    = '';
        $data    = [];
        $version = 0;

        if ($enterpriseTotal == 1) {
            $enterpriseInfo = $enterpriseList[0];
            $version        = $enterpriseInfo['version'];
            $adminId        = $enterpriseInfo['admin_id'];

            if ($version == 1) {

                // 旧版切到旧库,防止角色里的权限不一致
                $thisDatabase = Config::get('database.connections.mysql.database');
                $oldDatabase = str_replace('gw_', '', $thisDatabase);

                Log::info('旧系统登录切换DB', [
                    '当前DB' => $thisDatabase,
                    '待切DB' => $oldDatabase
                ]);

                Config::set('database.connections.mysql.database', $oldDatabase);
                DB::reconnect();

                Log::info('旧系统登录切换DB完成', Config::get('database.connections.mysql'));

                $hash = $enterpriseInfo['is_split_flow'] ? $enterpriseInfo['hash'] : 'public';
                $data = $this->getAuthorityDataForV1($enterpriseInfo, $adminId);
            } else {
                $hash = $enterpriseInfo['is_split_flow'] ? $enterpriseInfo['hash'] : 'store';
            }
        }

        // 删除关闭应用标识缓存
        Redis::connection()->del(['enterprise:application:close:' . $this->enterpriseId]);

        return success([
            'enterprise_total' => $enterpriseTotal,
            'version'          => $version,
            'token'            => $loginReturnData['token'] ?? '',
            'hash'             => $hash,
            'enterprise_list'  => $enterpriseList,
            'data'             => (object)$data,
        ]);
    }
}
