package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	middleware "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus/mw"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"io/ioutil"
	internalHelper "micro-scrm/internal/pkg/helper"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// Recovery 崩溃恢复中间件
func Recovery() gin.HandlerFunc {
	return middleware.Register("Recovery", func(c *gin.Context) {
		requestBody, err := ioutil.ReadAll(c.Request.Body)
		c.Request.Body.Close()
		if err == nil {
			buf := bytes.NewBuffer(requestBody)
			c.Request.Body = ioutil.NopCloser(buf)
		}

		defer func() {
			if err := recover(); err != nil {
				stack := stack(2)

				notice := os.Getenv("qy_notice_key")
				if notice != "" {
					// 定义 文件名、行号、方法名
					fileName, line, functionName := "", 0, ""
					var pc uintptr
					var ok bool
					pc, fileName, line, ok = runtime.Caller(3)
					if ok {
						functionName = runtime.FuncForPC(pc).Name()
						functionName = filepath.Ext(functionName)
						functionName = strings.TrimPrefix(functionName, ".")
					}

					content := fmt.Sprintf("# %s项目异常 \n "+
						">请求时间：<font color='comment'>%s</font>\n"+
						">请求地址：<font color='comment'>%s</font>\n"+
						">请求参数：<font color='comment'>%s</font>\n"+
						">文件名：<font color='comment'>%s</font>\n"+
						">行号：<font color='comment'>%d</font>\n"+
						">函数名：<font color='comment'>%s</font>\n"+
						">错误信息：<font color='comment'>%v</font>\n"+
						">调用栈：<font color='comment'>%v</font>",
						os.Getenv("app"),
						internalHelper.GetFormatDateByParams(0),
						c.Request.Host+c.Request.RequestURI,
						string(requestBody),
						fileName,
						line,
						functionName,
						err,
						string(stack),
					)

					go func() {
						message := gin.H{
							"msgtype": "markdown",
							"markdown": gin.H{
								"content":        content,
								"mentioned_list": []string{"@all"}}}
						endPoint := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + notice
						httpClient := clients.NewHTTPClient()
						_, err = httpClient.Post(ctx.Wrap(c), endPoint, clients.JSONBody(message))
						if err != nil {
							fmt.Println("异常报警错误", err)
						}
					}()
				}

				log.Error(ctx.Wrap(c)).Str("stack", string(stack)).Msgf("[panic]: %v", err)
				err := errorx.NewErrorMsg(500, fmt.Sprintf("服务器发生错误 [panic]: %v", err))
				ginplus.ResError(c, err, 500)
			}
		}()
		requestID := c.GetHeader(ctx.RequestIDKey.String())
		requestID = ctx.Request(requestID).ID()
		c.Set(ctx.RequestIDKey.String(), requestID)
		c.Set(ctx.LocaleKey, c.DefaultQuery("locale", "zh"))
		c.Set(ctx.LangKey, ginplus.AcceptLanguage(c.GetHeader("Accept-Language")))
		c.Next()
	})
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
