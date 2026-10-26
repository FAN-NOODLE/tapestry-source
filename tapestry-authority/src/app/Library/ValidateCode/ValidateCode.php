<?php
namespace App\Library\ValidateCode;

//验证码类
class ValidateCode
{
    private $charset = 'abcdefghkmnprstuvwxyzABCDEFGHKMNPRSTUVWXYZ23456789';//随机字符库
    private $code;//验证码
    private $codelen = 4;//验证码长度
    private $width = 55;//宽度
    private $height = 20;//高度
    private $img;//图形资源句柄
//    private $font;//指定的字体
//    private $fontsize = 20;//指定字体大小
    private $fontcolor;//指定字体颜色

    //构造方法初始化
    public function __construct()
    {
//        $this->font =  base_path().'/app/Library/ValidateCode/Elephant.ttf';//字体包路径
        $this->createCode();
    }
    //生成随机码
    private function createCode()
    {
        $_len = strlen($this->charset) - 1;
        for ($i = 0;$i < $this->codelen;++$i) {
            $this->code .= $this->charset[mt_rand(0, $_len)];
        }
    }
    //生成背景
    private function createBg()
    {
        $this->img = imagecreatetruecolor($this->width, $this->height);
        $color = imagecolorallocate($this->img, mt_rand(157, 255), mt_rand(157, 255), mt_rand(157, 255));
        imagefilledrectangle($this->img, 0, $this->height, $this->width, 0, $color);
    }
    //生成文字
    private function createFont()
    {
        $_x = $this->width / $this->codelen;
        for ($i = 0;$i < $this->codelen;++$i) {
            $this->fontcolor = imagecolorallocate($this->img, mt_rand(0, 156), mt_rand(0, 156), mt_rand(0, 156));
//            imagettftext($this->img, $this->fontsize, mt_rand(-30, 30), $_x * $i + mt_rand(1, 5), $this->height / 1.4, $this->fontcolor, $this->font, $this->code[$i]);
            imagestring($this->img, 5,  $_x * $i + mt_rand(1, 5), $this->height / 3, $this->code[$i], $this->fontcolor);
        }
    }
    //生成线条、雪花
    private function createLine()
    {
        //线条
        for ($i = 0;$i < 6;++$i) {
            $color = imagecolorallocate($this->img, mt_rand(0, 156), mt_rand(0, 156), mt_rand(0, 156));
            imageline($this->img, mt_rand(0, $this->width), mt_rand(0, $this->height), mt_rand(0, $this->width), mt_rand(0, $this->height), $color);
        }
        //雪花
        for ($i = 0;$i < 100;++$i) {
            $color = imagecolorallocate($this->img, mt_rand(200, 255), mt_rand(200, 255), mt_rand(200, 255));
            imagestring($this->img, mt_rand(1, 5), mt_rand(0, $this->width), mt_rand(0, $this->height), '*', $color);
        }
    }
    //输出
    private function outPut()
    {
        //字太小，只好缩小背景图生成好后再放大图片
        $new = imagecreatetruecolor($this->width * 2, $this->height * 2);
        imagecopyresized($new, $this->img,0, 0,0, 0,$this->width * 2, $this->height * 2, $this->width, $this->height);
        ob_start ();
        imagepng($new);
        $imageData = ob_get_contents();
        ob_end_clean ();
        imagedestroy($this->img);
        return  "data:image/png;base64,". base64_encode($imageData);

    }
    //对外生成
    public function doimg()
    {
        $this->createBg();
        $this->createLine();
        $this->createFont();
        return $this->outPut();
    }
    //获取验证码
    public function getCode()
    {
        return strtolower($this->code);
    }
}
