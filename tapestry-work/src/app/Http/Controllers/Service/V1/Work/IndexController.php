<?php

namespace App\Http\Controllers\Service\V1\Work;

use App\Http\Controllers\Controller;
use GuzzleHttp\Client;


class IndexController extends Controller
{

    public function index()
    {
        $params      = $this->request->post();
        $accessToken = $params['access_token'];
        $uri         = "https://qyapi.weixin.qq.com/cgi-bin/externalcontact/get_user_behavior_data?access_token=" . $accessToken;
        $partyId     = $params['party_id'];
        if (is_string($params['party_id'])) {
            $partyId = json_decode($params['party_id'], true);
        }
        $userid = $params['user_id'];
        if (is_string($params['user_id'])) {
            $userid = json_decode($params['user_id'], true);
        }
        $startTime          = $params['start_time'];
        $endTime            = $params['end_time'];
        $requestBody        = ["userid" => $userid, 'partyid' => $partyId, "start_time" => $startTime, "end_time" => $endTime];
        $client             = new Client();
        $options['headers'] = ["Content-Type" => "application/json"];
        $options['body']    = json_encode($requestBody);
        $response           = $client->post($uri, $options);

//        $client      = new Client();
//        $response    = $client->post($uri, null, $requestBody)->send();
        $req      = json_decode($response->getBody(), true);
        $behavior = $params['behavior'];
        if (is_string($params['behavior'])) {
            $behavior = json_decode($params['behavior'], true);
        }
        $new_apply_cnt         = $behavior['new_apply_cnt'] ?? 0;
        $new_contact_cnt       = $behavior['new_contact_cnt'] ?? 0;
        $chat_cnt              = $behavior['chat_cnt'] ?? 0;
        $message_cnt           = $behavior['message_cnt'] ?? 0;
        $reply_percentage      = $behavior['reply_percentage'] ?? 0.00;
        $avg_reply_time        = $behavior['avg_reply_time'] ?? 0;
        $negative_feedback_cnt = $behavior['negative_feedback_cnt'] ?? 0;
        $behaviorReq           = $req['behavior_data'];
        foreach ($behaviorReq as $item) {
            $new_apply_cnt         += $item['new_apply_cnt'];
            $new_contact_cnt       += $item['new_contact_cnt'];
            $chat_cnt              += $item['chat_cnt'];
            $message_cnt           += $item['message_cnt'];
            $reply_percentage      += $item['reply_percentage'];
            $avg_reply_time        += $item['avg_reply_time'];
            $negative_feedback_cnt += $item['negative_feedback_cnt'];
        }
        $result = [
            'new_apply_cnt'         => $new_apply_cnt,
            'new_contact_cnt'       => $new_contact_cnt,
            'chat_cnt'              => $chat_cnt,
            'message_cnt'           => $message_cnt,
            'reply_percentage'      => $reply_percentage,
            'avg_reply_time'        => $avg_reply_time,
            'negative_feedback_cnt' => $negative_feedback_cnt,
        ];
        return successJson($result);
    }
}