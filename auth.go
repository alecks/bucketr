// Copyright 2020 The Bucketr Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func authRoute(r chi.Router) {
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		decoder := json.NewDecoder(r.Body)
		_ = decoder.Decode(&req)

		if req.Username == "" || req.Password == "" {
			respond(w, http.StatusBadRequest, loginResponse{
				genericResponse: genericResponse{
					Error: "username and password are required",
				},
			})
			return
		}
		// bcrypt is limited to 72 bytes.
		if len(req.Password) > 72 {
			respond(w, http.StatusBadRequest, loginResponse{
				genericResponse: genericResponse{
					Error: "password can't be greater than 72 characters",
				},
			})
			return
		}

		// Find a user with the specified username.
		var user dbUser
		ctx, _ := context.WithTimeout(context.Background(), queryTimeout)
		err := mng.Database("bucketr").Collection("users").FindOne(ctx, bson.M{
			"username": req.Username,
		}).Decode(&user)

		stat := http.StatusOK
		// If the user doesn't exist, create it.
		if err == mongo.ErrNoDocuments || user.Password == "" {
			// Hash the password.
			hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				chk("[warn] couldn't hash password", err)
				respond(w, http.StatusInternalServerError, loginResponse{
					genericResponse: genericResponse{
						Error: "couldn't hash password",
					},
				})
				return
			}

			user.Username = req.Username
			user.Password = string(hash)

			// Insert the new user into the db.
			ctx, _ = context.WithTimeout(context.Background(), queryTimeout)
			if _, err := mng.Database("bucketr").Collection("users").InsertOne(ctx, user); err != nil {
				chk("[warn] create user", err)
				respond(w, http.StatusInternalServerError, loginResponse{
					genericResponse: genericResponse{
						Error: "couldn't create user",
					},
				})
				return
			}

			stat = http.StatusCreated
		} else if err := bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(req.Password),
		); err != nil {
			respond(w, http.StatusForbidden, loginResponse{
				genericResponse: genericResponse{
					Error: "invalid password",
				},
			})
			return
		}

		// Generate the jwt and return it.
		_, tknString, _ := tokenAuth.Encode(jwt.MapClaims{"username": req.Username})
		respond(w, stat, loginResponse{
			Token: tknString,
		})
	})
}
