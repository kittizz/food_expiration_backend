package woker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/alitto/pond"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"go.uber.org/fx"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/firebase"
)

type NotificationWorker struct {
	ctx       context.Context
	cancel    context.CancelFunc
	scheduler *gocron.Scheduler
	woker     *pond.WorkerPool

	firebase     *firebase.Firebase
	userUsecase  domain.UserUsecase
	itemUsecase  domain.ItemUsecase
	imageUsecase domain.ImageUsecase
}

func NewNotificationWorker(lc fx.Lifecycle, userUsecase domain.UserUsecase, itemUsecase domain.ItemUsecase, imageUsecase domain.ImageUsecase, firebase *firebase.Firebase) (*NotificationWorker, error) {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil, err
	}
	scheduler := gocron.NewScheduler(loc)
	woker := pond.New(10, 0)

	ctx, cancel := context.WithCancel(context.Background())
	noti := &NotificationWorker{
		ctx:          ctx,
		cancel:       cancel,
		scheduler:    scheduler,
		woker:        woker,
		userUsecase:  userUsecase,
		itemUsecase:  itemUsecase,
		firebase:     firebase,
		imageUsecase: imageUsecase,
	}

	noti.scheduler.Every(1).Minutes().StartImmediately().Do(noti.run)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			noti.scheduler.StartAsync()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			scheduler.Stop()
			woker.Stop()
			return nil
		},
	})
	return noti, nil
}

type NotiWoker struct {
	user  *domain.User
	items map[domain.NotificationStatus][]*domain.Item
}

type NOTIFICATION_TEMPLATE int

const (
	NOTIFICATION_TEMPLATE_SINGLE_ITEM_FOREWARN = iota
	NOTIFICATION_TEMPLATE_SINGLE_ITEM_EXPIRING
	NOTIFICATION_TEMPLATE_MULTIPLE_ITEM_FOREWARN
	NOTIFICATION_TEMPLATE_MULTIPLE_ITEM_EXPIRING
)

