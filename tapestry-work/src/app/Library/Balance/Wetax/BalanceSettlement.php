<?php


namespace App\Library\Balance\Wetax;


use GuzzleHttp\Client;
use GuzzleHttp\Exception\RequestException;
use Illuminate\Support\Facades\Log;
use function MongoDB\BSON\toJSON;

/**
 * 高灯科技-微纳税平台资金结算
 * Class BalanceSettlement
 * @package App\Library\Balance\Wetax
 */
class BalanceSettlement
{
    // 支付方式，1=银行卡 2=微信 3=支付宝 4 第三方微信
    const PAYMENT_WAY_BANK_CARD = 1;
    const PAYMENT_WAY_WECHAT = 2;
    const PAYMENT_WAY_ALIPAY = 3;
    const PAYMENT_WAY_THIRD_PARTY_WECHAT = 4;

    /**
     * @var SettlementConfig
     */
    private $config;

    /**
     * @var HttpHelper
     */
    private $httpHelper;

    public function __construct()
    {
        $this->config = SettlementConfig::getInstance();
        $this->httpHelper = new HttpHelper($this->config);
    }

    /**
     * 查询结算单状态
     * {
     *     "account_num": "ZD28f0f27", // 账单号
     *     "order_random_code": "3ff73660-12fd-4605-a851-2afa728f0f27", //客户订单号
     *     "settlement_code": "JS2016DE52B76C03", //结算单号
     *     "status": 860, //结算单状态
     *     "verification_code ":2, // 校验状态(1 代表校验中;2 代表校验结束)
     *     "service_amount": 0.100000, // 服务费
     *     "settlement_code": 1.000000, // 绩效费
     *     "total_money": 1.100000, // 合计金额
     *     "verification_error_code": "[]", // 错误码
     *     "is_verification": true, //账单号
     *     "fail_reason": "", //打款失败原因
     *     "hangup_flag": false, // 是否挂单
     *     "surplus_money": 0.0, // 本月剩余免增值税金额(挂单时有值，否则 为 0)
     *     "hangup_msg": null, // 挂单原因
     *     "hangup_begintime": null, // 挂单开始时间
     *     "project_code": null, // 项目编码
     *     "pay_time": null, // 支付时间(扣除商户余额的时间)
     *     "handle_time": null // 处理时间(打款成功或失败的时间)
     * }
     *
     * 说明：
     * 1. 结算单号是第三放生产的订单号，客户订单号是畅移生成的订单号，结算单号与客户订单号一一对应。
     * 2. 账单号和结算单号是一对多的关系，一个账单下可能会有多个结算单。
     * 3. 绩效费是用户收到的金额，服务费是第三方向畅移收取的费用，合计金额是服务费+绩效费，创建结算单的时候数据中传输的金额是用户收到的金额，也就是绩效费
     *
     * @param string $projectCode 项目编码
     * @param string $orderRandomCode 提现单号
     * @param string $settlementCode 结算单号
     * @return array
     */
    public function getBalance(string $projectCode, string $orderRandomCode, string $settlementCode = ''): array
    {
        $uri = '/api/balance/getbalance';
        $data = [
            'order_random_code' => $orderRandomCode,
            'settlement_code' => $settlementCode
        ];
        return $this->httpHelper->request($uri, $data, $projectCode);
    }

    /**
     * 结算单退款，退款结果会进行回调通知
     *
     * @param string $enterpriseHash 企业hash
     * @param string $projectCode 项目编码
     * @param string $orderRandomCode 客户订单编码
     * @param string $settlementCode 结算单编码
     * @return mixed
     */
    public function refundBalance(string $enterpriseHash, string $projectCode, string $orderRandomCode, string $settlementCode)
    {
        $uri = '/api/balance/refundbalance';
        $callbackUrl = $this->config->getSettlementRefundCallbackUrl() . $enterpriseHash;
        $data = [
            'settlement_code' => [$settlementCode],
            'order_random_code' => [$orderRandomCode]
        ];
        return $this->httpHelper->request($uri, $data, $projectCode, $callbackUrl);
    }

