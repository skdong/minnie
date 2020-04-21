/*
Copyright 2020 The skdong Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"log"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/spf13/cobra"
)

func NewAPIServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "minnie-apiserver",
		Long: `The Minnie API server
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Print("Hello Minnie")
			RunApi()
			return nil
		},
	}
	return cmd
}

type Resource struct {
	Id, Name string
}

type ResourceList struct {
	Resources []Resource
}

func NewResourceService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path("/resources").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Filter(webserviceLogging).Filter(measureTime)
	ws.Route(ws.GET("").Filter(NewCountFilter().routeConter).To(getAllResources))
	ws.Route(ws.GET("/{user-id}").Filter(routeLogging).Filter(NewCountFilter().routeConter).To(findResouce))
	return ws
}

func globalLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Printf("[global-filter (logger)] %s,%s\n", req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
}

func webserviceLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Printf("[webservice-filter (logger)] %s,%s\n", req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
}

func measureTime(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	now := time.Now()
	chain.ProcessFilter(req, resp)
	log.Printf("[webservice-filter (filter)] %v\n", time.Now().Sub(now))
}

func routeLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Printf("[route-filter (logger)] %s,%s\n", req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
}

type CountFilter struct {
	count   int
	counter chan int
}

func NewCountFilter() *CountFilter {
	c := new(CountFilter)
	c.counter = make(chan int)
	go func() {
		for {
			c.count += <-c.counter
		}
	}()
	return c
}

func (c *CountFilter) routeConter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	c.counter <- 1
	log.Printf("[route-filter (counter)] count: %d", c.count)
	chain.ProcessFilter(req, resp)
}

func getAllResources(request *restful.Request, response *restful.Response) {
	log.Print("getAllUsers")
	response.WriteEntity(ResourceList{[]Resource{{"42", "User"}, {"3.14", "Project"}}})
}

func findResouce(request *restful.Request, response *restful.Response) {
	log.Printf("findUser")
	response.WriteEntity(Resource{"42", "User"})
}

func RunApi() {
	restful.Filter(globalLogging)
	restful.Add(NewResourceService())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
