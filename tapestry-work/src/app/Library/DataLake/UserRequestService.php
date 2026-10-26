<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/9/1 0001
 * Time: 15:44
 */

namespace App\Library\DataLake;

use App\Library\DataLake\BaBind\BaBindQuery;
use App\Library\DataLake\Exception\QueryDataLakeBusinessError;
use App\Library\DataLake\Member\MemberRegister;
use App\Library\DataLake\Member\QueryBaList;
use App\Library\DataLake\Member\QueryMemberGrade;
use App\Library\DataLake\Member\QueryMemberGradeCalculator;
use App\Library\DataLake\Member\QueryMemberInfo;
use App\Library\DataLake\Member\QueryMemberOrderList;
use App\Library\DataLake\Member\QueryMemberPoints;
use App\Library\DataLake\Member\QueryMemberPointsDetailList;
use App\Library\DataLake\Member\QueryMemberPointsExpire;
use App\Library\DataLake\Member\QueryOrderDetail;
use App\Library\DataLake\Member\QueryStoreList;
use App\Library\DataLake\Member\QueryTierList;
use Throwable;


class UserRequestService implements UserRequestInterface
{
    /**
     * @var Request
     */
    private $request;

    public function __construct($enterpriseId = '')
    {
        $this->request = new Request($enterpriseId);
    }

    public function setEnterpriseId($enterpriseId)
    {
        $this->request->setEnterpriseId($enterpriseId);
    }

    /**
     * @param string $mobile
     * @return array
     * @throws Throwable
     */
    public function queryMemberInfo(string $queryType, string $value): array
    {
        $request = new QueryMemberInfo($queryType, $value);
        if (($response = $this->request->requestApi($request)) && $response->isOk()) {
            return $response->getData();
        }

        throw new QueryDataLakeBusinessError($response, '查询会员信息错误');
    }


    /**
     * @param string $mobile
     * @return array
     * @throws Throwable
     */
    public function queryGradeInfo(string $queryType, string $value): array
    {
        $request = new QueryMemberGrade($queryType, $value);
        if (($response = $this->request->requestApi($request)) && $response->isOk()) {
            return $response->getData();
        }
        throw new QueryDataLakeBusinessError($response, '查询会员等级错误');

    }

    /**
     * @param string $mobile
     * @return array
     * @throws Throwable
     */
    public function queryMemberGradeCalculator(string $queryType, string $value, string $grade, string $upgrade): array
    {
        $request = new QueryMemberGradeCalculator($queryType, $value, $grade, $upgrade);
        if (($response = $this->request->requestApi($request)) && $response->isOk()) {
            return $response->getData();
        }
        throw new QueryDataLakeBusinessError($response, '会员等级计算器错误');

    }

    /**
     * @param string $mobile
     * @return array
     * @throws Throwable
     */
    public function queryPointsExpire(
        string $queryType, string $value,
        string $pointTypeGroup, string $Startdate, string $Enddate
    )
    {
        $request = new QueryMemberPointsExpire($queryType, $value, $pointTypeGroup, $Startdate, $Enddate);
        if (($response = $this->request->requestApi($request)) && $response->isOk()) {
            return $response->getData();
        }
        throw new QueryDataLakeBusinessError($response, '会员等级计算器错误');

    }

    /**
     * @param string $mobile
     * @return array
     * @throws Throwable
     */
    public function queryPointsList(
        $queryType, $value, $pointTypeGroup,
        $Startdate, $Enddate, $pageNo, $limit
    ): array
    {
        $request = new QueryMemberPointsDetailList(
            $queryType, $value, $pointTypeGroup, $Startdate, $Enddate, $pageNo, $limit
        );
        if (($response = $this->request->requestApi($request)) && $response->isOk()) {
            return $response->getData();
        }
        throw new QueryDataLakeBusinessError($response, '会员等级计算器错误');

    }

    /**
     * @param string $mobile
     * @return array
     * @throws Throwable
     */
    public function getCurrentPoint(string $queryType, string $value): array
    {
        $request = new QueryMemberPoints($queryType, $value);
        if (($response = $this->request->requestApi($request)) && $response->isOk()) {
            return $response->getData();
        }
        throw new QueryDataLakeBusinessError($response, '查询会员积分错误');

    }


    public function getOrderList(string $queryType, string $value, int $page, int $size, string $tradeNo): array
    {
        $request = new QueryMemberOrderList($queryType, $value, $page, $size, $tradeNo);
        if (($response = $this->request->requestApi($request)) && $response->isOk()) {
            return $response->getData();
        }
        throw new QueryDataLakeBusinessError($response, '查询会员订单列表错误');
    }

    /**
     * @param string $mobile
     * @param string $orderNo
     * @return array
     * @throws Throwable
     */
    public function getOrderDetail(string $queryType, string $value, string $orderNo): array
    {
        $request = new QueryOrderDetail($queryType, $value, $orderNo);
        if (($response = $this->request->requestApi($request)) && $response->isOk()) {
            return $response->getData();
        }
        throw new QueryDataLakeBusinessError($response, '查询会员订单详情错误');
    }

    /**
     * @return array
     * @throws Throwable
     */
    public function getTierList(): array
    {
        $request = new QueryTierList();
        if (($response = $this->request->requestApi($request)) && $response->isOk()) {
            return $response->getData();
        }
        throw new QueryDataLakeBusinessError($response, '查询等级定义错误');
    }

    /**
     * @param string $storeNo
     * @param int $page
     * @param int $size
     * @return array
     * @throws Throwable
     */
    public function queryStoreList(string $storeNo, int $page, int $size): array
    {
        $request = new QueryStoreList($storeNo, $page, $size);
        if (($response = $this->request->requestApi($request)) && $response->isOk()) {
            return $response->getData();
        }
        throw new QueryDataLakeBusinessError($response, '查询门店错误');
    }

    /**
     * @param string $baCode
     * @param int $page
     * @param int $size
     * @return array
     * @throws Throwable
     */
    public function queryBAList(string $baCode, int $page, int $size): array
    {
        $request = new QueryBaList($baCode, $page, $size);
        if (($response = $this->request->requestApi($request)) && $response->isOk()) {
            return $response->getData();
        }
        throw new QueryDataLakeBusinessError($response, '查询ba错误');
    }

    /**
     * @param array $user
     * @return array
     * @throws Throwable
     */
    public function memberRegister(array $user): array
    {
        $request = new MemberRegister($user);
        if (($response = $this->request->requestApi($request)) && $response->isOk()) {
            return $response->getData();
        }
        throw new QueryDataLakeBusinessError($response, '注册错误');
    }

    // BA绑定关系

    /**
     * @throws Throwable
     */
    public function baBind(string $cusUnionID, string $cusUserID, string $baUserID, string $bindStatus, string $bindType, string $opDate, string $baCode, string $memberCode): array
    {
        $request = new BaBind\BaBind($cusUnionID, $cusUserID, $baUserID, $bindStatus, $bindType, $opDate, $baCode, $memberCode);
        if (($response = $this->request->requestApi($request)) && $response->isOk()) {
            return $response->getData() ?? [];
        }
        throw new QueryDataLakeBusinessError($response, '企业微信BA绑定解绑异常');
    }

    /**
     * @throws Throwable
     */
    public function baBindQuery(string $queryType, string $value, string $bindType): array
    {
        $request = new BaBindQuery($queryType, $value, $bindType);
        if (($response = $this->request->requestApi($request)) && $response->isOk()) {
            return $response->getData() ?? [];
        }
        throw new QueryDataLakeBusinessError($response, '企业微信BA绑定查询异常');
    }
}