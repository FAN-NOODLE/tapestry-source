package qyWork

type OpenUseridToUseridReq struct {
	OpenUserIDList []string `json:"open_userid_list"`
	SourceAgentID  int      `json:"source_agentid"`
}

type OpenUseridToUseridResp struct {
	UseridList            []UseridInfo `json:"userid_list"`
	InvalidOpenUseridList []string     `json:"invalid_open_userid_list"`
}

type UseridInfo struct {
	OpenUserid string `json:"open_userid"`
	Userid     string `json:"userid"`
}

type OpenExternalUseridToExternalUseridReq struct {
	ExternalUserid string `json:"external_userid"`
	SourceAgentID  int    `json:"source_agentid"`
}

type OpenExternalUseridToExternalUseridResp struct {
	ExternalUserid string `json:"external_userid"`
}
