<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/25 0025
 * Time: 17:08
 */

namespace App\Library\DataLake;


interface UserRequestInterface
{
    public function queryMemberInfo(string $queryType, string $value): array;

    public function queryGradeInfo(string $queryType, string $value): array;

    public function getCurrentPoint(string $queryType, string $value): array;

    public function getOrderList(string $queryType, string $value, int $page, int $size, string $tradeNo): array;

    public function getOrderDetail(string $queryType, string $value, string $orderNo): array;

    public function getTierList(): array;

    public function queryStoreList(string $storeNo, int $page, int $size): array;

    public function queryBAList(string $baCode, int $page, int $size): array;

    public function memberRegister(array $user): array;

    // 企业微信BA绑定关系
    public function baBind(string $cusUnionID, string $cusUserID, string $baUserID, string $bindStatus, string $bindType, string $opDate, string $baCode, string $memberCode): array;

    public function baBindQuery(string $queryType, string $value, string $bindType): array;
}