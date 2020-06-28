// Copyright 2020 The Bucketr Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var tokenAuth *jwtauth.JWTAuth
var mng *mongo.Client

const (
	queryTimeout = time.Second * 2
	warnTag      = "[warn]"
)

func main() {
	godotenv.Load()

	// Open the database.
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var err error
	mng, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	chk("open db", err)

	// Ping the db to ensure we have a connection.
	ctx, _ = context.WithTimeout(context.Background(), queryTimeout)
	chk("ping db", mng.Ping(ctx, readpref.Primary()))
	fmt.Println("db: open")

	// Seed the authenticator.
	signKey := make([]byte, 32)
	_, err = rand.Read(signKey)
	chk("rand read", err)
	tokenAuth = jwtauth.New("HS256", signKey, nil)

	// Wait for a sigint.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		fmt.Println("cleanup")
		// TODO: Softly close the program.
		os.Exit(0)
	}()

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	fmt.Println("http: ready to listen on", addr)
	chk("http", http.ListenAndServe(addr, router()))
}
