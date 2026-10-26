package appJob

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"micro-scrm/internal/pkg/xjob"
)

type cosPrivateField struct {
	Database string
	Table    string
	Field    string
}

type cosPrivateURL struct {
	RawURL string
	Key    string
}

type cosPrivateURLRewriter struct {
	oldHosts    map[string]struct{}
	publicHost  string
	privateHost string
	pattern     *regexp.Regexp
}

func newCosPrivateURLRewriter(oldOriginHost, oldCdnHost, privateHost string) cosPrivateURLRewriter {
	hosts := make([]string, 0, 2)
	oldHosts := make(map[string]struct{}, 2)
	for _, host := range []string{oldOriginHost, oldCdnHost} {
		if host == "" {
			continue
		}
		if _, exists := oldHosts[host]; exists {
			continue
		}
		oldHosts[host] = struct{}{}
		hosts = append(hosts, regexp.QuoteMeta(host))
	}

	return cosPrivateURLRewriter{
		oldHosts:    oldHosts,
		publicHost:  oldCdnHost,
		privateHost: privateHost,
		pattern:     regexp.MustCompile(`https?://(?:` + strings.Join(hosts, "|") + `)/[^"'\s<>,)\\]+`),
	}
}

func (r cosPrivateURLRewriter) rewritePublic(value string, copied map[string]struct{}) string {
	return r.pattern.ReplaceAllStringFunc(value, func(rawURL string) string {
		key, ok := r.objectKey(rawURL)
		if !ok {
			return rawURL
		}
		if _, exists := copied[key]; !exists {
			return rawURL
		}

		parsed, err := url.Parse(rawURL)
		if err != nil {
			return rawURL
		}
		parsed.Host = r.publicHost
		parsed.Path = "/" + cosPrivatePublicTargetKey(key)
		return parsed.String()
	})
}

func (r cosPrivateURLRewriter) extract(value string) []cosPrivateURL {
	urls := r.pattern.FindAllString(value, -1)
	result := make([]cosPrivateURL, 0, len(urls))
	seenKeys := make(map[string]struct{}, len(urls))
	for _, rawURL := range urls {
		key, ok := r.objectKey(rawURL)
		if !ok {
			continue
		}
		if _, exists := seenKeys[key]; exists {
			continue
		}
		seenKeys[key] = struct{}{}
		result = append(result, cosPrivateURL{RawURL: rawURL, Key: key})
	}
	return result
}

func (r cosPrivateURLRewriter) rewrite(value string, copied map[string]struct{}) string {
	return r.pattern.ReplaceAllStringFunc(value, func(rawURL string) string {
		key, ok := r.objectKey(rawURL)
		if !ok {
			return rawURL
		}
		if _, exists := copied[key]; !exists {
			return rawURL
		}

		urlStart := strings.Index(rawURL, "://") + len("://")
		parsed, err := url.Parse(rawURL)
		if err != nil {
			return rawURL
		}
		return rawURL[:urlStart] + r.privateHost + rawURL[urlStart+len(parsed.Host):]
	})
}

func (r cosPrivateURLRewriter) objectKey(rawURL string) (string, bool) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", false
	}
	if _, exists := r.oldHosts[parsed.Host]; !exists {
		return "", false
	}

	key := strings.TrimLeft(parsed.Path, "/")
	if key == "" || isCosPrivatePublicKey(key) {
		return "", false
	}
	return key, true
}

