// Copyright 2020 The Bucketr Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func useMiddleware(r chi.Router) {
	ratelimitRequests, err := strconv.Atoi(os.Getenv("RATELIMIT_REQUESTS"))
	if err != nil {
		ratelimitRequests = 100
		chk("[warn] RATELIMIT_REQUESTS = 100", err)
	}

	// Basic middleware stack; these can generally be disabled.
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	// Ratelimit clients; if you exceed RATELIMIT_REQUESTS in 1 min, you'll be ratelimited.
	r.Use(httprate.LimitByIP(ratelimitRequests, 1*time.Minute))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// Set a 30 second timeout. Clients will receive a 502 Gateway Timeout response.
	r.Use(middleware.Timeout(30 * time.Second))
	// CORS configuration. This should probably be stricter.
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))
}
