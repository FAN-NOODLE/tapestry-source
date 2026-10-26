package storeVisitSvc

import (
	"encoding/json"
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"html/template"
	storeVisitDao "micro-scrm/internal/app/dao/store_visit"
	"micro-scrm/internal/app/http/request"
	storeVisitReq "micro-scrm/internal/app/http/request/store_visit"
	"micro-scrm/internal/pkg/helper"
	pkgKafka "micro-scrm/internal/pkg/kafka"
	"os"
	"time"
)

type CreatePDFQueueSvc struct {
}

func NewCreatePDFQueueSvc() *CreatePDFQueueSvc {
	return &CreatePDFQueueSvc{}
}

// CreatePDFProduce PDF生成 生产者
func (p *CreatePDFQueueSvc) CreatePDFProduce(c ctx.Context, params storeVisitReq.CreatePdfReq) {
	str, _ := json.Marshal(params)
	err := pkgKafka.GetAppKafka().Producer(string(str), "create_store_visit_pdf")
	if err != nil {
		log.Info(c).Msgf("生成pdf数据进入队列失败,错误原因：%+v", err)
	}
}

func (p *CreatePDFQueueSvc) CreatePDFFile(c ctx.Context, params storeVisitReq.CreatePdfReq) error {
	info, err := NewPdfInfoSvc().getPdfInfo(c, storeVisitReq.VisitInfoReq{
		BaseReq: storeVisitReq.BaseReq{
			EnterpriseBaseReq: request.EnterpriseBaseReq{
				EnterpriseId: params.EnterpriseId,
			},
			Language: params.Language,
		},
		ID: params.ID,
	})
	if err != nil {
		return err
	}
	if len(info.ShowRecordList) == 0 {
		return nil
	}
	pdfGenerator, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(c).Msgf("生成pdf文件失败,错误原因：%+v", err)
		return err
	}
	pdfGenerator.Dpi.Set(300) // 高DPI确保清晰度
	pdfGenerator.MarginBottom.Set(15)
	pdfGenerator.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfGenerator.PageSize.Set(wkhtmltopdf.PageSizeA4)
	cover := "./storage/store_visit/" + helper.RuleThree(params.Language == "en", "en", "zh") + "/cover.html"
	tmpl, err := template.ParseFiles(cover)
	if err != nil {
		return err
	}
	coverName := "cover1.html"
	coverFile, err := os.Create(coverName)
	if err != nil {
		return err
	}
	data := struct {
		Name            string
		StoreVisitTime  string
		ReportTime      string
		CreatorName     string
		CreatorDeptName string
	}{
		Name:            info.StoreName,
		StoreVisitTime:  info.Date + " " + helper.FormatStoreVisitTimeSpan(params.Language, info.TimeSpan),
		ReportTime:      helper.GetFormatDateByParams(info.UpdatedAt),
		CreatorName:     info.CreatorName,
		CreatorDeptName: info.CreatorDeptName,
	}
	err = tmpl.Execute(coverFile, data)
	if err != nil {
		return err
	}
	defer coverFile.Close()
	defer os.Remove(coverName)
	pdfGenerator.Cover.Input = coverName
	pdfGenerator.Cover.Zoom.Set(1)
	footer := "./storage/store_visit/" + helper.RuleThree(params.Language == "en", "en", "zh") + "/footer.html"
	footerTmpl, err := template.ParseFiles(footer)
	if err != nil {
		return err
	}
	footerName := "footer1.html"
	footerFile, err := os.Create(footerName)
	if err != nil {
		return err
	}
	footerData := struct {
		CreatorName     string
		CreatorDeptName string
	}{
		CreatorName:     info.CreatorName,
		CreatorDeptName: info.CreatorDeptName,
	}
	err = footerTmpl.Execute(footerFile, footerData)
	if err != nil {
		return err
	}
	defer footerFile.Close()
	defer os.Remove(footerName)
	contentUrl := "./storage/store_visit/" + helper.RuleThree(params.Language == "en", "en", "zh") + "/content.html"
	// 创建模板并注册函数
	contentTmpl := template.Must(template.New("content.html").Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s) // 将字符串转为安全HTML
		},
	}).ParseFiles(contentUrl))
	contentName := "content1.html"
	contentFile, err := os.Create(contentName)
	if err != nil {
		return err
	}
	err = contentTmpl.Execute(contentFile, info)
	if err != nil {
		return err
	}
	defer contentFile.Close()
	defer os.Remove(contentName)
	contentPage := wkhtmltopdf.NewPage(contentName) //内容pdf
	contentPage.EnableLocalFileAccess.Set(true)     //needed to include js
	contentPage.PageOffset.Set(0)
	header := "./storage/store_visit/" + helper.RuleThree(params.Language == "en", "en", "zh") + "/header.html"
	contentPage.PageOptions.HeaderHTML.Set(header) //页眉
	contentPage.PageOptions.HeaderSpacing.Set(10)
	contentPage.PageOptions.FooterHTML.Set(footerName) //页脚
	contentPage.PageOptions.FooterRight.Set("[page]")
	contentPage.PageOptions.FooterFontSize.Set(10)
	contentPage.PageOptions.FooterSpacing.Set(0)
	pdfGenerator.AddPage(contentPage)
	err = pdfGenerator.Create()
	if err != nil {
		log.Fatal(c).Msgf("生成pdf文件失败,错误原因：%+v", err)
		return err
	}
	pdfName := fmt.Sprintf("%d_%d.pdf", params.ID, time.Now().Unix())
	err = pdfGenerator.WriteFile(pdfName)
	if err != nil {
		log.Fatal(c).Msgf("生成pdf文件失败,错误原因：%+v", err)
		return err
	}
	defer os.Remove(pdfName)
	cosUrl, err := helper.UploadFileToCos(c, pdfName, true)
	if err != nil {
		return err
	}
	err = storeVisitDao.NewMainDao().UpdateData(c, params.EnterpriseId, map[string]interface{}{
		"id": params.ID,
	}, map[string]interface{}{
		"pdf_url": cosUrl,
	})
	return err
}