func (w *NotificationWorker) run() {
	log.Info().Msg("Running notification worker")
	defer log.Info().Msgf("Shutting down notification worker")
	notiAt, err := time.Parse(time.TimeOnly, "00:05:00")
	if err != nil {
		log.Err(err)
		return
	}

	users, err := w.userUsecase.ListNotifications(w.ctx, notiAt)
	if err != nil {
		log.Err(err)
		return
	}

	usersId := lo.Map(users, func(u *domain.User, _ int) int {
		return u.ID
	})

	itemsNoti, err := w.itemUsecase.ListForNotification(w.ctx, usersId)
	if err != nil {
		log.Err(err)
		return
	}

	notiWoker := make(map[int]NotiWoker, len(usersId))

	if len(itemsNoti) == 0 {
		return
	}
	for _, i := range itemsNoti {
		n, ok := notiWoker[i.UserID]
		if !ok {
			u, found := lo.Find(users, func(u *domain.User) bool {
				return u.ID == i.UserID
			})
			if !found {
				continue
			}
			n = NotiWoker{
				user:  u,
				items: make(map[domain.NotificationStatus][]*domain.Item, 0),
			}
			n.items[domain.NOTIFICATION_STATUS_FOREWARN] = make([]*domain.Item, 0)
			n.items[domain.NOTIFICATION_STATUS_EXPIRING] = make([]*domain.Item, 0)
			notiWoker[i.UserID] = n
		}

		status := dateStatus(i.ExpireDate, *i.ForewarnDay)
		switch i.NotificationStatus {
		case domain.NOTIFICATION_STATUS_PLANNED:
			n.items[status] = append(n.items[status], i)
		case domain.NOTIFICATION_STATUS_FOREWARN:
			if status == domain.NOTIFICATION_STATUS_EXPIRING {
				n.items[domain.NOTIFICATION_STATUS_EXPIRING] = append(n.items[domain.NOTIFICATION_STATUS_EXPIRING], i)
			}

		case domain.NOTIFICATION_STATUS_EXPIRING:
			if status == domain.NOTIFICATION_STATUS_EXPIRING {
				n.items[domain.NOTIFICATION_STATUS_EXPIRING] = append(n.items[domain.NOTIFICATION_STATUS_EXPIRING], i)
			}
		}

		notiWoker[i.UserID] = n
	}
	group, groupCtx := w.woker.GroupContext(w.ctx)

	var updateToStatus = make(map[domain.NotificationStatus][]int, 0)

	for _, n := range notiWoker {

		group.Submit(func() error {
			parameterData, _ := json.Marshal(map[string]any{
				"isLocation": false,
				"title":      "รายการทั้งหมด",
				"locationId": 0,
				"isScan":     false,
				"isSearch":   false,
			})
			msg := &messaging.Message{
				Notification: &messaging.Notification{
					Title: "null",
					Body:  "null",
				},
				Token: *n.user.FcmToken, // it's a single device token
				Android: &messaging.AndroidConfig{
					Priority: "high",
					Notification: &messaging.AndroidNotification{
						Visibility: messaging.VisibilityPublic,
					},
				},
				Data: map[string]string{
					"initialPageName": "ItemList",
					"parameterData":   string(parameterData),
				},
			}
			fcm, err := w.firebase.Messaging(groupCtx)
			if err != nil {
				return err
			}

			send := fcm.Send
			lenEXPIRING := len(n.items[domain.NOTIFICATION_STATUS_EXPIRING])
			lenFOREWARN := len(n.items[domain.NOTIFICATION_STATUS_FOREWARN])
			if lenEXPIRING > 1 {
				msg.Notification.Title = "แจ้งเตือนวันหมดอายุ"
				msg.Notification.Body = "มีรายการกำลังจะหมดอายุในวันนี้ ดูเพิ่มเติม"
				msg.Notification.ImageURL = "https://th-bkk-1.xvercloud.com/food-expiration/images/logo100x100.png"
				message_id, err := send(groupCtx, msg)
				if err != nil {
					return err
				}
				log.Info().Interface("message_id", message_id).Msg("send notification: EXPIRING > 1")
			} else if lenEXPIRING == 1 {
				item := n.items[domain.NOTIFICATION_STATUS_EXPIRING][0]
				img, err := w.imageUsecase.Get(groupCtx, item.ImageID)
				if err != nil {
					return err
				}
				msg.Notification.Title = "แจ้งเตือนวันหมดอายุ"
				msg.Notification.Body = fmt.Sprintf("รายการ %s กำลังจะหมดอายุในวันนี ดูเพิ่มเติม", *item.Name)
				msg.Notification.ImageURL = fmt.Sprintf("https://%s/%s%s", viper.GetString("BUCKET_ENDPOINT"), viper.GetString("BUCKET_NAME"), img.Path)
				message_id, err := send(groupCtx, msg)
				if err != nil {
					return err
				}
				log.Info().Interface("message_id", message_id).Msg("send notification: EXPIRING 1")
			}

			if lenFOREWARN > 1 {
				msg.Notification.Title = "แจ้งเตือนวันหมดอายุล่วงหน้า"
				msg.Notification.Body = "มีรายการใกล้จะหมดอายุ ดูเพิ่มเติม"
				msg.Notification.ImageURL = "https://th-bkk-1.xvercloud.com/food-expiration/images/logo100x100.png"
				message_id, err := send(groupCtx, msg)
				if err != nil {
					return err
				}
				log.Info().Interface("message_id", message_id).Msg("send notification: FOREWARN > 1")
			} else if lenFOREWARN == 1 {
				item := n.items[domain.NOTIFICATION_STATUS_FOREWARN][0]
				img, err := w.imageUsecase.Get(groupCtx, item.ImageID)
				if err != nil {
					return err
				}
				msg.Notification.Title = "แจ้งเตือนวันหมดอายุล่วงหน้า"
				msg.Notification.Body = fmt.Sprintf("รายการ %s ใกล้จะหมดอายุภายใร %d นี้  ดูเพิ่มเติม", *item.Name, *item.ForewarnDay)
				msg.Notification.ImageURL = fmt.Sprintf("https://%s/%s%s", viper.GetString("BUCKET_ENDPOINT"), viper.GetString("BUCKET_NAME"), img.Path)
				message_id, err := send(groupCtx, msg)
				if err != nil {
					return err
				}
				log.Info().Interface("message_id", message_id).Msg("send notification: FOREWARN 1")
			}

			return nil
		})

		for status, itemGroup := range n.items {
			for _, i := range itemGroup {
				switch status {
				case domain.NOTIFICATION_STATUS_FOREWARN:
					updateToStatus[domain.NOTIFICATION_STATUS_EXPIRING] = append(updateToStatus[domain.NOTIFICATION_STATUS_EXPIRING], i.ID)
				case domain.NOTIFICATION_STATUS_EXPIRING:
					updateToStatus[domain.NOTIFICATION_STATUS_DONE] = append(updateToStatus[domain.NOTIFICATION_STATUS_DONE], i.ID)
				}
			}
		}
	}

	err = w.itemUsecase.UpdateNotificationStatus(w.ctx, updateToStatus)
	if err != nil {
		log.Err(err)
		return
	}
	group.Wait()
}

func dateStatus(date time.Time, preDay int) domain.NotificationStatus {
	today := time.Now()
	difference := date.Sub(today).Hours() / 24

	if difference < 0 {
		return domain.NOTIFICATION_STATUS_EXPIRING
	} else if difference <= float64(preDay) {
		return domain.NOTIFICATION_STATUS_FOREWARN
	} else {
		return domain.NOTIFICATION_STATUS_PLANNED
	}

}
