<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/4/9 0009
 * Time: 10:17
 */

namespace App\Library\WechatNotice\Format;


class MarkDownFormat implements FormatInterface
{
    public function format(string $content): array
    {
        $data = [
            'msgtype'  => 'markdown',
            'markdown' => [
                'content'        => $content,
                'mentioned_list' => ["@all"]
            ]
        ];
        return $data;
    }
}