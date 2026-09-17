package codeexec

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

// Runner invokes the sandboxed Python Lambda.
type Runner struct {
	client *lambda.Client
	fn     string
}

func NewRunner(ctx context.Context, region, fn, accessKey, secretKey string) (*Runner, error) {
	if region == "" || fn == "" || accessKey == "" || secretKey == "" {
		return &Runner{}, nil
	}
	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)
	if err != nil {
		return nil, err
	}
	return &Runner{client: lambda.NewFromConfig(cfg), fn: fn}, nil
}

func (r *Runner) Configured() bool { return r != nil && r.client != nil && r.fn != "" }

// invokePayload matches the Lambda handler's expected event shape.
type invokePayload struct {
	Mode      string      `json:"mode"`
	Source    string      `json:"source"`
	Stdin     string      `json:"stdin,omitempty"`
	Tests     []gradeTest `json:"tests,omitempty"`
	TimeoutMs int         `json:"timeoutMs"`
}

type gradeTest struct {
	Input string `json:"input"`
}

type gradeResponse struct {
	Results []RunResult `json:"results"`
}

func (r *Runner) Grade(ctx context.Context, source string, inputs []string) ([]RunResult, error) {
	
	if !r.Configured() {
		return nil, ErrNotConfigured
	}
	
	tests := make([]gradeTest, len(inputs))
	
	for i, in := range inputs {
		tests[i] = gradeTest{Input: in}
	}
	
	payload, _ := json.Marshal(invokePayload{
		Mode: "grade", Source: source, Tests: tests, TimeoutMs: runTimeoutMs,
	})
	
	out, err := r.client.Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   aws.String(r.fn),
		InvocationType: types.InvocationTypeRequestResponse,
		Payload:        payload,
	})
	
	if err != nil {
		return nil, fmt.Errorf("invoke grader: %w", err)
	}
	
	if out.FunctionError != nil && *out.FunctionError != "" {
		return nil, fmt.Errorf("grader error (%s): %s", *out.FunctionError, string(out.Payload))
	}
	
	var res gradeResponse
	
	if err := json.Unmarshal(out.Payload, &res); err != nil {
		return nil, fmt.Errorf("decode grader result: %w", err)
	}
	
	return res.Results, nil
}

// Run executes source once with optional stdin and returns its output.
func (r *Runner) Run(ctx context.Context, source, stdin string) (RunResult, error) {
	
	if !r.Configured() {
		return RunResult{}, ErrNotConfigured
	}
	
	payload, _ := json.Marshal(invokePayload{
		Mode: "run", Source: source, Stdin: stdin, TimeoutMs: runTimeoutMs,
	})
	
	out, err := r.client.Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   aws.String(r.fn),
		InvocationType: types.InvocationTypeRequestResponse,
		Payload:        payload,
	})
	
	if err != nil {
		return RunResult{}, fmt.Errorf("invoke runner: %w", err)
	}
	
	// A thrown handler surfaces here rather than as a transport error.
	if out.FunctionError != nil && *out.FunctionError != "" {
		return RunResult{}, fmt.Errorf("runner error (%s): %s", *out.FunctionError, string(out.Payload))
	}
	
	var res RunResult
	
	if err := json.Unmarshal(out.Payload, &res); err != nil {
		return RunResult{}, fmt.Errorf("decode runner result: %w", err)
	}
	
	return res, nil
}
