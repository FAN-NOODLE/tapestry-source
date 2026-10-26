package helper

import (
	"context"
	"encoding/json"
	"fmt"
	ffmpegGo "github.com/u2takey/ffmpeg-go"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Ffmpeg struct {
	filePath         string         // 文件地址或路径
	ffmpegKwArgs     map[string]any // 压缩参数
	compressFilePath string         // 压缩后的文件路径
	fileExt          string
}

func NewFfmpeg(filePath string) *Ffmpeg {
	return &Ffmpeg{
		filePath:     filePath,
		ffmpegKwArgs: ffmpegGo.KwArgs{"threads": 1},
	}
}

func (f *Ffmpeg) BuildKwArgs(args map[string]any) *Ffmpeg {
	if len(args) > 0 {
		for key, val := range args {
			f.ffmpegKwArgs[key] = val
		}
	}

	return f
}

// Compress 执行压缩
func (f *Ffmpeg) Compress() (string, int, error) {
	// 获取文件名和文件后缀
	fileName := filepath.Base(f.filePath)

	md5Name := fileName + Int64ToString(time.Now().UnixNano())
	// 压缩文件地址
	f.compressFilePath = fmt.Sprintf("storage/upload/compress-%s.%s", Md5V(md5Name), f.fileExt)
	// 执行压缩
	err := ffmpegGo.Input(f.filePath).Output(f.compressFilePath, f.ffmpegKwArgs).OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		return "", 0, err
	}

	infoProbe, err := ffmpegGo.Probe(f.compressFilePath)
	if err != nil {
		return "", 0, err
	}

	var probeInfo sourceRequest.FfmpegProbe
	err = json.Unmarshal([]byte(infoProbe), &probeInfo)
	if err != nil {
		return "", 0, err
	}

	// 获取压缩文件大小
	compressFilesize := StringToInt(probeInfo.Format.Size)

	return f.compressFilePath, compressFilesize, nil
}

// EditVideoCreateTime 修改文件创建地址
func (f *Ffmpeg) EditVideoCreateTime() (string, error) {
	// 获取文件名和文件后缀
	fileName := filepath.Base(f.filePath)

	md5Name := fileName + Int64ToString(time.Now().UnixNano())
	// 压缩文件地址
	f.compressFilePath = fmt.Sprintf("storage/upload/compress-%s.%s", Md5V(md5Name), f.fileExt)

	f.BuildKwArgs(map[string]any{
		"metadata": "creation_time=" + time.Now().In(time.UTC).Format("2006-01-02T15:04:05Z"),
		"c":        "copy",
	})
	// 输出新视频
	err := ffmpegGo.Input(f.filePath).Output(f.compressFilePath, f.ffmpegKwArgs).OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		return "", err
	}

	return f.compressFilePath, nil
}

func (f *Ffmpeg) SetFileExt(ext string) *Ffmpeg {
	f.fileExt = ext
	return f
}

func (f *Ffmpeg) VideoBuildKwArgs() *Ffmpeg {
	f.fileExt = "mp4"

	f.ffmpegKwArgs["c:v"] = "libx264"
	infoProbe, _ := ffmpegGo.Probe(f.filePath)
	var probeInfo sourceRequest.FfmpegProbe
	_ = json.Unmarshal([]byte(infoProbe), &probeInfo)
	bitRate := 2000
	div := 1
	for _, value := range probeInfo.Streams {
		if value.CodecType == "video" {
			bitRateInt, _ := strconv.Atoi(value.BitRate)
			duration, _ := strconv.Atoi(value.Duration)
			if duration > 60 {
				div = 2
			}

			rateInt := bitRateInt / 1000 / div
			bitRate = rateInt
			// 码率大于2000的，最大为2000
			if rateInt > 2000 {
				bitRate = 2000
			}
			// 码率小于500的，最小为500
			if rateInt < 500 {
				bitRate = 500
			}

		}
	}
	bitRateStr := strconv.FormatInt(int64(bitRate), 10)
	bitRateStr += "k"
	f.ffmpegKwArgs["b:v"] = bitRateStr
	return f
}

func (f *Ffmpeg) ImageBuildKwArgs(args map[string]any) *Ffmpeg {
	f.fileExt = "jpg"
	f.ffmpegKwArgs["q"] = "100"
	// 判断是否存在压缩参数
	if len(args) > 0 {
		for key, val := range args {
			f.ffmpegKwArgs[key] = val
		}
	}

	return f
}

// UploadCos 上传文件到cos
func (f *Ffmpeg) UploadCos() (string, error) {
	// 删除本地文件
	defer func() { _ = os.Remove(f.compressFilePath) }()

	return UploadFileToCos(context.Background(), f.compressFilePath, false)
}

func (f *Ffmpeg) MakeCover() (string, string, error) {
	// 获取文件名和文件后缀
	fileName := filepath.Base(f.filePath)

	md5Name := fileName + Int64ToString(time.Now().UnixNano())
	// 压缩文件地址
	f.compressFilePath = fmt.Sprintf("storage/upload/compress-%s.%s", Md5V(md5Name), f.fileExt)
	imageUrl := fmt.Sprintf("storage/upload/compress-%s.%s", Md5V(md5Name), "jpg")
	// 执行压缩
	stream := ffmpegGo.Input(f.filePath)

	err := stream.Output(imageUrl, map[string]interface{}{"ss": 0,
		"vframes": 1}).OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		return "", "", err
	}
	err = stream.Output(f.compressFilePath, f.ffmpegKwArgs).OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		return "", imageUrl, err
	}

	infoProbe, err := ffmpegGo.Probe(f.compressFilePath)
	if err != nil {
		return "", "", err
	}

	var probeInfo sourceRequest.FfmpegProbe
	err = json.Unmarshal([]byte(infoProbe), &probeInfo)
	if err != nil {
		return "", "", err
	}

	return f.compressFilePath, imageUrl, nil
}
