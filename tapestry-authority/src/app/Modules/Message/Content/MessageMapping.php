<?php

namespace App\Modules\Message\Content;

class MessageMapping
{
    /*
     * $meta:
     *  {   "key1": {"from":"a", translate:"raw"},
     *      "key2":{"from":"b", translate:"raw"},
     *      "key3":{"from":"hello world!", translate:"const"}
     *      "key3":{"from":"hello world!", translate:"const"}
     *      "key3":{"from":"a", translate:"format", template: "hello %s"}
     * }
     *
     *     tanslate: raw/date/thing/text
     *
     * $obj:
     *  { "a":11, "b": 22 }
     *
     * return:
     *  ["key1"=>[
     *      "value"=> xxxx ],
     *  "key2"=>[
     *      "value"=> xxxx ],
     * ]
     * */
    public function translate($meta, $msg)
    {
        $res = [];
        foreach ($meta as $key => $item) {
            $translate = $item["translate"];
            $template = $item["template"] ?? '';
            $value = !empty($item["from"]) ? $this->xpathValue($msg, $item["from"]) : $template;

            $res[$key]['value'] = $this->{$translate}($value, $template);
        }
        return $res;

    }

    public function xpathValue($msg, $from)
    {
        return $msg[$from] ?? '';
    }

    public function phrase($value)
    {
        return $this->formatData($value, 5);
    }

    public function date($value, $template)
    {
        if ('now' == $template) {
            $value = time();
        }
        if (empty($value)) {
            $value = time();
        }
        return date('Y年m月d日 H:i', $value);
    }

    public function time($value, $template = '')
    {
        if ('now' == $template) {
            $value = time();
        }
        if (empty($value)) {
            $value = time();
        }
        return date('Y年m月d日 H:i', $value);
    }

    public function thing($value, $template = '')
    {
        if (is_int($value) && 10 == strlen($value)) {
            $value = $this->date($value, $template);
        }
        if ($template) {
            $value = sprintf($template, $value);
        }
        return $this->formatData($value, 20);
    }

    public function amount($value)
    {
        $amount = '¥' . ($value / 100);
        return $this->formatData($amount, 11);
    }

    public function characterString($value)
    {
        return $this->formatData($value, 32);
    }

    public function name($value)
    {
        return $this->formatData($value, 10);
    }

    public function number($value)
    {
        return $this->formatData($value, 32);
    }

    /**
     * 格式化事物类型字符长度
     *
     * @param string $thingValue
     * @param int $length 长度
     * @param string $replace
     * @return string
     */
    protected function formatData(string $thingValue, int $length, string $replace = '...')
    {
        // 长度大于20个字符
        if (mb_strlen($thingValue) > $length) {
            if (!empty($replace)) {
                $thingValue = mb_substr($thingValue, 0, $length - 3) . $replace;
            } else {
                $thingValue = mb_substr($thingValue, 0, $length);
            }
        }

        return $thingValue;
    }

    /**
     * 格式化time类型时间格式
     *
     * @param int $value 时间戳
     * @return String
     */
    protected function restTime(int $value)
    {
        $result = '00:00';
        $value = $value - time();
        if ($value > 0) {
            $hour = floor($value / 3600);
            $minute = floor(($value - 3600 * $hour) / 60);

            $result = $hour . ':' . $minute;
        }

        return $result;
    }
}
