<?php

namespace App\Traits;

trait ToolsTrait
{
    /**
     * 过滤字符串中表情
     *
     * @param $str string 昵称
     * @return string
     */
    private function filterEmoji($str)
    {
        $str = preg_replace_callback('/./u', function (array $match) {
            return strlen($match[0]) >= 4 ? '' : $match[0];
        }, $str);

        return $str;
    }

    /**
     * 数组转XML
     *
     * @param $arr array 数组
     * @return string
     */
    private function arrayToXml($arr)
    {
        $xml = $xml ?? '';
        foreach ($arr as $key => $val) {
            if (is_array($val)) {
                $xml .= '<' . $key . '>' . $this->arrayToXml($val) . '</' . $key . '>';
            } else {
                $xml .= '<' . $key . '>' . $val . "</" . $key . '>';
            }
        }

        return $xml;
    }

    /**
     * XML转数组
     *
     * @param $xml string XML
     * @return bool|array
     */
    private function xmlToArray($xml)
    {
        $parser = xml_parser_create();
        if (!xml_parse($parser, $xml, true)) {
            xml_parser_free($parser);

            return false;
        }
        libxml_disable_entity_loader(true);
        $xmlString = simplexml_load_string($xml, 'SimpleXMLElement', LIBXML_NOCDATA);
        $arr = json_decode(json_encode($xmlString), true);

        return $arr;
    }

    /**
     * 去除空格和特殊字符
     *
     * @param $str
     * @return string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    private function trimStr($str)
    {
        $search  = array(" ", "　", "\n", "\r", "\t", "\/");
        $replace = array("", "", "", "", "", "");

        return $this->filterEmoji(str_replace($search, $replace, $str));
    }
}
