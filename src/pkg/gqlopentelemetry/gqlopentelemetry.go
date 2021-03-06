package gqlopentelemetry

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type Tracer struct {
	ComplexityExtensionName string
	DataDog                 bool
}

var _ interface {
	graphql.HandlerExtension
	graphql.ResponseInterceptor
	graphql.FieldInterceptor
} = Tracer{}

func (a Tracer) ExtensionName() string {
	return "OpenTelemetry"
}

func (a Tracer) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (a Tracer) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	tracer := otel.Tracer("response")
	ctx, span := tracer.Start(ctx, operationName(ctx))
	defer span.End()
	if !span.IsRecording() {
		return next(ctx)
	}

	oc := graphql.GetOperationContext(ctx)
	span.SetAttributes(
		attribute.String("request.query", oc.RawQuery),
	)
	complexityExtension := a.ComplexityExtensionName
	if complexityExtension == "" {
		complexityExtension = "ComplexityLimit"
	}
	complexityStats, ok := oc.Stats.GetExtension(complexityExtension).(*extension.ComplexityStats)
	if !ok {
		// complexity extension is not used
		complexityStats = &extension.ComplexityStats{}
	}

	if complexityStats.ComplexityLimit > 0 {
		span.SetAttributes(
			attribute.Int64("request.complexityLimit", int64(complexityStats.ComplexityLimit)),
			attribute.Int64("request.operationComplexity", int64(complexityStats.Complexity)),
		)
	}

	for key, val := range oc.Variables {
		span.SetAttributes(
			attribute.String(fmt.Sprintf("request.variables.%s", key), fmt.Sprintf("%+v", val)),
		)
	}

	return next(ctx)
}

func (a Tracer) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	fc := graphql.GetFieldContext(ctx)
	tracer := otel.Tracer("field")
	ctx, span := tracer.Start(ctx, fc.Field.ObjectDefinition.Name+"/"+fc.Field.Name)
	defer span.End()
	if !span.IsRecording() {
		return next(ctx)
	}

	span.SetAttributes(
		attribute.String("resolver.path", fc.Path().String()),
		attribute.String("resolver.object", fc.Field.ObjectDefinition.Name),
		attribute.String("resolver.field", fc.Field.Name),
		attribute.String("resolver.alias", fc.Field.Alias),
	)
	if a.DataDog {
		span.SetAttributes(
			// key from gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext#ResourceName
			attribute.String("resource.name", operationName(ctx)),
		)
	}
	for _, arg := range fc.Field.Arguments {
		if arg.Value != nil {
			span.SetAttributes(
				attribute.String(fmt.Sprintf("resolver.args.%s", arg.Name), arg.Value.String()),
			)
		}
	}

	resp, err := next(ctx)

	errList := graphql.GetFieldErrors(ctx, fc)
	if len(errList) != 0 {
		// UNKNOWN, HTTP Mapping: 500 Internal Server Error
		span.SetStatus(codes.Error, errList.Error())
		span.SetAttributes(
			attribute.Bool("resolver.hasError", true),
			attribute.Int64("resolver.errorCount", int64(len(errList))),
		)
		for idx, err := range errList {
			span.SetAttributes(
				attribute.String(fmt.Sprintf("resolver.error.%d.message", idx), err.Error()),
				attribute.String(fmt.Sprintf("resolver.error.%d.kind", idx), fmt.Sprintf("%T", err)),
			)
		}
	}

	return resp, err
}

func operationName(ctx context.Context) string {
	requestContext := graphql.GetOperationContext(ctx)
	requestName := "nameless-operation"
	if requestContext.Doc != nil && len(requestContext.Doc.Operations) != 0 {
		op := requestContext.Doc.Operations[0]
		if op.Name != "" {
			requestName = op.Name
		}
	}

	return requestName
}
