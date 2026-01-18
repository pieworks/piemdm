package resp

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tomnomnom/linkheader"
)

// type response struct {
// 	Code    int         `json:"code"`
// 	Message string      `json:"message"`
// 	Data    interface{} `json:"data"`
// }

func HandleSuccess(c *gin.Context, data interface{}) {
	SetHeader(c)

	if data == nil {
		data = map[string]string{}
	}

	// resp := response{Code: 0, Message: "success", Data: data}
	// c.JSON(http.StatusOK, resp)
	c.JSON(http.StatusOK, data)
}

func HandleError(c *gin.Context, httpCode int, message string, data interface{}) {
	SetHeader(c)

	if data == nil {
		data = map[string]string{}
	}
	// resp := response{Code: code, Message: message, Data: data}
	// c.JSON(httpCode, resp)
	// 如果 data 中有 errors 则输出errors
	// 如果 data 中有 rate 则输出 rate
	// 如果 data 中有 documentation_url 则输出 documentation_url
	c.JSON(httpCode, gin.H{
		"message": message,
		// "errors": []map[string]string{
		// 	{
		// 		"resource": "CommitComment",
		// 		"field":    "body",
		// 		"code":     "blank",
		// 	},
		// },
		// "rate": map[string]int{
		// 	"limit":     5000,
		// 	"remaining": 0,
		// 	"reset":     1634234567, // UTC 时间戳，表示速率限制重置时间
		// },
		// "documentation_url": "https://docs.github.com/rest/repos/commits#create-a-commit-comment"
	})
}

func SetHeader(c *gin.Context) {
	c.Header("Accept", "application/json")
	c.Header("Content-Type", "application/json; charset=utf-8")
	// 可以在 Http 文件中添加 cors 跨域设置
	// 已在 middleware/cors.go 中统一处理,此处移除以避免冲突
	// c.Header("Access-Control-Allow-Origin", "*")
	// c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	// c.Header("Cache-Control", "private, max-age=60, s-maxage=60")
	// c.Header("Authentication-Token-Expiration", "2025-08-09 18:35:36 +0800")
	// c.Header("Link", "<http://127.0.0.1:8081/articles/index?per_page=10&page=2&direction=asc&sort=full_name>; rel=\"next\", <http://127.0.0.1:8081/articles/index?per_page=10&page=682&direction=asc&sort=full_name>; rel=\"last\"")
	// c.Header("X-Accepted-Permissions", "metadata=read")
	c.Header("X-Ratelimit-Limit", "5000")
	c.Header("X-Ratelimit-Remaining", "4983")
	c.Header("X-Ratelimit-Reset", "1746963906")
	c.Header("X-Ratelimit-Used", "17")
	c.Header("X-Ratelimit-Resource", "core")
	// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubdomains; preload")
	// c.Header("Referrer-Policy", "origin-when-cross-origin, strict-origin-when-cross-origin")
	// c.Header("Content-Security-Policy", "default-src 'none'")
	// c.Header("Transfer-Encoding", "chunked")
	// 可以在最前的地方生成
	c.Header("X-Request-Id", "9F59:18EA40:6CC22A:84B06B:68208D6D")
	// 在默认的请求上， 浏览器只能访问以下默认的 响应头:
	// - Cache-Control
	// - Content-Language
	// - Content-Type
	// - Expires
	// - Last-Modified
	// - Pragma
	// 如果想让浏览器能访问到其他的 响应头的话 需要在服务器上设置 Access-Control-Expose-Headers
	c.Header("Access-Control-Expose-Headers", "Link, Accept, Content-Type, Authorization, X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Used, X-RateLimit-Resource, X-RateLimit-Reset, X-Request-Id")
}

// GeneratePaginationLinks 生成分页链接
func GeneratePaginationLinks(r *http.Request, page, pageSize, total int) linkheader.Links {
	host := r.Host
	baseUrl := "http://" + host + r.URL.Path
	query := r.URL.Query()
	links := linkheader.Links{}
	// pageTotal := math.Ceil(float64(total) / float64(pageSize))

	if page > 1 {
		query.Set("page", "1")
		links = append(links, linkheader.Link{URL: baseUrl + "?" + query.Encode(), Rel: "first"})
		query.Set("page", strconv.Itoa(page-1))
		links = append(links, linkheader.Link{URL: baseUrl + "?" + query.Encode(), Rel: "prev"})
	}

	query.Set("page", strconv.Itoa(page+1))
	links = append(links, linkheader.Link{URL: baseUrl + "?" + query.Encode(), Rel: "next"})

	// if page < total {
	query.Set("page", strconv.Itoa(int(total)))
	links = append(links, linkheader.Link{URL: baseUrl + "?" + query.Encode(), Rel: "last"})
	// }

	return links
}
