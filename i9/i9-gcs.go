/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package i9

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type srGCS struct {
	ctx     context.Context
	handler *storage.BucketHandle
}

func newGcsInstance(bucketName string, credential []byte) (*srGCS, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(credential))
	if err != nil {
		return nil, err
	}

	handler := client.Bucket(bucketName)

	instance := &srGCS{
		ctx:     ctx,
		handler: handler,
	}

	return instance, nil
}

func (slf *srGCS) Write(destination string, data []byte, timeout ...time.Duration) error {
	ctx := context.Background()
	if len(timeout) > 0 {
		c, cancel := context.WithTimeout(ctx, timeout[0])
		defer cancel()

		ctx = c
	}

	writer := slf.handler.Object(destination).NewWriter(ctx)
	defer func() {
		_ = writer.Close()
	}()

	_, err := io.Copy(writer, bytes.NewReader(data))
	return err
}

func (slf *srGCS) WriteFromFile(destination string, fromFilePath string, autoIncludeExtension bool, timeout ...time.Duration) error {
	data, err := os.ReadFile(fromFilePath)
	if err != nil {
		return err
	}

	if autoIncludeExtension {
		ext := filepath.Ext(fromFilePath)
		destination += ext
	}

	return slf.Write(destination, data, timeout...)
}
