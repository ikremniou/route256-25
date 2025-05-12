package mw

import (
	"context"
	"fmt"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func Validate(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if message, ok := req.(proto.Message); ok {
		if err := protovalidate.Validate(message); err != nil {
			validationErr, ok := err.(*protovalidate.ValidationError)
			if !ok {
				return nil, status.Error(codes.InvalidArgument, "validation failed") // Generic error
			}

			var errorDetails string
			for _, violation := range validationErr.Violations {
				errorDetails += fmt.Sprintf("field: %s, rule: %s; ", violation.FieldValue, violation.RuleValue)
			}

			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	return handler(ctx, req)
}
