package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log2 "github.com/hoang-hs/base/src/common/log"
	"github.com/hoang-hs/base/src/pkg/alert"
	"runtime"
	"strings"
	"time"

	"connectrpc.com/connect"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/workflow"
)

const (
	frameCountCaller = 4
)

var (
	ErrPanic = errors.New("panic occurred")
)

type RecoveryInterceptor struct {
	sender alert.Sender
	interceptor.WorkerInterceptorBase
}

func NewRecoveryInterceptor(sender alert.Sender) *RecoveryInterceptor {
	return &RecoveryInterceptor{sender: sender}
}

// for api
func (r *RecoveryInterceptor) Unary() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (resp connect.AnyResponse, err error) {
			defer func() {
				if errRecover := recover(); errRecover != nil {
					text := fmt.Sprintf("[Panic]\ntime: [%s]\nerror: [%v] \nstack:\n%v \nrequest: [%v]",
						time.Now().Format(time.RFC3339), errRecover, getCallers(4), req.Spec().Procedure)
					log2.Error("panic", log2.String("error", text))
					r.sender.SendMessage(ctx, text)
					err = connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", ErrPanic, errRecover))
				}
			}()
			return next(ctx, req)
		}
	}
}

// for temporal worker
var _ interceptor.WorkerInterceptor = (*RecoveryInterceptor)(nil)

func (r *RecoveryInterceptor) InterceptActivity(_ context.Context, next interceptor.ActivityInboundInterceptor) interceptor.ActivityInboundInterceptor {
	return &ActivityRecoveryInbound{
		sender:                     r.sender,
		ActivityInboundInterceptor: &interceptor.ActivityInboundInterceptorBase{Next: next},
	}
}

func (r *RecoveryInterceptor) InterceptWorkflow(_ workflow.Context, next interceptor.WorkflowInboundInterceptor) interceptor.WorkflowInboundInterceptor {
	return &interceptor.WorkflowInboundInterceptorBase{Next: next}
}

type ActivityRecoveryInbound struct {
	sender alert.Sender
	interceptor.ActivityInboundInterceptor
}

func (r ActivityRecoveryInbound) ExecuteActivity(ctx context.Context, in *interceptor.ExecuteActivityInput) (interface{}, error) {
	defer func() {
		if errRecover := recover(); errRecover != nil {
			info := activity.GetInfo(ctx)
			text := fmt.Sprintf("[Panic]\ntime: [%s]\nerror: [%v] \nstack:\n%v\nactivity: [%v]",
				time.Now().Format(time.RFC3339), errRecover, getCallers(4), info.ActivityType.Name)
			log2.Error("panic", log2.String("error", text))
			r.sender.SendMessage(ctx, text)
			panic(errRecover)
		}
	}()
	return r.ActivityInboundInterceptor.ExecuteActivity(ctx, in)
}

func (r *RecoveryInterceptor) RecoveryConsumer(ctx context.Context, m *kafka.Message) {
	if errRecover := recover(); errRecover != nil {
		text := fmt.Sprintf("[Panic]\ntime: [%s]\nerror: [%v] \nstack:\n%v\ntopic: [%s]\npartition: [%d]\noffset: [%v]",
			time.Now().Format(time.RFC3339), errRecover, getCallers(4), *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		log2.Error("panic", log2.String("error", text))
		r.sender.SendMessage(ctx, text)
	}
}

//nolint:unparam
func getCallers(start int) string {
	var pcs [frameCountCaller]uintptr
	n := runtime.Callers(start, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var callers []string
	for {
		frame, more := frames.Next()
		callers = append(callers, fmt.Sprintf("%s:%d", frame.Function, frame.Line))
		if !more {
			break
		}
	}

	return strings.Join(callers, "\n")
}