    /**
     * 发起结算单
     *
     * @param string $enterpriseHash 企业hash
     * @param string $projectCode 项目编码
     * @param string $orderRandomCode 客户订单号，每个结算单唯一
     * @param string $name 姓名
     * @param int $certificateType 证件类型，暂支持1=身份证，3=港澳居民来往内地通行证，4=港澳居民居住证，5=台湾居民来往大陆通行证，6=台湾居民居住证
     * @param string $certificateNum 证件号
     * @param string $phoneNum 手机号
     * @param float $settleAmount 金额，单位元，小数后两位
     * @param int $paymentWay 收款方式 1=银行卡 2=微信 3=支付宝 4 第三方微信，需客户APPID与高灯微信支付账户绑定，详情咨询高灯运营人员
     * @param string $bankName 银行名称，收款方式为银行卡时必填
     * @param string $bankCardNum 银行卡号，收款方式为银行卡时必填，只支持借记卡
     * @param string $paymentAccount 支付宝收款账户，收款方式为支付宝时，payment_account与alipay_accountid 必填其一
     * @param string $alipayAccountId 支付宝唯一Id，收款方式为支付宝时，payment_account与alipay_accountid 必填其一
     * @param string $wxOpenId 当payment_way=4 时， wx_openId 和 wx_appid 必填
     * @param string $wxAppid 当payment_way=4 时，wx_openId 和 wx_appid 必填
     * @return mixed
     */
    public function settleBalance(string $enterpriseHash,
                                  string $projectCode,
                                  string $orderRandomCode,
                                  string $name,
                                  int $certificateType,
                                  string $certificateNum,
                                  string $phoneNum,
                                  float $settleAmount,
                                  int $paymentWay,
                                  string $bankName = '',
                                  string $bankCardNum = '',
                                  string $paymentAccount = '',
                                  string $alipayAccountId = '',
                                  string $wxOpenId = '',
                                  string $wxAppid = ''
    )
    {
        $uri = '/api/balance/CreateForBatch';
        $callbackUrl = $this->config->getSettlementCallbackUrl() . $enterpriseHash;
        $data = [
            [
                'order_random_code' => $orderRandomCode,
                'name' => $name,
                'certificate_type' => $certificateType,
                'certificate_num' => $certificateNum,
                'phone_num' => $phoneNum,
                'settle_amount' => $settleAmount,
                'payment_way' => $paymentWay,
                'bank_name' => $bankName,
                'bankcard_num' => $bankCardNum,
                'payment_account' => $paymentAccount,
                'alipay_accountid' => $alipayAccountId,
                'wx_openId' => $wxOpenId,
                'wx_appid' => $wxAppid
            ]
        ];
        return $this->httpHelper->request($uri, $data, $projectCode, $callbackUrl);
    }

    /**
     * 批量结算
     *
     * @param string $enterpriseHash 企业hash
     * @param string $projectCode
     * @param array $settleData
     * @return mixed
     */
    public function batchSettleBalance(string $enterpriseHash, string $projectCode, array $settleData)
    {
        $uri = '/api/balance/CreateForBatch';
        $callbackUrl = $this->config->getSettlementCallbackUrl() . $enterpriseHash;
        return $this->httpHelper->request($uri, $settleData, $projectCode, $callbackUrl);
    }

    /**
     * 授权签约，返回结果为true表示签约成功，失败会抛出异常
     *
     * @param string $projectCode 项目编码
     * @param string $name 用户姓名
     * @param int $certificateType 证件类型(默认为 1)
     * @param string $certificateNum 证件号
     * @param string $phoneNum 手机号
     * @return bool
     */
    public function signAgreement(string $projectCode, string $name, int $certificateType, string $certificateNum,
                                  string $phoneNum): bool
    {
        $uri = '/api/balance/signagreement';
        $data = [
            'name' => $name,
            'certificate_type' => $certificateType,
            'certificate_num' => $certificateNum,
            'phone_num' => $phoneNum,
            'merchant_id' => $this->config->getAppKey()
        ];
        $this->httpHelper->request($uri, $data, $projectCode);
        return true;
    }

    /**
     * 获取签约结果
     * 响应结果：
     * [
     *     'is_signed' => true, // 签约结果，true 已签约，false 未签约
     *     'agreement_pic' => '' // 签约生成 PDF 格式文件链接(10 分钟有效期)
     * ]
     *
     * @param string $projectCode 项目编码
     * @param string $name 用户姓名
     * @param string $certificateNum 身份证号
     * @return array
     */
    public function queryAgreement(string $projectCode, string $name, string $certificateNum): array
    {
        $uri = '/api/balance/QueryAgreement';
        $data = [
            'name' => $name,
            'certificate_num' => $certificateNum
        ];
        return $this->httpHelper->request($uri, $data, $projectCode);
    }

    /**
     * 获取项目列表
     * 响应结果：
     * [
     *     {
     *         'project_code': "XM19BK0A8658AEA0", // 项目编号
     *         'project_name': "Test 项目", // 项目名称
     *         'create_time': "2019-11-20 20:46:06", // 创建时间
     *         'allow_merchant_balance': false // 是否允许使用商户余额
     *     }
     * ]
     *
     * @param $projectName
     * @return array
     */
    public function getProjectList(string $projectName): array
    {
        $uri = '/api/balance/getprojectlist';
        $data = [
            'project_name' => $projectName
        ];
        return $this->httpHelper->request($uri, $data);
    }

    /**
     * 根据项目编码获取项目信息
     * @param string $projectCode
     * @return array
     */
    public function getProjectByCode(string $projectCode): array
    {
        $uri = '/api/balance/getprojectbycode';
        $data = [
            'project_code' => $projectCode
        ];
        $project = $this->httpHelper->request($uri, $data);
        if (!is_array($project)) {
            $msg = is_string($project) && $project != null ? $project : "没有查询到该项目";
            throw new SettlementException(SettlementException::DATA_ERROR, $msg);
        }
        return $this->httpHelper->request($uri, $data);
    }
}