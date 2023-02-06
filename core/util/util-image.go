/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package util

import (
	"bytes"
	"encoding/base64"
	"image"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func (*srUtil) ImageBase64Decoder(imageBase64 string) (buff bytes.Buffer, config image.Config, format string, err error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(imageBase64))

	buff = bytes.Buffer{}
	_, err = buff.ReadFrom(reader)
	if err != nil {
		return
	}

	config, format, err = image.DecodeConfig(bytes.NewReader(buff.Bytes()))
	return
}
