// Copyright 2020 The Bucketr Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func bucketsRoute(r chi.Router) {
	r.Use(jwtauth.Verifier(tokenAuth))
	r.Use(jwtauth.Authenticator)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())

		// Find all buckets belonging to the requester.
		ctx, _ := context.WithTimeout(context.Background(), queryTimeout)
		cur, err := mng.Database("bucketr").Collection("buckets").Find(
			ctx, bson.M{
				"owner": claims["username"],
			},
		)
		if err != nil {
			chk("[warn] get buckets", err)
			respond(w, http.StatusInternalServerError, bucketsResponse{
				genericResponse: genericResponse{
					Error: "couldn't get buckets",
				},
			})
			return
		}

		var buckets []dbBucket
		if err := cur.All(ctx, &buckets); err != nil {
			chk("[warn] cur.all", err)
			respond(w, http.StatusInternalServerError, bucketsResponse{
				genericResponse: genericResponse{
					Error: "error decoding buckets",
				},
			})
			return
		}

		respond(w, http.StatusOK, bucketsResponse{
			Buckets: buckets,
		})
	})

	r.Route("/{bucket}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())

			var bucket dbBucket
			ctx, _ := context.WithTimeout(context.Background(), queryTimeout)
			// Find one bucket owned by the requester with the name specified in the url param.
			if err := mng.Database("bucketr").Collection("buckets").FindOne(ctx, bson.M{
				"name":  chi.URLParam(r, "bucket"),
				"owner": claims["username"],
			}).Decode(&bucket); err != nil && err != mongo.ErrNoDocuments {
				chk("[warn] get bucket", err)
				respond(w, http.StatusInternalServerError, bucketResponse{
					genericResponse: genericResponse{
						Error: "couldn't get bucket",
					},
				})
				return
			}

			respond(w, http.StatusOK, bucketResponse{
				Bucket: bucket,
			})
		})

		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			bucket := chi.URLParam(r, "bucket")

			// Delete the bucket owned by the requester with the name specified in the url param.
			ctx, _ := context.WithTimeout(context.Background(), queryTimeout)
			if _, err := mng.Database("bucketr").Collection("buckets").DeleteOne(ctx, bson.M{
				"name":  bucket,
				"owner": claims["username"],
			}); err != nil {
				chk("[warn] delete bucket", err)
				respond(w, http.StatusInternalServerError, bucketResponse{
					genericResponse: genericResponse{
						Error: "couldn't delete bucket",
					},
				})
				return
			}

			respond(w, http.StatusOK, bucketResponse{
				Bucket: dbBucket{
					Name:  bucket,
					Owner: claims["username"].(string),
				},
			})
		})

		r.Put("/{key}", func(w http.ResponseWriter, r *http.Request) {
			var req keyRequest
			decoder := json.NewDecoder(r.Body)
			_ = decoder.Decode(&req)

			_, claims, _ := jwtauth.FromContext(r.Context())
			bucket := chi.URLParam(r, "bucket")
			key := chi.URLParam(r, "key")

			// Upsert a key in the bucket. This creates the specified bucket if required.
			ctx, _ := context.WithTimeout(context.Background(), queryTimeout)
			upsert := true
			if _, err := mng.Database("bucketr").Collection("buckets").UpdateOne(
				ctx, bson.M{
					"name":  bucket,
					"owner": claims["username"],
				},
				bson.D{{
					Key: "$set",
					Value: bson.D{{
						Key:   "store." + key,
						Value: req.Value,
					}},
				}},
				&options.UpdateOptions{
					Upsert: &upsert,
				},
			); err != nil {
				chk("[warn] couldn't upsert bucket", err)
				respond(w, http.StatusConflict, keyResponse{
					genericResponse: genericResponse{
						Error: "couldn't upsert bucket",
					},
				})
				return
			}

			respond(w, http.StatusCreated, keyResponse{
				Key:   key,
				Value: req.Value,
				Bucket: dbBucket{
					Name:  bucket,
					Owner: claims["username"].(string),
				},
			})
		})

		r.Delete("/{key}", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			bucket := chi.URLParam(r, "bucket")
			key := chi.URLParam(r, "key")

			// Delete a key in the bucket specified.
			ctx, _ := context.WithTimeout(context.Background(), queryTimeout)
			if _, err := mng.Database("bucketr").Collection("buckets").UpdateOne(
				ctx,
				bson.M{
					"name":  bucket,
					"owner": claims["username"],
				},
				bson.D{{
					Key: "$unset",
					Value: bson.D{{
						Key:   "store." + key,
						Value: "",
					}},
				}},
			); err != nil {
				chk("[warn] unset key", err)
				respond(w, http.StatusInternalServerError, keyResponse{
					genericResponse: genericResponse{
						Error: "couldn't unset key",
					},
				})
				return
			}

			respond(w, http.StatusOK, keyResponse{
				Key: key,
				Bucket: dbBucket{
					Name:  bucket,
					Owner: claims["username"].(string),
				},
			})
		})
	})
}
