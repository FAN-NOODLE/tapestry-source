package helper

import (
	"archive/zip"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"micro-scrm/internal/modules/marketingPromotion/storage"
	"micro-scrm/internal/modules/marketingPromotion/storage/cos"
	"micro-scrm/internal/pkg/cosroute"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// FileSize 获取文件大小
func FileSize(filePath string) int64 {
	if IsRemoteFile(filePath) {
		fileName := filepath.Base(filePath)
		localPath := fmt.Sprintf("uploads/%s/%s", time.Now().Format("2006/01/02"), fileName)
		err := DownloadFileToLocal(filePath, localPath)
		if err != nil {
			return 0
		}
		filePath = localPath
		defer func() {
			_ = os.Remove(localPath)
		}()
	}

	file, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return file.Size()
}

// CreateFile 创建文件，不存在路径自动创建
func CreateFile(filePath string) (*os.File, error) {
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// IsRemoteFile 判断是否为远程文件
func IsRemoteFile(path string) bool {
	protocols := []string{"http://", "https://"}

	for _, protocol := range protocols {
		if strings.HasPrefix(path, protocol) {
			return true
		}
	}

	return false
}

func ReadRemoteFile(c context.Context, fileURL string) ([]byte, error) {
	if isCOSRemoteFileURL(fileURL) {
		return readCOSRemoteFile(c, fileURL)
	}

	req, err := http.NewRequestWithContext(c, http.MethodGet, fileURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download remote file failed: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

var readCOSRemoteFile = func(c context.Context, fileURL string) ([]byte, error) {
	tmpFile, err := os.CreateTemp("", "micro-scrm-remote-*")
	if err != nil {
		return nil, err
	}
	tmpPath := tmpFile.Name()
	defer func() {
		_ = os.Remove(tmpPath)
	}()
	if err = tmpFile.Close(); err != nil {
		return nil, err
	}

	cs := storage.Init(cos.COS)
	if cs == nil {
		return nil, fmt.Errorf("cos storage is not initialized")
	}
	if err = cs.Download(c, fileURL, tmpPath); err != nil {
		return nil, err
	}

	return os.ReadFile(tmpPath)
}

func isCOSRemoteFileURL(fileURL string) bool {
	_, err := cosroute.ParseObjectURL(fileURL)
	return err == nil
}

// DownloadFileToLocal 下载文件到本地
func DownloadFileToLocal(url string, filePath string) error {
	data, err := ReadRemoteFile(ctx.New(), url)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("download remote file failed: empty content")
	}

	file, err := CreateFile(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	if _, err = file.Write(data); err != nil {
		return err
	}
	return nil
}

// UploadFileToCos 上传本地文件到cos
// localPath 本地文件绝对路径
func UploadFileToCos(c context.Context, localPath string, private bool) (string, error) {
	// 获取文件名和文件后缀
	fileName := filepath.Base(localPath)

	cs := storage.Init(cos.COS)
	uploadFile := fmt.Sprintf(RuleThree(private, "private/", "")+"uploads/%s/%s", time.Now().Format("2006/01/02"), fileName)
	path, err := cs.Upload(c, localPath, uploadFile)
	if err != nil {
		return "", err
	}

	return cs.GetURL(path), nil
}

// GetImageSize 获取图片尺寸
func GetImageSize(imagePath string) (int, int, error) {
	// 打开图片文件
	img, err := imaging.Open(imagePath, imaging.AutoOrientation(true))
	if err != nil {
		return 0, 0, err
	}
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	return width, height, nil
}

// CalculateNewDimensions 计算尺寸
func CalculateNewDimensions(originalWidth, originalHeight, maxWidth, maxHeight int) (int, int) {
	var newWidth, newHeight int
	aspectRatio := float64(originalWidth) / float64(originalHeight)
	if originalWidth > originalHeight {
		newWidth = maxWidth
		newHeight = int(float64(newWidth) / aspectRatio)
		if newHeight > maxHeight {
			newHeight = maxHeight
			newWidth = int(float64(newHeight) * aspectRatio)
		}
	} else {
		newHeight = maxHeight
		newWidth = int(float64(newHeight) * aspectRatio)
		if newWidth > maxWidth {
			newWidth = maxWidth
			newHeight = int(float64(newWidth) / aspectRatio)
		}
	}
	return newWidth, newHeight
}

func ResizedImage(fileLocalPath string, maxWidth, maxHeight int) error {
	// 打开图片文件
	file, err := os.Open(fileLocalPath)
	if err != nil {
		return err
	}
	// 解码图片
	img, _, err := image.Decode(file)
	// 获取图片尺寸
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	// 计算图片尺寸
	newWidth, newHeight := CalculateNewDimensions(width, height, maxWidth, maxHeight)

	resizedImage := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)
	// 创建新的图片文件
	outFile, err := os.Create(fileLocalPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
		_ = outFile.Close()
	}()

	// 将调整分辨率后的图片编码为JPEG格式并输出到创建的文件
	err = jpeg.Encode(outFile, resizedImage, &jpeg.Options{Quality: 100})
	if err != nil {
		return err
	}
	return nil
	////下载文件到本地
	//err := helper.DownloadFileToLocal(filePath, fileLocalPath)
	//if err != nil {
	//	return "", err
	//}
	////压缩文件
	//err = helper.ResizedImage(fileLocalPath, 1440, 1080)
	//if err != nil {
	//	return "", err
	//}
	//// 上传cos
	//cosFileAddr, err := helper.UploadFileToCos(c, fileLocalPath)
	//if err != nil {
	//	return "", err
	//}
}

// CreateZip creates a zip file containing the files from the given paths.
func CreateZip(zipFileName string, files map[string]string) error {
	// 创建 ZIP 文件
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return err // 返回错误
	}
	defer zipFile.Close() // 确保在函数退出时关闭 ZIP 文件

	// 创建 ZIP 写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close() // 确保在函数退出时关闭 ZIP 写入器

	// 遍历文件映射
	for localPath, zipPath := range files {
		// 打开本地文件
		fileToZip, err := os.Open(localPath)
		if err != nil {
			return err // 返回错误
		}
		defer fileToZip.Close() // 确保在函数退出时关闭文件

		// 获取 ZIP 中的目录路径
		dirPath := filepath.Dir(zipPath)
		if dirPath != "" {
			// 创建目录路径（如果目录不为空）
			_, err = zipWriter.Create(dirPath + "/")
			if err != nil {
				return err // 返回错误
			}
		}

		// 在 ZIP 中创建文件
		w, err := zipWriter.Create(zipPath)
		if err != nil {
			return err // 返回错误
		}

		// 将本地文件内容复制到 ZIP 文件中
		_, err = io.Copy(w, fileToZip)
		if err != nil {
			return err // 返回错误
		}
	}
	return nil // 成功
}

func UnZipToDir(zipFilePath string, unZipDir string) ([]string, error) {
	var filePathSlice []string
	// 打开 ZIP 文件
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return filePathSlice, err
	}
	defer r.Close()

	// 遍历 ZIP 文件中的每个文件
	for _, f := range r.File {
		// 生成文件解压后的路径
		filePath := filepath.Join(unZipDir, f.Name)

		// 如果文件是一个目录，创建目录
		if f.FileInfo().IsDir() {
			continue
		}

		_, err := os.Stat(unZipDir)
		if err != nil {
			if os.IsNotExist(err) {
				// 创建目标文件
				if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
					return filePathSlice, err
				}
			}
		}

		// 打开目标文件并解压内容
		outFile, err := os.Create(filePath)
		if err != nil {
			return filePathSlice, err
		}
		defer outFile.Close()

		// 打开 ZIP 文件中的具体文件
		rc, err := f.Open()
		if err != nil {
			return filePathSlice, err
		}
		defer rc.Close()

		// 将压缩文件的内容拷贝到目标文件
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return filePathSlice, err
		}
		filePathSlice = append(filePathSlice, filePath)
	}
	return filePathSlice, err

}

// DeleteFilesInDirectory 删除目录下文件（不包含目录中的目录）
func DeleteFilesInDirectory(dir string) error {
	// 获取目录中的所有文件
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("unable to read directory: %w", err)
	}

	// 遍历所有文件
	for _, file := range files {
		if file.Name() == ".gitignore" {
			continue
		}
		// 获取文件的完整路径
		filePath := filepath.Join(dir, file.Name())

		// 如果是文件，则删除
		if !file.IsDir() {
			err := os.Remove(filePath)
			if err != nil {
				return fmt.Errorf("unable to remove file %s: %w", filePath, err)
			}
			fmt.Println("Deleted file:", filePath)
		}
	}
	return nil
}

// DownloadImage 下载图片的简单实现
func DownloadImage(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP状态码: %d", resp.StatusCode)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
