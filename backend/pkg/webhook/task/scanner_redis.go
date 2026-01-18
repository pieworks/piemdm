package task

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"piemdm/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// ProviderSet is cron providers.
// var ProviderSet = wire.NewSet(NewScanner, service.NewWebhookService)

type WebhookService interface {
	Get(id uint) (*model.Webhook, error)
}

type EntityService interface {
	Get(c *gin.Context, tableCode string, id uint) (map[string]any, error)
}

type WebhookDeliveryService interface {
	Create(webhookDelivery *model.WebhookDelivery) error
}

type Scanner struct {
	rdb                    *redis.Client
	webhookService         WebhookService
	entityService          EntityService
	webhookDeliveryService WebhookDeliveryService
}

const (
	ScannerSize = 10
)

func NewScanner(rdb *redis.Client, webhookService WebhookService, entityService EntityService, webhookDeliveryService WebhookDeliveryService) *Scanner {
	return &Scanner{
		rdb:                    rdb,
		webhookService:         webhookService,
		entityService:          entityService,
		webhookDeliveryService: webhookDeliveryService,
	}
}

func (s *Scanner) Run() {
	fmt.Println("Scanner Run.")

	err := s.scannerRedis()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (s *Scanner) scannerRedis() error {
	fmt.Println("scannerRedis.")
	go func() {
		for {
			listName := "WebhookQueue"
			// LPOP：在队列最左侧取出一个值(移除并返回列表的第一个元素)
			ctx := context.Background()
			val, err := s.rdb.BRPop(ctx, 10*60*time.Second, listName).Result()
			if err != nil {
				fmt.Println("s.rdb.BRPop: ", err.Error())
			}
			fmt.Println("scannerRedis:", time.Now())

			if val != nil {
				req := model.WebhookReq{}
				// json.Unmarshal([]byte(val), &req)
				if json.Unmarshal([]byte(val[1]), &req); err != nil {
					log.Println("unmarshal: ", err)
					return
				}

				go s.DeliverOne(req)
			}
		}
	}()

	return nil
}

func (s *Scanner) DeliverOne(webhookReq model.WebhookReq) error {
	// 	// 根据entity_id 进行push

	// 	// 获取 Webhook config
	hook, err := s.webhookService.Get(webhookReq.HookID)
	if err != nil {
		return err
	}
	fmt.Println("hook, err: ", hook, err)

	// 获取 entity
	// RPUSH "WebhookQueue" "{\"HookID\":334502,\"ApprovalCode\":\"B2E74AA0-4C82-4E1B-B337-DB28E4CE621B\",\"TableCode\":\"entity_115\",\"EntityID\":2003000562}"
	c := &gin.Context{}
	entity, err := s.entityService.Get(c, webhookReq.TableCode, webhookReq.EntityID)
	if err != nil {
		return err
	}
	// log.Printf("relation: %+v\n", relation)

	go s.DeliverPost(webhookReq, hook, entity)
	return nil
}

// 请求方式
// 目前支持 GET 和 POST 两种请求方法
// Get：请求指定的页面信息，并返回响应主体，可用来检索或获取信息
// Post：向指定资源提交数据进行处理，数据被包含在请求文本中。用来创建新资源、修改现有资
//
// 请求头
// HTTP 请求中的 Header 参数，以键–值 (key–value) 的形式填入
//
// 数据类型
// plain
// application/json
// application/xxx-form-urlencoded
// application/xml
//
// 请求体: 请求的 body 信息，支持上述四种格式
// Basic Auth 用户名: 部分 HTTP 服务器存在 Basic Auth 用户名认证，可通过填充此参数进行认证
// Basic Auth  密码: 部分 HTTP 服务器存在 Basic Auth 密码认证，可通过填充此参数进行认证
//
// 返回值示例
// HTTP 响应体，将按照 JSON 规则对响应体进行解析并返回相应的内容，对应的 key值可以被流程后续节点消费
// 注：为确保返回值内容能够被引用，需保证输入的内容符合 JSON 格式
func (s *Scanner) DeliverPost(webhookReq model.WebhookReq, hook *model.Webhook, entity map[string]any) {
	now := time.Now()
	start := time.Now() // 获取当前时间
	uuid := uuid.New()
	deliveryCode := strings.ToUpper(uuid.String())
	// 转成Json 进行进行POST调用
	entityJson, _ := json.Marshal(entity)

	posturl := hook.Url
	req, err := http.NewRequest("POST", posturl, bytes.NewReader(entityJson))
	if err != nil {
		log.Fatalf("impossible to build request: %s", err)
	}
	req.Header.Set("Content-Type", hook.ContentType+"; charset=UTF-8")
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-GitHub-Delivery", deliveryCode)
	req.Header.Set("X-GitHub-Event", "Release")
	req.Header.Set("X-GitHub-Hook-ID", fmt.Sprintf("%d", hook.ID))
	// User-Agent: GitHub-Hookshot/044aadd
	// req.Header.Set("X-GitHub-Hook-Installation-Target-ID", "667721582")
	// req.Header.Set("X-GitHub-Hook-Installation-Target-Type", "repository")

	// create http client
	// do not forget to set timeout; otherwise, no timeout!
	client := http.Client{Timeout: 10 * time.Second}
	// send the request
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	// get request Header
	reqKeys := make([]string, 0)
	for key := range req.Header {
		reqKeys = append(reqKeys, key)
	}
	sort.Strings(reqKeys)
	reqHeader := ""
	for _, key := range reqKeys {
		reqHeader = reqHeader + fmt.Sprintf("%+v: %+v\n", key, req.Header[key][0])
	}

	// get request Body
	reqBodyCopy, err := req.GetBody()
	if err != nil {
		log.Println("return a new copy of Body error: ", err)
	}
	reqBodyByte, err := io.ReadAll(reqBodyCopy)
	if err != nil {
		log.Println("read request body error: ", err)
	}
	reqBody := string(reqBodyByte)

	// get response Header
	resHeader := ""
	resKeys := make([]string, 0)
	for key := range res.Header {
		resKeys = append(resKeys, key)
	}
	sort.Strings(resKeys)
	for _, key := range resKeys {
		resHeader = resHeader + fmt.Sprintf("%+v: %+v\n", key, res.Header[key][0])
	}

	// get response Body
	resBodyByte, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	// save to database
	delivery := model.WebhookDelivery{
		HookID:          hook.ID,
		DeliveryCode:    deliveryCode,
		Event:           "Release",
		EntityID:        uint(webhookReq.EntityID),
		RequestHeaders:  reqHeader,
		RequestPayload:  reqBody,
		ResponseStatus:  res.StatusCode,
		ResponseHeaders: resHeader,
		ResponseBody:    string(resBodyByte),
		DeliveredAt:     &now,
		CompletedAt:     &now,
	}
	if err := s.webhookDeliveryService.Create(&delivery); err != nil {
		// 保存失败的时候，重新把数据放入队列
		// TODO
		// 发送失败的时候，重新放入队列，保存失败，是否可以不用保存。
		// 或者采用先保存，然后发送，然后更新记录。
		log.Println("DeliverPost Save to Database err. ")

		// // 将 map[string]any 转换为 JSON 字符串
		// dataJson, err := json.Marshal(webhookReq)
		// if err != nil {
		// 	log.Println("DeliverPost Failed to marshal data to JSON:", err)
		// 	return
		// }

		// // 存储消息
		// ctx := context.Background()
		// res, err := s.rdb.LPush(ctx, "WebhookQueue", dataJson).Result()
		// if err != nil {
		// 	log.Printf("DeliverPost redis.LPush Error: %+v\n", err)
		// }
		// log.Printf("DeliverPost repush data to redis: %+v\n", res)
	}

	elapsed := time.Since(start)
	log.Println("该函数执行完成耗时：", elapsed)
}
