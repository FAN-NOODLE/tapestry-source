<?php
/*
 * @Author: zhaojie
 * @Date: 2020-03-20 19:36:30
 * @LastEditTime: 2020-03-21 02:28:24
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \gw-shop\src\app\Library\ImageTools.php
 */
namespace App\Library;

use Illuminate\Support\Facades\Log;

class ImageTools
{

    /**
     * 从图片文件创建`Image资源
     * @param string $file 图片文件，支持url
     * @return bool|resource    成功返回图片image资源，失败返回false
     */
    function createImageFromFile($file)
    {
        if (preg_match('/http(s)?:\/\//', $file)) {
            $fileSuffix = $this->getNetworkImgType($file);
        } else {
            $fileSuffix = pathinfo($file, PATHINFO_EXTENSION);
        }

        if (!$fileSuffix) return false;

        switch ($fileSuffix) {
            case 'jpg':
            case 'jpeg':
                $theImage = @imagecreatefromjpeg($file);
                break;
            case 'png':
                $theImage = @imagecreatefrompng($file);
                if (empty($theImage)) {
                    $theImage = @imagecreatefromjpeg($file);
                }
                break;
            case 'gif':
                $theImage = @imagecreatefromgif($file);
                break;
            default:
                $theImage = @imagecreatefromstring(file_get_contents($file));
                break;
        }

        return $theImage;
    }


    /**
     * 获取网络图片类型
     * @param $url string 网络图片url,支持不带后缀名url
     * @return bool
     */
    public function getNetworkImgType($url)
    {
        $ch = curl_init(); //初始化curl
        curl_setopt($ch, CURLOPT_URL, $url); //设置需要获取的URL
        curl_setopt($ch, CURLOPT_NOBODY, 1);
        curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 3);//设置超时
        curl_setopt($ch, CURLOPT_TIMEOUT, 3);
        curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false); //支持https
        curl_exec($ch);//执行curl会话
        $http_code = curl_getinfo($ch);//获取curl连接资源句柄信息
        curl_close($ch);//关闭资源连接

        if ($http_code['http_code'] == 200) {
            $theImgType = explode('/', $http_code['content_type']);

            if ($theImgType[0] == 'image') {
                return $theImgType[1];
            } else {
                return false;
            }
        } else {
            return false;
        }
    }

    /**
     * 获取网络图片资源
     *
     * @param $url
     *
     * @return array
     */
    public function getImageSource($url)
    {
        $ch = curl_init(); //初始化curl
        curl_setopt($ch, CURLOPT_URL, $url); //设置需要获取的URL
        curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 10);//设置超时
        curl_setopt($ch, CURLOPT_TIMEOUT, 10);
        curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false); //支持https
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        $content   = curl_exec($ch);//执行curl会话
        $http_code = curl_getinfo($ch);//获取curl连接资源句柄信息
        curl_close($ch);//关闭资源连接
        $image = imagecreatefromstring($content);
        return [$image, $http_code];
    }

    /**
     * 获取网络图片
     * @param string $url 网络图片url,支持不带后缀名url
     * @return bool|string
     */
    public function getImage($url)
    {
        list($content, $http_code) = $this->getImageSource($url);

        if ($http_code['http_code'] == 200) {
            $theImgType = explode('/', $http_code['content_type']);
            if ($theImgType[0] == 'image') {
                $fileSuffix = $theImgType[1];
                $path = storage_path('download/passwordcoupon/');
                if(!is_dir($path)) {
                    mkdir($path, 0777, true);
                }
                $save_to    = storage_path('download/passwordcoupon/temp-' . uniqid());
                switch ($fileSuffix) {
                    case 'jpeg':
                        $save_to .= '.jpeg';
                        imagejpeg($content, $save_to);
                        break;
                    case 'jpg':
                        $save_to .= '.jpg';
                        imagejpeg($content, $save_to);
                        break;
                    case 'png':
                        $save_to .= '.png';
                        imagepng($content, $save_to);
                        break;
                    case 'gif':
                        $save_to .= '.gif';
                        imagegif($content, $save_to);
                        break;
                    default:
                        $save_to .= '.' . $fileSuffix;
                        Log::info('网络图片类型未识别:' . $save_to);
                        break;
                }
//                $downloaded_file = fopen($save_to, 'w');
//                fwrite($downloaded_file, $content);
//                fclose($downloaded_file);
                return $save_to;
            } else {
                return false;
            }
        } else {
            return false;
        }
    }

    /**
     * 方图片处理成圆图
     * @param $imgpath
     * @return string
     */
    public function toCircular($imgpath){
        $ext     = pathinfo($imgpath);
        if (empty($ext)) {
            return '';
        }
        $src_img = null;
        switch ($ext['extension']) {
            case 'png':
                $src_img = imagecreatefrompng($imgpath);
                break;
            default:
                $src_img = imagecreatefromjpeg($imgpath);
                break;
        }
        $img = $this->toCircularSource($src_img);
        // 如果不是png的图片，修改成png后缀
        if (strtolower($ext['extension']) != 'png') {
            unlink($imgpath);
            $basename = str_replace($ext['extension'], 'png', $ext['basename']);
            $imgpath = $ext['dirname'] . '/' . $basename;
        }
        imagepng($img, $imgpath);
        return $imgpath;
    }

    /**
     * 方图片处理成圆图
     * @param $imgpath
     * @return string
     */
    public function toCircularBySource($imageSource){

        $img = $this->toCircularSource($imageSource);
        return $img;
    }

    public function toCircularSource($image)
    {
        $w   = imagesx($image);
        $h   = imagesy($image);
        $w   = min($w, $h);
        $h   = $w;
        $img = imagecreatetruecolor($w, $h);
        imagesavealpha($img, true);
        //拾取一个完全透明的颜色,最后一个参数127为全透明
        $bg = imagecolorallocatealpha($img, 255, 255, 255, 127);
        imagefill($img, 0, 0, $bg);
        imagecolortransparent($img, $bg);
        $r   = $w / 2; //圆半径
//        $y_x = $r; //圆心X坐标
//        $y_y = $r; //圆心Y坐标
        for ($x = 0; $x < $w; $x++) {
            for ($y = 0; $y < $h; $y++) {
                $rgbColor = imagecolorat($image, $x, $y);
                if (((($x - $r) * ($x - $r) + ($y - $r) * ($y - $r)) < ($r * $r))) {
                    imagesetpixel($img, $x, $y, $rgbColor);
                }
            }
        }
        return $img;
    }
}