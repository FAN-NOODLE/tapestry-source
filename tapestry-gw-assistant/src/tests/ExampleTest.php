<?php


use App\Models\Order;
use app\Services\Order\OrderMachine;
use App\Services\Order\OrderStateConfig;

class ExampleTest extends TestCase
{
    /**
     * 基础url
     */
    const APP_BASE_URL = 'http://crs.shop.me/app/v1';
    const APP_MACHINE_STATE_URL = '/new-order/index';
    const SERVICE_BASE_URL = 'http://crs.shop.me/service/v1';
    const SERVICE_MACHINE_STATE_URL = '/new-order/index';
    const MOCK_USER_ID = 'Mock-UserId';
    const MOCK_ENTERPRISE_ID = 'Mock-EnterpriseId';
    private $userId = 307;
    private $enterpriseId = 2;
    /**
     * 构造请求头
     *
     * @param $params array
     * @return array
     */
    public function buildHeader($params = [])
    {
        $params = array_merge($params, [self::MOCK_USER_ID => $this->userId,self::MOCK_ENTERPRISE_ID => $this->enterpriseId]);
        $header = $this->transformHeadersToServerVars($params);
        return $header;
    }

    public function common_create_order(){
        $data = [
            'payment_channel'   => 1,
            'user_coupon_id' => 0,
            'order_type' =>1,
            'good_data' =>[
                [
                    "spu_id"=> 406,
                    "sku_id"=> 549,
                    "quantity"=> 1,
                    "type"=> 1,
                    "car_id"=> 2576
                ]
            ],
            'mem_uid' =>"",
            'type' =>1,
            'delivery_name' =>"123",
            'delivery_mobile' =>"13022223333",
            'delivery_address' =>"qqweqweqwe",
            'zip_code' =>"100000",
            'logistics_price' =>0,
            'user_address_id' =>475,
            'province' =>"北京市",
            'city' =>"北京市",
            'area' =>"东城区",
            'check_shop_activity_id' =>0,
        ];

        $response = $this->call('POST', self::APP_BASE_URL . '/new-order/add', $data,[],[],$this->buildHeader());
        return $response;
    }

    public function getOrderId($response_content){
        $response_content = json_decode($response_content,true);
        $order_no = $response_content['data']['order_no'];
        $order_info =Order::query()->where('order_no',$order_no)->first();
        return $order_info->id;
    }

    /**
     * 将订单更新为需要的状态
     * @param $response_content
     * @param $status
     *
     * @return mixed
     */
    public function updateOrder($response_content,$status){
        $response_content = json_decode($response_content,true);
        $order_no = $response_content['data']['order_no'];
        $order_info =Order::query()->where('order_no',$order_no)->first();
        $order_info->status = $status;
        $order_info->save();
        return $order_info->id;
    }


    /**
     * 用户申请仅退款
     * @param $order_id
     *
     * @return mixed
     */
    public function common_apply_refund($order_id){
        $data = [
            'id'=>$order_id,
            'event'=>OrderStateConfig::EVENT_CUSTOMER_APPLY_REFUND,
            'reason'=> '不想要了',
        ];
        $response = $this->call('POST', ExampleTest::APP_BASE_URL . ExampleTest::APP_MACHINE_STATE_URL, $data,[],[],$this->buildHeader());
    }

    /**
     * 用户申请退货退款
     * @param $order_id
     *
     * @return mixed
     */
    public function common_apply_return($order_id){
        $data = [
            'id'=>$order_id,
            'event'=>OrderStateConfig::EVENT_CUSTOMER_APPLY_RETURN,
            'reason'=> '不想要了',
        ];
        $response = $this->call('POST', ExampleTest::APP_BASE_URL . ExampleTest::APP_MACHINE_STATE_URL, $data,[],[],$this->buildHeader());
    }

    /**
     * 用户再次申请退货退款
     * @param $order_id
     *
     * @return mixed
     */
    public function common_apply_return_again($order_id){
        $data = [
            'id'=>$order_id,
            'event'=>OrderStateConfig::EVENT_CUSTOMER_APPLY_RETURN_AGAIN,
            'reason'=> '不想要了',
        ];
        $response = $this->call('POST', ExampleTest::APP_BASE_URL . ExampleTest::APP_MACHINE_STATE_URL, $data,[],[],$this->buildHeader());
    }
}
