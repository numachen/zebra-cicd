package middleware

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

type ctxKey string

const tracerKey ctxKey = "zipkin-tracer"

// InitZipkin 初始化 Zipkin tracer，并返回 tracer 与 reporter closer
// zipkinURL: e.g. "http://localhost:9411/api/v2/spans"
// serviceName: e.g. "zebra-cicd"
// hostPort: e.g. "192.168.30.85:9527"
func InitZipkin(zipkinURL, serviceName, hostPort string) (*zipkin.Tracer, io.Closer, error) {
	reporter := zipkinhttp.NewReporter(zipkinURL)

	endpoint, err := zipkin.NewEndpoint(serviceName, hostPort)
	if err != nil {
		_ = reporter.Close()
		return nil, nil, err
	}

	tracer, err := zipkin.NewTracer(
		reporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithSampler(zipkin.AlwaysSample),
	)
	if err != nil {
		_ = reporter.Close()
		return nil, nil, err
	}
	return tracer, reporter, nil
}

// ZipkinMiddleware 从请求 header 抽取（支持 b3 single-header 或 multi-header via ExtractHTTP），
// 创建 SERVER span 并注入到请求上下文
func ZipkinMiddleware(tracer *zipkin.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		if tracer == nil {
			// 如果 tracer 未初始化，直接继续以避免空指针 panic
			c.Next()
			return
		}

		var span zipkin.Span

		// 使用 b3.ExtractHTTP 从 http.Header 中抽取 SpanContext（兼容 single-header 和 X-B3-*）
		if sc, err := b3.ExtractHTTP(c.Request.Header); err == nil {
			// ExtractHTTP 返回的 sc 可以直接作为 Parent 传入（视具体版本返回类型）
			span = tracer.StartSpan(
				c.Request.URL.Path,
				zipkin.Kind("SERVER"),
				zipkin.Parent(sc),
				zipkin.WithTimestamp(time.Now()),
			)
		}

		if span == nil {
			// 未抽取到上游上下文则新建根 span
			span = tracer.StartSpan(
				c.Request.URL.Path,
				zipkin.Kind("SERVER"),
				zipkin.WithTimestamp(time.Now()),
			)
		}
		defer span.Finish()

		// 把 tracer 放进 context，方便后续创建子 span 使用
		ctx := context.WithValue(c.Request.Context(), tracerKey, tracer)
		// 把当前 span 放进 zipkin context（供 zipkin.SpanFromContext 使用）
		ctx = zipkin.NewContext(ctx, span)
		c.Request = c.Request.WithContext(ctx)

		// 将 b3 header 写回响应，便于链路调试（可选）
		_ = b3.InjectHTTP(c.Writer.Header(), span.Context())

		c.Next()
	}
}

// TracerFromContext 从 context 中取出 *zipkin.Tracer（如果存在）
func TracerFromContext(ctx context.Context) *zipkin.Tracer {
	if v := ctx.Value(tracerKey); v != nil {
		if t, ok := v.(*zipkin.Tracer); ok {
			return t
		}
	}
	return nil
}

// StartChildSpan 从传入的 context 创建一个子 span，并把它和新的 context 返回
// - name: 子 span 名称，例如 "mysql.query"
func StartChildSpan(ctx context.Context, name string) (context.Context, zipkin.Span) {
	// 优先从 context 中取 tracer
	tracer := TracerFromContext(ctx)
	parent := zipkin.SpanFromContext(ctx)

	// 如果没有 tracer，但有 parent span，尝试从 parent 取 tracer（fallback）
	if tracer == nil && parent != nil {
		if t := parent.Tracer(); t != nil {
			tracer = t
		}
	}

	if tracer == nil {
		// 无 tracer：不能创建真实的 zipkin span
		return ctx, nil
	}

	if parent != nil {
		child := tracer.StartSpan(name, zipkin.Parent(parent.Context()), zipkin.WithTimestamp(time.Now()))
		ctx = zipkin.NewContext(ctx, child)
		return ctx, child
	}

	child := tracer.StartSpan(name, zipkin.WithTimestamp(time.Now()))
	ctx = zipkin.NewContext(ctx, child)
	return ctx, child
}

// InjectToHeader 将当前 span 的 b3 信息注入到 http.Request Header（用于 outbound requests）
func InjectToHeader(req *http.Request, span zipkin.Span) {
	if span == nil || req == nil {
		return
	}
	_ = b3.InjectHTTP(req.Header, span.Context())
}
