<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/4/9 0009
 * Time: 10:17
 */

namespace App\Library\WechatNotice\Format;


class TextFormat implements FormatInterface
{
    public function format(string $content): array
    {
        $data = [
            'msgtype' => 'text',
            'text'    => [
                'content'        => $content,
                'mentioned_list' => ["@all"]
            ]
        ];
        return $data;
    }
}