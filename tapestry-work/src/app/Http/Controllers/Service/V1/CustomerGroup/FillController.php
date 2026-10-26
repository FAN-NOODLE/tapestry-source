<?php


namespace App\Http\Controllers\Service\V1\CustomerGroup;


use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Cmq\CmqQueue;

class FillController extends BaseController
{

    public function index()
    {
         $queue = (new CmqQueue(config('cmq_setting.customer_group_relation')))->receiveMessage();
        if($queue) {
            //(new CmqQueue(config('cmq_setting.customer_group_relation')))->deleteMessage($queue->receiptHandle);
            $queueBody = json_decode($queue->msgBody, true);
            return $queueBody;
        }

        return $queue;

        /*$params = [
            'action' => 'custom_group_relation',
            'source' => 'crs',
            'query'  => [
                'enterprise_id' => 1,
                'relation'      => [
                    [
                        'user_id'         => 3,
                        'ec_info'         => [
                            [
                                'ec_userid' => '54654645645645',
                                'domain'    => 'd45345465trrt'
                            ]
                        ],
                        'customer_groups' => [1, 2, 3, 4],
                    ],
                    [
                        'user_id'         => 4,
                        'ec_info'         => [
                            [
                                'ec_userid' => '546546456456452',
                                'domain'    => 'd45345465trrt'
                            ]
                        ],
                        'customer_groups' => [1, 2, 3, 4],
                    ]
                ]
            ],
        ];

        (new CmqQueue(config('cmq_setting.customer_group_relation')))->sendMessage(json_encode($params));
        $queue = (new CmqQueue(config('cmq_setting.customer_group_relation')))->receiveMessage();
        if($queue) {
            //(new CmqQueue(config('cmq_setting.customer_group_relation')))->deleteMessage($queue->receiptHandle);
            $queueBody = json_decode($queue->msgBody, true);
            return $queueBody;
        }

        return $queue;*/

    }
}