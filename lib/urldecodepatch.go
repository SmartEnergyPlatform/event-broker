/*
 * Copyright 2018 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package lib

import (
	"net/http"
	"net/url"
)

func NewUrlDecodeMiddleWare(handler http.Handler) *UrlDecodeMiddleWare {
	return &UrlDecodeMiddleWare{handler: handler}
}

type UrlDecodeMiddleWare struct {
	handler http.Handler
}

func (this *UrlDecodeMiddleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if this.handler != nil {
		r.URL.Path, _ = url.QueryUnescape(r.URL.Path)
		r.URL.RawPath, _ = url.QueryUnescape(r.URL.RawPath)
		this.handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Forbidden", 403)
	}
}
