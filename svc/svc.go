/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package svc

func InitCLogSVC(grpcAddress string, usingTLS bool) (CLogSVC, error) {
	sr := &srCLog{
		address:  grpcAddress,
		usingTLS: usingTLS,
	}

	err := sr.buildConnection()
	if err != nil {
		return nil, err
	}

	return sr, nil
}