func isCosPrivatePublicKey(key string) bool {
	key = strings.TrimLeft(key, "/")
	for _, prefix := range cosPrivatePublicPrefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

func cosPrivatePublicTargetKey(key string) string {
	return "uploads/" + strings.TrimLeft(key, "/")
}

func cosPrivateShouldMoveImagePublic(field cosPrivateField, key string) bool {
	isSelectedField := field.Database == "micro_scrm" && field.Table == "members" && field.Field == "qr_poster_code"
	isSelectedField = isSelectedField || strings.HasPrefix(field.Database, "micro_scrm_enterprise_") && field.Table == "activity_qrcode_point" && field.Field == "code_url"
	if !isSelectedField {
		return false
	}
	extension := strings.ToLower(path.Ext(key))
	return extension == ".png" || extension == ".jpg"
}

func cosPrivateBaseFields() []cosPrivateField {
	return []cosPrivateField{
		{Database: "micro_scrm", Table: "chat_group_msg_task", Field: "content"},
		{Database: "micro_scrm", Table: "crm_sftp_downloaded_file", Field: "file_url"},
		{Database: "micro_scrm", Table: "dashboard_log", Field: "file_url"},
		{Database: "micro_scrm", Table: "download_resource", Field: "source_data"},
		{Database: "micro_scrm", Table: "drainage_welcomes", Field: "content"},
		{Database: "micro_scrm", Table: "excitation_integral_rule_description", Field: "content"},
		{Database: "micro_scrm", Table: "excitation_level_rule", Field: "level_icon"},
		{Database: "micro_scrm", Table: "excitation_medal_style", Field: "light_icon"},
		{Database: "micro_scrm", Table: "excitation_medal_style", Field: "unlit_icon"},
		{Database: "micro_scrm", Table: "marketing_strategy_channels", Field: "content"},
		{Database: "micro_scrm", Table: "members", Field: "qr_poster_code"},
		{Database: "micro_scrm", Table: "order_achievement_sftp_log", Field: "file_url"},
		{Database: "micro_scrm", Table: "reserve_order_data_acxiom_sftp_log", Field: "file_url"},
		{Database: "micro_scrm", Table: "template", Field: "template_picture_url"},
		//{Database: "micro_work", Table: "contact_qr_code", Field: "transform_url"},
		//{Database: "micro_work", Table: "store_worker", Field: "qr_code"},
		//{Database: "micro_work", Table: "marketing_plan_task_resource", Field: "poster_url"},
		//{Database: "micro_work", Table: "member_banner_relation", Field: "suncode_url"},
	}
}

func cosPrivateEnterpriseFields(database string) []cosPrivateField {
	return []cosPrivateField{
		{Database: database, Table: "activity_qrcode_main", Field: "banner_url"},
		{Database: database, Table: "activity_qrcode_point", Field: "code_url"},
		{Database: database, Table: "activity_qrcode_welcomes", Field: "content"},
		{Database: database, Table: "display_task", Field: "review_file"},
		{Database: database, Table: "display_task_allot_data", Field: "image_url"},
		{Database: database, Table: "display_task_allot_remark", Field: "avatar"},
		{Database: database, Table: "display_task_allot_remark", Field: "content"},
		{Database: database, Table: "drainage_contact_group_welcomes", Field: "content"},
		{Database: database, Table: "external_pay_order", Field: "tickets"},
		{Database: database, Table: "invite_record", Field: "background_url"},
		{Database: database, Table: "invite_template", Field: "background_url"},
		{Database: database, Table: "invite_template", Field: "share_url"},
		{Database: database, Table: "notification_crowd", Field: "file_url"},
		{Database: database, Table: "notification_info", Field: "attachment"},
		{Database: database, Table: "notification_info", Field: "page_url"},
		{Database: database, Table: "notification_info", Field: "pdf_url"},
		{Database: database, Table: "notification_info", Field: "zip_url"},
		{Database: database, Table: "notification_info_merge", Field: "cover"},
		{Database: database, Table: "online_warehouse_order", Field: "transfer_file"},
		{Database: database, Table: "online_warehouse_order", Field: "share_poster_url"},
		{Database: database, Table: "plan_task_conditions", Field: "content"},
		{Database: database, Table: "plan_task_import_file", Field: "file_url"},
		{Database: database, Table: "plan_task_import_file", Field: "download_url"},
		{Database: database, Table: "plan_task_template", Field: "template_content"},
		{Database: database, Table: "plan_tasks", Field: "download_url"},
		{Database: database, Table: "reserve_service", Field: "reserve_img"},
		{Database: database, Table: "store_visit_main", Field: "pdf_url"},
		{Database: database, Table: "store_visit_media", Field: "url"},
		{Database: database, Table: "store_visit_media", Field: "video_image_url"},
		{Database: database, Table: "user_group_file", Field: "file_url"},
	}
}

func cosPrivateMigrationFields(handler *xjob.JobHandler) ([]cosPrivateField, error) {
	fields := cosPrivateBaseFields()
	enterpriseIDs, err := GetEnterpriseIds(handler)
	if err != nil {
		return nil, err
	}
	for _, enterpriseID := range enterpriseIDs {
		fields = append(fields, cosPrivateEnterpriseFields(fmt.Sprintf("micro_scrm_enterprise_%d", enterpriseID))...)
	}
	return fields, nil
}
