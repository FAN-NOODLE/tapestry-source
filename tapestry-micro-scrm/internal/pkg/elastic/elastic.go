package elastic

import (
	"encoding/json"
	"micro-scrm/internal/config"

	"github.com/elastic/go-elasticsearch/v8"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
)

var (
	ShareEsClient *Client
)

type Client struct {
	cli *elasticsearch.Client
}

type ClientOpt func(esc *Client)

func NewElasticClient(config config.Elastic, opts ...ClientOpt) error {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: config.EsAddress,
		Username:  config.EsUser,
		Password:  config.EsPassword,
	})
	if err != nil {
		return err
	}

	esc := &Client{cli: es}
	for _, opt := range opts {
		opt(esc)
	}

	res, err := esc.GetEsClient().Info()
	if err != nil {
		log.Error(ctx.New()).Msgf("NewElasticClient get es info Error response: %s", err)
		return err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.IsError() {
		log.Error(ctx.New()).Msgf("NewElasticClient res Error: %s", res.String())
		return err
	}

	var r map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Error(ctx.New()).Msgf("NewElasticClient parsing res body err: %s", err)
		return err
	}
	log.Info(ctx.New()).Msgf("es Client: %s", elasticsearch.Version)
	log.Info(ctx.New()).Msgf("es Server: %s", r["version"].(map[string]interface{})["number"])
	ShareEsClient = esc
	return nil
}

func (es *Client) GetEsClient() *elasticsearch.Client {
	return es.cli
}
