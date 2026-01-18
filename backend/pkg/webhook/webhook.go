package webhook

import (
	"log"

	"piemdm/pkg/webhook/task"
)

// ProviderSet is cron providers.
// var ProviderSet = wire.NewSet(NewWebhook, task.NewScanner)

type Webhook struct {
	Scanner *task.Scanner
}

func NewWebhook(scanner *task.Scanner) *Webhook {
	return &Webhook{
		Scanner: scanner,
	}
	// var rdb *redis.Client
	// ctx := context.Background()
	// result := rdb.BRPop(ctx, 600, "WebhookQueue")
	// fmt.Printf("result: %#v\n\n", result)

	// go func() {
	// 	listName := "WebhookQueue"

	// 	for {
	// 		// total, err := redis.LLen(listName)
	// 		// if err != nil {
	// 		// 	log.Println(err)
	// 		// 	return
	// 		// }
	// 		// if total == 0 {
	// 		// 	time.Sleep(2 * time.Second)
	// 		// 	log.Println("total:", time.Now(), total)
	// 		// 	continue
	// 		// }

	// 		// LPOP：在队列最左侧取出一个值(移除并返回列表的第一个元素)
	// 		ctx := context.Background()
	// 		result := rdb.BRPop(ctx, 600, listName)
	// 		fmt.Printf("result: %#v\n\n", result)
	// 		// if err != nil {
	// 		// 	log.Println("redis BRPop timeout, reconnected. ", time.Now(), err)
	// 		// 	continue
	// 		// }
	// 		// resultByte := result[1].([]uint8)

	// 		// req := WebhookReq{}
	// 		// json.Unmarshal([]byte(resultByte), &req)
	// 		// if err != nil {
	// 		// 	log.Println("unmarshal: ", err)
	// 		// 	return
	// 		// }

	// 		// log.Printf("redis->BRPop: %+v\n", req)

	// 		// go DeliverOne(req)
	// 	}

	// }()

	// log.Println("initialize webhook done.")
}

func (s *Webhook) Start() error {
	log.Println("Webhook Start.")
	s.Scanner.Run()
	return nil
}

// func DeliverOne(webhookReq WebhookReq) {
// 	// 根据entity_id 进行push

// 	// 获取 Webhook config
// 	var sw service.WebhookService
// 	hook := sw.GetWebhook(webhookReq.HookID)

// 	// 获取 entity
// 	var es service.EntityService
// 	entity, _, relation := es.GetEntity(webhookReq.TableCode, webhookReq.EntityID, []string{})
// 	// log.Printf("entity: %+v\n", entity)
// 	log.Printf("relation: %+v\n", relation)

// 	go DeliverPost(webhookReq, hook, entity)

// }

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
// func DeliverPost(webhookReq WebhookReq, hook *model.Webhook, entity map[string]any) {
// 	start := time.Now() // 获取当前时间
// 	guuid := uuid.NewV4()
// 	deliveryCode := strings.ToUpper(guuid.String())
// 	// 转成Json 进行进行POST调用
// 	entityJson, _ := json.Marshal(entity)

// 	posturl := hook.Url
// 	req, err := http.NewRequest("POST", posturl, bytes.NewReader(entityJson))
// 	if err != nil {
// 		log.Fatalf("impossible to build request: %s", err)
// 	}
// 	req.Header.Set("Content-Type", hook.ContentType+"; charset=UTF-8")
// 	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	req.Header.Set("X-GitHub-Delivery", deliveryCode)
// 	req.Header.Set("X-GitHub-Event", "Release")
// 	req.Header.Set("X-GitHub-Hook-ID", fmt.Sprintf("%d", hook.ID))
// 	// User-Agent: GitHub-Hookshot/044aadd
// 	// req.Header.Set("X-GitHub-Hook-Installation-Target-ID", "667721582")
// 	// req.Header.Set("X-GitHub-Hook-Installation-Target-Type", "repository")

// 	// create http client
// 	// do not forget to set timeout; otherwise, no timeout!
// 	client := http.Client{Timeout: 10 * time.Second}
// 	// send the request
// 	res, err := client.Do(req)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer res.Body.Close()

// 	// get request Header
// 	reqKeys := make([]string, 0)
// 	for key, _ := range req.Header {
// 		reqKeys = append(reqKeys, key)
// 	}
// 	sort.Strings(reqKeys)
// 	reqHeader := ""
// 	for _, key := range reqKeys {
// 		reqHeader = reqHeader + fmt.Sprintf("%+v: %+v\n", key, req.Header[key][0])
// 	}

// 	// get request Body
// 	reqBodyCopy, err := req.GetBody()
// 	if err != nil {
// 		log.Println("return a new copy of Body error: ", err)
// 	}
// 	reqBodyByte, err := io.ReadAll(reqBodyCopy)
// 	if err != nil {
// 		log.Println("read request body error: ", err)
// 	}
// 	reqBody := string(reqBodyByte)

// 	// get response Header
// 	resHeader := ""
// 	resKeys := make([]string, 0)
// 	for key, _ := range res.Header {
// 		resKeys = append(resKeys, key)
// 	}
// 	sort.Strings(resKeys)
// 	for _, key := range resKeys {
// 		resHeader = resHeader + fmt.Sprintf("%+v: %+v\n", key, res.Header[key][0])
// 	}

// 	// get response Body
// 	resBodyByte, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	// save to database
// 	delivery := model.WebhookDelivery{
// 		HookID:          hook.ID,
// 		DeliveryCode:    deliveryCode,
// 		Event:           "Release",
// 		EntityID:        webhookReq.EntityID,
// 		RequestHeaders:  reqHeader,
// 		RequestPayload:  reqBody,
// 		ResponseStatus:  res.StatusCode,
// 		ResponseHeaders: resHeader,
// 		ResponseBody:    string(resBodyByte),
// 		DeliveredAt:     time.Now(),
// 		CompletedAt:     time.Now(),
// 	}
// 	var se service.WebhookDeliveryService
// 	ok := se.AddWebhookDelivery(&delivery)
// 	if !ok {
// 		// 保存失败的时候，重新把数据放入队列
// 		log.Println("DeliverPost Save to Database err. ")

// 		// 将 map[string]any 转换为 JSON 字符串
// 		dataJson, err := json.Marshal(webhookReq)
// 		if err != nil {
// 			log.Println("DeliverPost Failed to marshal data to JSON:", err)
// 			return
// 		}

// 		// 存储消息
// 		res, err := redis.LPush("WebhookQueue", dataJson)
// 		if err != nil {
// 			log.Printf("DeliverPost redis.LPush Error: %+v\n", err)
// 		}
// 		log.Printf("DeliverPost repush data to redis: %+v\n", res)

// 	}

// 	elapsed := time.Since(start)
// 	log.Println("该函数执行完成耗时：", elapsed)

// }
