<?php


namespace App\Services;


use App\Exceptions\Error;
use CrsSdk\MicroUser\User;
use Illuminate\Support\Facades\Log;

class UserRemoteService
{

    public static function getUserInfoByUnionId(int $enterpriseId, string $unionId)
    {
        if (!empty($enterpriseId) && !empty($unionId)) {
            try {
                return (new User())->userDetail([
                    'enterprise_id' => $enterpriseId,
                    'union_id' => $unionId
                ]);
            } catch (Error $e) {
                Log::error('用户信息查询失败', [
                    'code' => $e->getCode(),
                    'status_code' => $e->getStatusCode(),
                    'message' => $e->getMessage(),
                    'exception' => $e->getTraceAsString()
                ]);
            }
        }
        return [];
    }

    public static function getUserInfoById(int $enterpriseId, int $userId)
    {
        if (!empty($enterpriseId) && !empty($userId)) {
            try {
                return (new User())->userDetail([
                    'enterprise_id' => $enterpriseId,
                    'id' => $userId
                ]);
            } catch (Error $e) {
                Log::error('用户信息查询失败', [
                    'code' => $e->getCode(),
                    'status_code' => $e->getStatusCode(),
                    'message' => $e->getMessage(),
                    'exception' => $e->getTraceAsString()
                ]);
            }
        }
        return [];
    }

    /**
     * 获取用户列表
     */
    public static function listUserInfo(int $enterpriseId, array $userIdList)
    {
        if (!empty($enterpriseId) && !empty($userIdList)) {
            $params = [
                'enterprise_id' => $enterpriseId,
                'conditions' => [
                    ['id', 'in', $userIdList]
                ],
                'is_paginate' => 0
            ];
            try {
                return (new User())->userList($params)['list'];
            } catch (Error $e) {
                Log::error('用户信息列表查询失败', [
                    'code' => $e->getCode(),
                    'status_code' => $e->getStatusCode(),
                    'message' => $e->getMessage(),
                    'exception' => $e->getTraceAsString()
                ]);
            }
        }
        return [];
    }

    public static function listUserInfoByUnionId(int $enterpriseId, array $unionIdList)
    {
        if (!empty($enterpriseId) && !empty($unionIdList)) {
            $params = [
                'enterprise_id' => $enterpriseId,
                'conditions' => [
                    ['union_id', 'in', $unionIdList]
                ],
                'is_paginate' => 0
            ];
            try {
                return (new User())->userList($params)['list'];
            } catch (Error $e) {
                Log::error('用户信息列表查询失败', [
                    'code' => $e->getCode(),
                    'status_code' => $e->getStatusCode(),
                    'message' => $e->getMessage(),
                    'exception' => $e->getTraceAsString()
                ]);
            }
        }
        return [];
    }

    public static function listUserUnionId(int $enterpriseId, array $userIdList)
    {
        $userInfoList = self::listUserInfo($enterpriseId, $userIdList);
        $result = [];
        foreach ($userInfoList as $userInfo) {
            $result[$userInfo['id']] = $userInfo['union_id'] ?? '';
        }
        return $result;
    }
}