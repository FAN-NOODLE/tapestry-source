package callback

import (
	"errors"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/lock"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	workexternal "micro-scrm/internal/modules/work-external"
	"micro-scrm/internal/pkg/helper"
	pkgKafka "micro-scrm/internal/pkg/kafka"
	"micro-scrm/internal/pkg/sdk/work"
	"time"
)

// HandleCallbackMessage 处理外部联系人信息
func HandleCallbackMessage(c ctx.Context, data *work.MemberExternalCallback, handle WorkExternalHandleInterface) error {
	switch true {
	case data.Event == ExternalCallbackEventOfChangeExternalContact &&
		data.ChangeType == AddExternalContact:
		// 添加好友
		return handle.AddExternalContactHandle(c, data)

	case data.Event == ExternalCallbackEventOfChangeExternalContact &&
		data.ChangeType == DelExternalContact:
		// 员工删除好友
		return handle.DelExternalContactHandle(c, data)
	case data.Event == ExternalCallbackEventOfChangeExternalContact &&
		data.ChangeType == DelFollowUser:
		// 好友删除员工
		return handle.DelFollowUserHandle(c, data)
	case data.Event == ExternalCallbackEventOfChangeExternalChat:
		// 入群
	}

	return nil
}

type WorkExternalHandleInterface interface {
	// 添加企业客户事件
	AddExternalContactHandle(c ctx.Context, data *work.MemberExternalCallback) error
	// 编辑企业客户事件
	EditExternalContactHandle(c ctx.Context, data *work.MemberExternalCallback) error
	// 外部联系人免验证添加成员事件
	AddHalfExternalContactHandle(c ctx.Context, data *work.MemberExternalCallback) error
	// 删除企业客户事件
	DelExternalContactHandle(c ctx.Context, data *work.MemberExternalCallback) error
	// 删除跟进成员事件
	DelFollowUserHandle(c ctx.Context, data *work.MemberExternalCallback) error
	// 客户接替失败事件
	TransferFailHandle(c ctx.Context, data *work.MemberExternalCallback) error
}

// =============================== SyncToKafka 具体实现========================================
type SyncToKafka struct {
	service *Service
}

func NewSyncToKafka(s *Service) *SyncToKafka {
	return &SyncToKafka{
		s,
	}
}

func (s SyncToKafka) AddExternalContactHandle(c ctx.Context, data *work.MemberExternalCallback) error {
	// 入数据库
	record, err := s.service.AddRelationRecord(c, data)
	if err != nil {
		return err
	}

	// 同步到 kafka
	msg := workexternal.KafkaMsg{WorkRelationAddRecord: record}
	if err := s.ProducerMsg(msg); err != nil {
		return err
	}
	return nil
}

func (s SyncToKafka) EditExternalContactHandle(c ctx.Context, data *work.MemberExternalCallback) error {
	return errors.New("SyncToKafka not support EditExternalContactHandle ")
}

func (s SyncToKafka) AddHalfExternalContactHandle(c ctx.Context, data *work.MemberExternalCallback) error {
	return errors.New("SyncToKafka not support AddHalfExternalContactHandle ")
}

func (s SyncToKafka) DelExternalContactHandle(c ctx.Context, data *work.MemberExternalCallback) error {
	return s.delRelationHandle(c, data)
}

func (s SyncToKafka) DelFollowUserHandle(c ctx.Context, data *work.MemberExternalCallback) error {
	return s.delRelationHandle(c, data)
}

func (s SyncToKafka) TransferFailHandle(c ctx.Context, data *work.MemberExternalCallback) error {
	return errors.New("SyncToKafka not support TransferFailHandle ")
}

// ProducerMsg 推送消息到 kafka
func (s *SyncToKafka) ProducerMsg(msg workexternal.KafkaMsg) error {
	return pkgKafka.GetAppKafka().Producer(msg.ToJSONString(), "scrm-relation-update")
}

// delRelationHandle 删除关系操作
func (s *SyncToKafka) delRelationHandle(c ctx.Context, data *work.MemberExternalCallback) error {
	// 添加分布式锁
	key := helper.Md5V(data.ToUserName + data.ChangeType + data.UserID + data.ExternalUserID)
	mtx, err := lock.Mutex(
		c,
		lock.Key(key),
		lock.TTL(5*time.Second),
		lock.Retry(lock.NoRetry()),
	)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("产生分布式锁失败,错误信息：%s", err.Error())
		return err
	}
	err = mtx.Lock()
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("外部联系人数据加锁失败,错误信息：%s", err.Error())
		return err
	}

	defer mtx.Unlock()
	// 查找上次记录
	record, err := s.service.FindLastAddRelationRecord(c, data)
	// 未查询到数据则直接返回
	if record == nil || record.InCK() {
		return nil
	}

	// 同步到 kafka
	msg := workexternal.KafkaMsg{WorkRelationAddRecord: record, DelTime: data.CreateTime}
	if err := s.ProducerMsg(msg); err != nil {
		return err
	}

	s.service.FinishAddRelation(c, record)
	return nil
}
