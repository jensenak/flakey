package flakey

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

var (
	address = "http://localhost:5000"
)

func Workflow(ctx workflow.Context) error {
	rp := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 1, // Don't back off
		MaximumInterval:    time.Second * 2,
		MaximumAttempts:    30,
	}
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
		HeartbeatTimeout:    60 * time.Second,
		RetryPolicy:         rp,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("Flakey workflow started")

	var color string
	err := workflow.ExecuteActivity(ctx, Start).Get(ctx, &color)
	if err != nil {
		logger.Error("Start failed", err)
		return err
	}

	var steps []string
	err = workflow.ExecuteActivity(ctx, GetSteps, color).Get(ctx, &steps)
	if err != nil {
		logger.Error("Get steps failed", err)
		return err
	}

	var data string
	err = workflow.ExecuteActivity(ctx, RunSteps, steps).Get(ctx, &data)
	if err != nil {
		logger.Error("Run steps failed", err)
		return err
	}

	err = workflow.ExecuteActivity(ctx, Submit, color, data).Get(ctx, nil)
	if err != nil {
		logger.Error("Submit failed", err)
		return err
	}

	logger.Info("Done")

	return nil
}

func Start(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Start")

	resp, err := http.Get(fmt.Sprintf("%s/", address))
	if err != nil {
		logger.Error("Error in start", err)
		return "", err
	}
	defer resp.Body.Close()

	var data map[string]string

	json.NewDecoder(resp.Body).Decode(&data)

	if color, ok := data["color"]; ok {
		return color, nil
	}

	return "", nil
}

func GetSteps(ctx context.Context, color string) ([]string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("GetSteps")

	resp, err := http.Get(fmt.Sprintf("%s/color/%s", address, color))
	if err != nil {
		logger.Error("Error getting steps", err)
		return []string{}, err
	}
	defer resp.Body.Close()

	var data map[string][]string

	json.NewDecoder(resp.Body).Decode(&data)

	if steps, ok := data["steps"]; ok {
		return steps, nil
	}

	return []string{}, nil
}

func RunSteps(ctx context.Context, steps []string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("RunSteps")

	var result string

	for _, step := range steps {
		resp, err := http.Get(fmt.Sprintf("%s/step/%s", address, step))
		if err != nil {
			logger.Error("Error running step", step, err)
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			logger.Error("RunSteps failed!", step)
			return "", errors.New(resp.Status)
		}
		var data map[string]string

		json.NewDecoder(resp.Body).Decode(&data)
		if word, ok := data["word"]; ok {
			result = fmt.Sprintf("%s%s", result, word)
		}
	}

	return result, nil
}

func Submit(ctx context.Context, color, data string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Submit")

	reqBody, err := json.Marshal(map[string]string{
		"color": color,
		"data":  data,
	})
	if err != nil {
		logger.Error("Can't make json")
		return err
	}

	resp, err := http.Post(fmt.Sprintf("%s/done", address), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		logger.Error("Couldn't submit", err)
		return err
	}

	resp.Body.Close()
	return nil
}
