// Copyright 2020 The Bucketr Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type genericResponse struct {
	Error string `json:"error,omitempty"`
}

type loginResponse struct {
	genericResponse
	Token string `json:"token,omitempty"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type bucketsResponse struct {
	genericResponse
	Buckets []dbBucket `json:"buckets,omitempty"`
}

type bucketResponse struct {
	genericResponse
	Bucket dbBucket `json:"bucket,omitempty"`
}

type keyResponse struct {
	genericResponse
	Key    string   `json:"key,omitempty"`
	Value  string   `json:"value,omitempty"`
	Bucket dbBucket `json:"bucket,omitempty"`
}

type keyRequest struct {
	Value string `json:"value"`
}

type dbUser struct {
	Username string `bson:"username" json:"username,omitempty"`
	Password string `bson:"password" json:"-"`
}

type dbBucket struct {
	Name  string            `bson:"name" json:"name,omitempty"`
	Owner string            `bson:"owner" json:"owner,omitempty"`
	Store map[string]string `bson:"store" json:"store,omitempty"`
}
